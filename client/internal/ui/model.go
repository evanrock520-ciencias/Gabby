package ui

import (
	"client/internal/domain"
	"client/internal/protocol"
	"client/internal/ui/messages"
	"client/internal/ui/panels"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Focus int

const (
	Rooms Focus = iota
	Users
	Chat
)

var (
	defaultKeymaps = []panels.Keymap{
		panels.NewKeymap("q", "Exit"),
		panels.NewKeymap("c", "Create Room"),
		panels.NewKeymap("i", "Invitate"),
	}
	modalKeymaps = []panels.Keymap{
		panels.NewKeymap("enter", "Confirm"),
		panels.NewKeymap("esc", "Cancel"),
	}
)

type Model struct {
	width       int
	heigth      int
	focus       Focus
	roomsModel  panels.ListModel
	usersModel  panels.ListModel
	chatModel   panels.ChatModel
	footerModel panels.FooterModel
	activeModal *panels.ModalModel
	session     domain.SessionState
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.chatModel.Init(), m.activeModal.TextInput.Focus())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Darle proporción a la interfaz
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.heigth = msg.Height

		// Los paneles tienen bordes que ocupan dos carácteres más de altura y de ancho.
		m.roomsModel.SetSize(m.width/4-2, (m.heigth/2)-3)
		(m.usersModel.SetSize(m.width/4-2, m.heigth-(m.heigth/2)-3-1))

		m.chatModel.SetSize(m.width-(m.width/4)-2, m.heigth-5)

		m.footerModel.SetSize(m.width, 3)

		if m.activeModal != nil {
			m.activeModal.SetSize(m.width, max(0, m.heigth-3)) // -3 por el footer
		}

	case messages.ListResultMsg:
		switch result := msg.(type) {
		case messages.EnterRoomMsg:
			if result.Roomname != "" {
				if room, ok := m.session.Rooms[result.Roomname]; ok {
					m.chatModel.SetRoom(room)
				}
			}
		case messages.EnterDMMsg:
			if result.Username != "" {
				dmRoom, ok := m.session.DMs[result.Username]
				if !ok {
					newDMRoom := domain.NewRoom("@" + result.Username)
					dmRoom = &newDMRoom
					if m.session.DMs == nil {
						m.session.DMs = make(map[string]*domain.Room)
					}
					m.session.DMs[result.Username] = dmRoom
				}
				m.chatModel.SetRoom(dmRoom)
			}
		}

	case messages.ModalResultMsg:
		m.activeModal = nil
		m.footerModel.Keymaps = defaultKeymaps
		m.syncFocus()

		switch result := msg.(type) {
		case messages.CreateRoomMsg:
			if result.Roomname != "" {
				m.roomsModel.Items = append(m.roomsModel.Items, result.Roomname)
				if m.session.Rooms == nil {
					m.session.Rooms = make(map[string]*domain.Room)
				}
				newRoom := domain.NewRoom(result.Roomname)
				m.session.Rooms[result.Roomname] = &newRoom
			}
		case messages.SetUserMsg:
			if result.Username != "" {
				m.footerModel.Username = result.Username
				m.session.CurrentUser.Username = result.Username
				m.session.CurrentUser.Status = protocol.ACTIVE
				m.chatModel.SetUsername(result.Username)
			} else {
				return m, tea.Quit
			}
		}

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

		// Si hay un modal activo, consume todo el teclado
		if m.activeModal != nil {
			updatedModal, cmd := m.activeModal.Update(msg)
			m.activeModal = &updatedModal
			if !m.activeModal.IsCapturingInput() {
				m.activeModal = nil
				m.footerModel.Keymaps = defaultKeymaps
				m.syncFocus()
			}
			return m, cmd
		}

		// Si el panel activo está en modo captura de texto, se delega todo el teclado
		if active := m.activePanel(); active != nil && active.IsCapturingInput() {
			var cmd tea.Cmd
			switch m.focus {
			case Rooms:
				m.roomsModel, cmd = m.roomsModel.Update(msg)
			case Users:
				m.usersModel, cmd = m.usersModel.Update(msg)
			case Chat:
				m.chatModel, cmd = m.chatModel.Update(msg)
			}
			return m, cmd
		}

		// Modo Navegación
		cmd := m.handleNavegation(msg)
		return m, cmd

	default:
		var cmd tea.Cmd
		if m.activeModal != nil {
			updatedModal, cmd := m.activeModal.Update(msg)
			m.activeModal = &updatedModal
			return m, cmd
		}
		m.chatModel, cmd = m.chatModel.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *Model) openModal(modal panels.ModalModel) tea.Cmd {
	modal.SetSize(m.width, max(0, m.heigth-3))
	m.chatModel.TextInput.Blur()
	m.footerModel.Keymaps = modalKeymaps
	cmd := modal.TextInput.Focus()
	m.activeModal = &modal
	return cmd
}

func (m *Model) syncFocus() {
	m.usersModel.SetFocus(m.focus == Users)
	m.roomsModel.SetFocus(m.focus == Rooms)
	m.chatModel.SetFocus(m.focus == Chat)
}

func (m Model) activePanel() panels.InputCapturer {
	switch m.focus {
	case Rooms:
		return m.roomsModel
	case Users:
		return m.usersModel
	case Chat:
		return m.chatModel
	default:
		return nil
	}
}

func (m Model) View() string {
	var mainView string

	if m.activeModal != nil {
		mainView = m.activeModal.View()
	} else {
		roomsView := m.roomsModel.View()
		usersView := m.usersModel.View()
		chatView := m.chatModel.View()
		sidebar := lipgloss.JoinVertical(lipgloss.Top, roomsView, usersView)
		mainView = lipgloss.JoinHorizontal(lipgloss.Top, sidebar, chatView)
	}

	footerView := m.footerModel.View()
	return lipgloss.JoinVertical(lipgloss.Top, mainView, footerView)
}

// Función para poblar los datos de la interfaz
func NewModel(session domain.SessionState) Model {
	if session.Rooms == nil {
		session.Rooms = make(map[string]*domain.Room)
	}
	if session.DMs == nil {
		session.DMs = make(map[string]*domain.Room)
	}

	var users []string
	for _, user := range session.Users {
		users = append(users, user.Username)
	}

	var rooms []string
	for _, room := range session.Rooms {
		rooms = append(rooms, room.Name)
	}

	usersModel := panels.NewListModel("[2] Users", users, func(value string) messages.ListResultMsg {
		return messages.EnterDMMsg{Username: value}
	})
	roomsModel := panels.NewListModel("[1] Rooms", rooms, func(value string) messages.ListResultMsg {
		return messages.EnterRoomMsg{Roomname: value}
	})
	globalRoom, ok := session.Rooms["Global"]
	if !ok {
		newGlobal := domain.NewRoom("Global")
		globalRoom = &newGlobal
		session.Rooms["Global"] = globalRoom
	}
	chatModel := panels.NewChatModel(globalRoom, session.CurrentUser.Username)
	footerModel := panels.FooterModel{Keymaps: defaultKeymaps, Username: "", Status: protocol.ACTIVE}
	loginModel := panels.NewModalModel("Login", "Username", 8, func(value string) messages.ModalResultMsg {
		return messages.SetUserMsg{Username: value}
	}, false)
	m := Model{
		focus:       Rooms,
		usersModel:  usersModel,
		roomsModel:  roomsModel,
		chatModel:   chatModel,
		footerModel: footerModel,
		activeModal: &loginModel,
		session:     session,
	}
	m.syncFocus()
	return m
}

func (m *Model) handleNavegation(msg tea.KeyMsg) tea.Cmd {
	// Modo Navegación
	switch msg.String() {
	case "q":
		return tea.Quit
	case "1":
		m.focus = Rooms
		m.syncFocus()
		m.chatModel.TextInput.Blur()
	case "2":
		m.focus = Users
		m.syncFocus()
		m.chatModel.TextInput.Blur()
	case "3":
		m.focus = Chat
		m.syncFocus()
		return m.chatModel.TextInput.Focus()
	case "c":
		return m.openModal(panels.NewModalModel("Create Room", "Roomname", 16, func(value string) messages.ModalResultMsg {
			return messages.CreateRoomMsg{Roomname: value}
		}, true))
	case "i":
		// TODO: Implementar una secuencia de Modals tipo wizard (porque está selección tiene 2 pasos)
		return m.openModal(panels.NewModalModel("Invitate", "Roomname", 17, func(value string) messages.ModalResultMsg {
			return messages.InvitateMsg{Roomname: value}
		}, true))
	default:
		var cmd tea.Cmd
		switch m.focus {
		case Rooms:
			m.roomsModel, cmd = m.roomsModel.Update(msg)
		case Users:
			m.usersModel, cmd = m.usersModel.Update(msg)
		case Chat:
			m.chatModel, cmd = m.chatModel.Update(msg)
		}
		return cmd
	}

	return nil
}
