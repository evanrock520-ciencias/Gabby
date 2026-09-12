package panels

import (
	"strings"

	"client/internal/domain"
	"client/internal/ui/messages"
	"client/internal/ui/styles"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ChatModel struct {
	Panel
	TextInput     textinput.Model
	Viewport      viewport.Model
	DisplayedRoom *domain.Room
	Username      string
	ready         bool
}

func NewChatModel(initialRoom *domain.Room, username string) ChatModel {
	ti := textinput.New()
	ti.Placeholder = "Message..."
	ti.Blur()
	ti.CharLimit = 156
	return ChatModel{
		TextInput:     ti,
		DisplayedRoom: initialRoom,
		Username:      username,
	}
}

func (m ChatModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *ChatModel) SetUsername(username string) {
	m.Username = username
}

func (m *ChatModel) SetRoom(room *domain.Room) {
	m.DisplayedRoom = room
	m.refreshMessages()
}

func (m *ChatModel) Blur() {
	m.TextInput.Blur()
}

func (m *ChatModel) Focus() tea.Cmd {
	return m.TextInput.Focus()
}

func (m *ChatModel) refreshMessages() {
	if !m.ready || m.DisplayedRoom == nil {
		return
	}
	var sb strings.Builder
	for _, msg := range m.DisplayedRoom.Messages {
		userRender := styles.HeaderTitleStyle.Render(msg.Username)
		msgRender := styles.TextStyle.Render(msg.Message)

		clientMessageStyle := styles.MessageStyle.Width(m.width - 5)
		if msg.Username == m.Username {
			clientMessageStyle = styles.CurrentUserMessageStyle.Width(m.width - 5).Align(lipgloss.Right)
		}

		sb.WriteString(clientMessageStyle.Render(userRender + "\n" + msgRender))
		sb.WriteString("\n")
	}

	m.Viewport.SetContent(sb.String())
	m.Viewport.GotoBottom()
}

func (m *ChatModel) SetSize(width int, height int) {
	m.Panel.SetSize(width, height)

	viewportWidth := max(0, width-4)
	viewportHeight := max(0, height-4)

	if !m.ready {
		m.Viewport = viewport.New(viewportWidth, viewportHeight)
		m.ready = true
		m.refreshMessages()
	} else {
		m.Viewport.Width = viewportWidth
		m.Viewport.Height = viewportHeight
	}
}

func (m *ChatModel) AddMessage(msg domain.ChatMessage) {
	if m.DisplayedRoom != nil {
		m.DisplayedRoom.AddMessage(msg)
		m.refreshMessages()
	}
}

func (m ChatModel) IsCapturingInput() bool {
	return m.TextInput.Focused()
}

func (m ChatModel) Update(msg tea.Msg) (ChatModel, tea.Cmd) {
	var (
		textinputCmd tea.Cmd
		viewportCmd  tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !m.TextInput.Focused() {
			switch msg.String() {
			case "e":
				if m.DisplayedRoom != nil && m.DisplayedRoom.Name != "Global" && !strings.HasPrefix(m.DisplayedRoom.Name, "@") {
					return m, func() tea.Msg {
						return messages.LeftRoomMsg{
							Roomname: m.DisplayedRoom.Name,
						}
					}
				}
			}
		}

		switch msg.Type {
		case tea.KeyEsc:
			m.TextInput.Blur()
			return m, nil
		case tea.KeyEnter:
			if !m.TextInput.Focused() {
				cmd := m.TextInput.Focus()
				return m, cmd
			}
			text := strings.TrimSpace(m.TextInput.Value())
			if text != "" {
				chatMessage := domain.ChatMessage{Message: text, Username: m.Username}
				m.AddMessage(chatMessage)
			}
			m.TextInput.SetValue("")
			return m, nil
		}
	}

	if m.TextInput.Focused() {
		m.TextInput, textinputCmd = m.TextInput.Update(msg)
	} else {
		m.Viewport, viewportCmd = m.Viewport.Update(msg)
	}

	if _, ok := msg.(tea.MouseMsg); ok {
		m.Viewport, viewportCmd = m.Viewport.Update(msg)
	}

	return m, tea.Batch(textinputCmd, viewportCmd)
}

func (m ChatModel) View() string {
	boxStyle := styles.BoxStyle.Width(m.width).Height(m.height)
	roomName := ""
	if m.DisplayedRoom != nil {
		roomName = m.DisplayedRoom.Name
	}
	header := styles.ChatNameStyle.Width(m.width - 2).Render(roomName)

	if m.IsFocused() {
		boxStyle = boxStyle.BorderForeground(styles.PrimaryColor)
	}

	messageBarStyle := styles.TextInputStyle.Width(m.width - 4) // -4 Por los bordes del panel y del input.
	if m.IsCapturingInput() {
		messageBarStyle = messageBarStyle.BorderForeground(styles.SecondaryColor)
	}

	viewportContent := m.Viewport.View()
	if !m.ready {
		viewportContent = ""
	}

	return boxStyle.Render(header + "\n" + viewportContent + "\n" + messageBarStyle.Render(m.TextInput.View()))
}
