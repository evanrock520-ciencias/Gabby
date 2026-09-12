package panels

import (
	"client/internal/ui/messages"
	"client/internal/ui/styles"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type InputModal struct {
	Layout
	TextInput textinput.Model
	Prompt    string
	onConfirm func(value string) messages.ModalResultMsg
	optional  bool
}

func NewInputModal(prompt string, placeholder string, charLimit int, action func(value string) messages.ModalResultMsg, optional bool) *InputModal {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Blur()
	ti.CharLimit = charLimit
	return &InputModal{
		Prompt:    prompt,
		TextInput: ti,
		onConfirm: action,
		optional:  optional,
	}
}

func (m *InputModal) Init() tea.Cmd {
	return nil
}

func (m *InputModal) IsCapturingInput() bool {
	return m.TextInput.Focused()
}

func (m *InputModal) Focus() tea.Cmd {
	return m.TextInput.Focus()
}

func (m *InputModal) Update(msg tea.Msg) (Modal, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if !m.optional {
				return m, tea.Quit
			}
			m.TextInput.Reset()
			m.TextInput.Blur()
			return m, nil
		case tea.KeyEnter:
			text := m.TextInput.Value()
			return m, func() tea.Msg {
				return m.onConfirm(text)
			}
		}
	}

	var cmd tea.Cmd
	m.TextInput, cmd = m.TextInput.Update(msg)
	return m, cmd
}

func (m *InputModal) View() string {
	innerWidth := 50
	if m.width > 0 && innerWidth > m.width-4 {
		innerWidth = max(20, m.width-4)
	}
	innerWidth = max(10, innerWidth-6)

	inputStyle := styles.TextInputStyle.Width(innerWidth - 2).BorderForeground(styles.PrimaryColor)
	inputBox := inputStyle.Render(m.TextInput.View())

	return renderModalCard(m.width, m.height, m.Prompt, inputBox)
}
