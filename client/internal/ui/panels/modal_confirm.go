package panels

import (
	"client/internal/ui/messages"
	"client/internal/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ConfirmModal struct {
	Layout
	Prompt   string
	Choices  []string
	OnAccept func() messages.ModalResultMsg
	OnCancel func() messages.ModalResultMsg

	cursor int
	active bool
}

func NewConfirmModal(prompt string, choices []string, accept func() messages.ModalResultMsg, cancel func() messages.ModalResultMsg) *ConfirmModal {
	return &ConfirmModal{
		Prompt:   prompt,
		Choices:  choices,
		OnAccept: accept,
		OnCancel: cancel,
		active:   true,
	}
}

func (m *ConfirmModal) Init() tea.Cmd {
	return nil
}

func (m *ConfirmModal) IsCapturingInput() bool {
	return m.active
}

func (m *ConfirmModal) Focus() tea.Cmd {
	return nil
}

func (m *ConfirmModal) Update(msg tea.Msg) (Modal, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left":
			if m.cursor > 0 {
				m.cursor--
			}
		case "right":
			if m.cursor < len(m.Choices)-1 {
				m.cursor++
			}
		case "enter":
			m.active = false
			if m.cursor == 0 {
				return m, func() tea.Msg { return m.OnAccept() }
			}
			return m, func() tea.Msg { return m.OnCancel() }
		case "esc":
			m.active = false
			return m, nil
		}
	}
	return m, nil
}

func (m *ConfirmModal) View() string {
	acceptButton := styles.NormalButtonStyle.Render(m.Choices[0])
	if m.cursor == 0 {
		acceptButton = styles.AcceptButtonStyle.Render(m.Choices[0])
	}

	cancelButton := styles.NormalButtonStyle.Render(m.Choices[1])
	if m.cursor == 1 {
		cancelButton = styles.CancelButtonStyle.Render(m.Choices[1])
	}
	body := lipgloss.JoinHorizontal(lipgloss.Center, acceptButton, "    ", cancelButton)

	return renderModalCard(m.width, m.height, m.Prompt, body)
}
