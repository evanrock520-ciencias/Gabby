package panels

import (
	"client/internal/ui/styles"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ChatModel struct {
	Width     int
	Height    int
	Focused   bool
	TextInput textinput.Model
}

func NewChatModel() ChatModel {
	ti := textinput.New()
	ti.Placeholder = "Message..."
	ti.Blur()
	ti.CharLimit = 156
	return ChatModel{
		TextInput: ti,
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
	boxStyle := styles.BoxStyle.Width(m.Width).Height(m.Height)

	if m.Focused {
		boxStyle = boxStyle.BorderForeground(styles.PrimaryColor)
	}

	messageBarStyle := styles.TextInputStyle.Width(m.Width - 4) // -4 Por los bordes de el panel y del input.
	if m.IsCapturingInput() {
		messageBarStyle = messageBarStyle.BorderForeground(styles.AccentColor)
	}

	return boxStyle.Render(messageBarStyle.Render(m.TextInput.View()))
}
