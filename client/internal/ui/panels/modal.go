package panels

import (
	"client/internal/ui/styles"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ModalModel struct {
	Panel
	TextInput textinput.Model
	Prompt    string
	// Quizás convendria tambien tener una lista para selección
}

func NewModalModel(prompt string, placeholder string, charLimit int) ModalModel {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Blur()
	ti.CharLimit = charLimit
	return ModalModel{
		Prompt:    prompt,
		TextInput: ti,
	}
}

func (m ModalModel) Init() tea.Cmd {
	return nil
}

func (m ModalModel) IsCapturingInput() bool {
	return m.TextInput.Focused()
}

func (m ModalModel) Update(msg tea.Msg) (ModalModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		// TODO: Definir que hacer con los valores
		case tea.KeyEsc, tea.KeyEnter:
			m.TextInput.Reset()
			m.TextInput.Blur()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.TextInput, cmd = m.TextInput.Update(msg)
	return m, cmd
}

func (m ModalModel) View() string {
	cardWidth := 50
	if m.width > 0 && cardWidth > m.width-4 {
		cardWidth = max(20, m.width-4)
	}

	innerWidth := max(10, cardWidth-6)

	header := styles.HeaderTitleStyle.Width(innerWidth).Align(lipgloss.Center).Render(m.Prompt)
	divider := styles.DividerStyle.Render(strings.Repeat("─", innerWidth))
	inputStyle := styles.TextInputStyle.Width(innerWidth - 2).BorderForeground(styles.PrimaryColor)
	inputBox := inputStyle.Render(m.TextInput.View())

	cardBody := lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		divider,
		"",
		inputBox,
	)

	cardStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(styles.PrimaryColor).Padding(1, 2).Width(cardWidth)
	card := cardStyle.Render(cardBody)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center, // Lugar en el que se coloca horizontalmente
		lipgloss.Center, // Lugar en el que se coloca verticalmente
		card,
	)
}
