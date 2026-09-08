package panels

import (
	"strings"

	"client/internal/domain"
	"client/internal/ui/styles"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type ChatModel struct {
	Panel
	TextInput     textinput.Model
	ChatName      string
	Viewport      viewport.Model
	DisplayedRoom domain.Room
	ready         bool
}

func NewChatModel() ChatModel {
	ti := textinput.New()
	ti.Placeholder = "Message..."
	ti.Blur()
	ti.CharLimit = 156
	return ChatModel{
		TextInput:     ti,
		DisplayedRoom: domain.NewRoom("Global"),
	}
}

func (m ChatModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *ChatModel) SetSize(width int, height int) {
	m.Panel.SetSize(width, height)

	viewportWidth := max(0, width-4)
	viewportHeight := max(0, height-4)

	if !m.ready {
		m.Viewport = viewport.New(viewportWidth, viewportHeight)
		m.ready = true
	} else {
		m.Viewport.Width = viewportWidth
		m.Viewport.Height = viewportHeight
	}
}

func (m *ChatModel) AddMessage(msg domain.ChatMessage) {
	m.DisplayedRoom.AddMessage(msg)
	if m.ready {
		var sb strings.Builder
		for _, msg := range m.DisplayedRoom.Messages {
			userRender := styles.HeaderTitleStyle.Render(" " + msg.Username)
			msgRender := styles.TextStyle.Render(msg.Message)
			sb.WriteString(styles.MessageStyle.Render(userRender + "\n " + msgRender))
			sb.WriteString("\n")
		}

		m.Viewport.SetContent(sb.String())
		m.Viewport.GotoBottom()
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
			// TODO: Sincronizar el username con el usuario del cliente
			chatMessage := domain.ChatMessage{Message: text, Username: "Evan"}
			if text != "" {
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
	header := styles.ChatNameStyle.Width(m.width - 2).Render(m.ChatName)

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
