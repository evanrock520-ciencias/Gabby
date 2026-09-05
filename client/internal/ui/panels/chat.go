package panels

import (
	"client/internal/ui/styles"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ChatModel struct {
	Panel
	TextInput textinput.Model
	ChatName  string
}

func NewChatModel() ChatModel {
	ti := textinput.New()
	ti.Placeholder = "Message..."
	ti.Blur()
	ti.CharLimit = 156
	return ChatModel{
		TextInput: ti,
		ChatName:  "Global",
	}
}

func (m ChatModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m ChatModel) IsCapturingInput() bool {
	return m.TextInput.Focused()
}

func (m ChatModel) Update(msg tea.Msg) (ChatModel, tea.Cmd) {
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
			// text := m.TextInput.Value()
			m.TextInput.SetValue("")
			return m, nil
		}
	}

	var cmd tea.Cmd
	if m.TextInput.Focused() {
		m.TextInput, cmd = m.TextInput.Update(msg)
	}
	return m, cmd
}

func (m ChatModel) View() string {
	boxStyle := styles.BoxStyle.Width(m.width).Height(m.height)
	header := styles.ChatNameStyle.Width(m.width - 2).Render(m.ChatName)
	contentStyle := lipgloss.NewStyle().Width(m.width - 2).Height(m.height - 4)

	if m.IsFocused() {
		boxStyle = boxStyle.BorderForeground(styles.PrimaryColor)
	}

	messageBarStyle := styles.TextInputStyle.Width(m.width - 4) // -4 Por los bordes de el panel y del input.
	if m.IsCapturingInput() {
		messageBarStyle = messageBarStyle.BorderForeground(styles.AccentColor)
	}

	return boxStyle.Render(header + "\n" + contentStyle.Render() + "\n" + messageBarStyle.Render(m.TextInput.View()))
}
