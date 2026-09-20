package ui

import (
	"client/internal/domain"
	"client/internal/network"
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
	activeModal panels.Modal
	session     domain.SessionState
	conn        *network.ConnectionManager
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.chatModel.Init(), m.activeModal.Focus(), m.waitForServerMsg())
}

func (m Model) waitForServerMsg() tea.Cmd {
	return func() tea.Msg {
		return <-m.conn.Messages()
	}
}

func (m Model) sendToServer(msg protocol.ClientMessage) tea.Cmd {
	m.conn.Send(msg)
	return m.waitForServerMsg()
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

		// Delegar el mensaje para que se actualizen los paneles de listas
		// Y calculen correctamente el offset tras un resize
		m.roomsModel, _ = m.roomsModel.Update(msg)
		m.usersModel, _ = m.usersModel.Update(msg)

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

		// Si hay un modal activo, consume todo el teclado
		if m.activeModal != nil {
			updatedModal, cmd := m.activeModal.Update(msg)
			m.activeModal = updatedModal
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
			m.activeModal = updatedModal
			return m, cmd
		}
		m.chatModel, cmd = m.chatModel.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *Model) openModal(modal panels.Modal) tea.Cmd {
	modal.SetSize(m.width, max(0, m.heigth-3))
	m.chatModel.Blur()
	cmd := modal.Focus()
	m.activeModal = modal
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
func NewModel(session domain.SessionState, conn *network.ConnectionManager) Model {
	globalRoom := session.AddRoom("Global")

	usersModel := panels.NewListModel("[2] Users", userItems(session.Users), func(value string) tea.Msg {
		return messages.EnterDMMsg{Username: value}
	})
	roomsModel := panels.NewListModel("[1] Rooms", roomItems(session.Rooms), func(value string) tea.Msg {
		return messages.EnterRoomMsg{Roomname: value}
	})
	chatModel := panels.NewChatModel(globalRoom, session.CurrentUser.Username)
	footerModel := panels.FooterModel{Username: "", Status: domain.ACTIVE}
	loginModel := panels.NewInputModal("Login", "Username", 8, func(value string) tea.Msg {
		return messages.SetUserMsg{Username: value}
	}, false)
	m := Model{
		focus:       Rooms,
		usersModel:  usersModel,
		roomsModel:  roomsModel,
		chatModel:   chatModel,
		footerModel: footerModel,
		activeModal: loginModel,
		session:     session,
		conn:        conn,
	}
	m.syncFocus()
	return m
}

func (m *Model) handleNavegation(msg tea.KeyMsg) tea.Cmd {
	// Modo Navegación
	switch msg.String() {
	case "q":
		m.openModal(panels.NewConfirmModal("Leave the Chat", []string{"Yes", "Not"}, func() tea.Msg {
			return messages.LeftChat{}
		}, func() tea.Msg {
			return nil
		}))
	case "1":
		return m.setFocus(Rooms)
	case "2":
		return m.setFocus(Users)
	case "3":
		return m.setFocus(Chat)
	case "c":
		return m.openModal(panels.NewInputModal("Create Room", "Roomname", 16, func(value string) tea.Msg {
			return messages.CreateRoomMsg{Roomname: value}
		}, true))
	case "i":
		return m.openModal(panels.NewWizardModal(
			func(store *panels.WizardStore) []panels.Modal {
				return []panels.Modal{
					panels.NewListModal("Select a Room", roomItems(m.session.Rooms), false, func(values []string) tea.Msg {
						store.Set("room", values)
						return nil
					}),
					panels.NewListModal("Select Users", userItems(m.session.Users), true, func(values []string) tea.Msg {
						store.Set("users", values)
						return nil
					}),
				}
			},
			func(store panels.WizardStore) tea.Msg {
				rooms, _ := store.Get("room")
				users, _ := store.Get("users")
				return messages.InvitateMsg{Roomname: rooms[0], Users: users}
			},
		))
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
