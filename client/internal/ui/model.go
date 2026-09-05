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
	return nil
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
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "1":
			m.focus = Rooms
		case "2":
			m.focus = Users
		case "3":
			m.focus = Chat
		default:
			var cmd tea.Cmd
			switch m.focus {
			case Rooms:
				m.roomsModel, cmd = m.roomsModel.Update(msg)
			case Users:
				m.usersModel, cmd = m.usersModel.Update(msg)
			case Chat:
				var chatMod tea.Model
				chatMod, cmd = m.chatModel.Update(msg)
				m.chatModel = chatMod.(panels.ChatModel)
			}
			return m, cmd
		}
	}
	return m, nil
}

func (m Model) View() string {
	// Decide el focus de la interfaz
	m.usersModel.Focused = m.focus == Users
	m.roomsModel.Focused = m.focus == Rooms
	m.chatModel.Focused = m.focus == Chat

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
	return Model{usersModel: usersModel, roomsModel: roomsModel}
}
