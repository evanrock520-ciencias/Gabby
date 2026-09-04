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
	roomsModel panels.RoomsModel
	usersModel panels.UsersModel
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
