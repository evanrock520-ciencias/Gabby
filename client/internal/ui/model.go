package ui

import (
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

type Model struct {
	width      int
	heigth     int
	focus      Focus
	roomsModel panels.ListModel
	usersModel panels.ListModel
	chatModel  panels.ChatModel
}

func (m Model) Init() tea.Cmd {
	return m.chatModel.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Darle proporción a la interfaz
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.heigth = msg.Height

		// Los paneles tienen bordes que ocupan dos carácteres más de altura y de ancho.
		m.roomsModel.Width = m.width/4 - 2
		m.roomsModel.Height = (m.heigth)/2 - 2

		m.usersModel.Width = m.width/4 - 2
		m.usersModel.Height = (m.heigth)/2 - 2

		m.chatModel.Width = m.width - (m.width / 4) - 2
		m.chatModel.Height = m.heigth - 2

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
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
		switch msg.String() {
		case "q":
			return m, tea.Quit
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
			cmd := m.chatModel.TextInput.Focus()
			return m, cmd
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
			return m, cmd
		}

	default:
		var cmd tea.Cmd
		m.chatModel, cmd = m.chatModel.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *Model) syncFocus() {
	m.usersModel.Focused = m.focus == Users
	m.roomsModel.Focused = m.focus == Rooms
	m.chatModel.Focused = m.focus == Chat
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
	roomsView := m.roomsModel.View()
	usersView := m.usersModel.View()
	chatView := m.chatModel.View()
	sidebar := lipgloss.JoinVertical(lipgloss.Top, roomsView, usersView)
	appView := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, chatView)
	return appView
}

// Función para poblar los datos de la interfaz
func NewModel() Model {
	usersModel := panels.NewListModel("[2] Users", []string{"Yahel", "Derek", "Luis", "Sofia"})
	roomsModel := panels.NewListModel("[1] Rooms", []string{"Sala 1", "Sala 2", "Sala 3", "Sala 4"})
	chatModel := panels.NewChatModel()
	m := Model{
		focus:      Rooms,
		usersModel: usersModel,
		roomsModel: roomsModel,
		chatModel:  chatModel,
	}
	m.syncFocus()
	return m
}
