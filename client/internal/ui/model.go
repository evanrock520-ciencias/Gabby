package ui

import (
	"client/internal/domain"
	"client/internal/ui/messages"
	"client/internal/ui/panels"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Focus int

const (
	Rooms Focus = iota
	Users
	Chat
	totalFocus
)

func (f Focus) Next() Focus {
	return (f + 1) % totalFocus
}

func (f Focus) Prev() Focus {
	return (f - 1 + totalFocus) % totalFocus
}

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

	case messages.GlobalResultMsg:
		switch result := msg.(type) {
		case messages.ChangeStatusMsg:
			m.SetStatus(result.Status)
		}

	case messages.ChatResultMsg:
		switch result := msg.(type) {
		case messages.LeftRoomMsg:
			if result.Roomname == "Global" || strings.HasPrefix(result.Roomname, "@") || result.Roomname == "" {
				return m, nil
			}
			m.roomsModel.RemoveItem(result.Roomname)
			m.session.RemoveRoom(result.Roomname)
			if m.chatModel.DisplayedRoom != nil && m.chatModel.DisplayedRoom.Name == result.Roomname {
				if globalRoom, ok := m.session.GetRoom("Global"); ok {
					m.chatModel.SetRoom(globalRoom)
				} else {
					m.chatModel.SetRoom(nil)
				}
			}
		}

	case messages.ListResultMsg:
		switch result := msg.(type) {
		case messages.EnterRoomMsg:
			if result.Roomname != "" {
				if room, ok := m.session.GetRoom(result.Roomname); ok {
					m.chatModel.SetRoom(room)
					m.setFocus(Chat)
				}
			}
		case messages.EnterDMMsg:
			if result.Username != "" {
				dmRoom := m.session.GetOrCreateDM(result.Username)
				m.chatModel.SetRoom(dmRoom)
				m.setFocus(Chat)
				return m, m.chatModel.Focus()
			}
		}

	case messages.ModalResultMsg:
		m.closeModal()

		switch result := msg.(type) {
		case messages.CreateRoomMsg:
			if result.Roomname != "" {
				room := m.session.AddRoom(result.Roomname)
				m.roomsModel.AddItem(NewRoomItem(room))
			}
		case messages.SetUserMsg:
			if result.Username != "" {
				m.SetUsername(result.Username)
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
				m.closeModal()
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
			m.footerModel.SetKeys(m.CurrentKeymaps())
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
	m.chatModel.Blur()
	cmd := modal.TextInput.Focus()
	m.activeModal = &modal
	m.footerModel.SetKeys(m.CurrentKeymaps())
	return cmd
}

func (m *Model) setFocus(newFocus Focus) tea.Cmd {
	m.focus = newFocus
	m.syncFocus()

	m.chatModel.Blur()
	return nil
}

func (m *Model) syncFocus() {
	m.usersModel.SetFocus(m.focus == Users)
	m.roomsModel.SetFocus(m.focus == Rooms)
	m.chatModel.SetFocus(m.focus == Chat)
	m.footerModel.SetKeys(m.CurrentKeymaps())
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
	globalRoom := session.AddRoom("Global")

	usersModel := panels.NewListModel("[2] Users", userItems(session.Users), func(value string) messages.ListResultMsg {
		return messages.EnterDMMsg{Username: value}
	})
	roomsModel := panels.NewListModel("[1] Rooms", roomItems(session.Rooms), func(value string) messages.ListResultMsg {
		return messages.EnterRoomMsg{Roomname: value}
	})
	chatModel := panels.NewChatModel(globalRoom, session.CurrentUser.Username)
	footerModel := panels.FooterModel{Username: "", Status: domain.ACTIVE}
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
		return m.setFocus(Rooms)
	case "2":
		return m.setFocus(Users)
	case "3":
		return m.setFocus(Chat)
	case "c":
		return m.openModal(panels.NewModalModel("Create Room", "Roomname", 16, func(value string) messages.ModalResultMsg {
			return messages.CreateRoomMsg{Roomname: value}
		}, true))
	case "i":
		// TODO: Implementar una secuencia de Modals tipo wizard (porque está selección tiene 2 pasos)
		return m.openModal(panels.NewModalModel("Invitate", "Roomname", 17, func(value string) messages.ModalResultMsg {
			return messages.InvitateMsg{Roomname: value}
		}, true))
	case "s":
		return func() tea.Msg {
			return messages.ChangeStatusMsg{
				Status: m.session.CurrentUser.Status.Next(),
			}
		}
	case "tab":
		m.setFocus(m.focus.Next())
	case "shift+tab":
		m.setFocus(m.focus.Prev())
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

func (m *Model) closeModal() {
	m.activeModal = nil
	m.syncFocus()
}

func (m *Model) SetStatus(status domain.Status) {
	m.session.SetStatus(status)
	m.footerModel.Status = status
}

func (m *Model) SetUsername(username string) {
	m.footerModel.Username = username
	m.session.SetCurrentUser(username, domain.ACTIVE)
	m.chatModel.SetUsername(username)
}

func userItems(users []domain.User) []panels.Item {
	items := make([]panels.Item, len(users))
	for i, user := range users {
		items[i] = NewUserItem(user)
	}
	return items
}

func roomItems(rooms map[string]*domain.Room) []panels.Item {
	items := make([]panels.Item, 0, len(rooms))
	for _, room := range rooms {
		items = append(items, NewRoomItem(room))
	}
	return items
}
