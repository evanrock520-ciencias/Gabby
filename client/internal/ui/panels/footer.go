package panels

import (
	"client/internal/protocol"
	"client/internal/ui/styles"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FooterModel struct {
	Layout
	Keymaps map[string]string
	// Probablemente convenga utilizar una referencia a un tipo usuario para actualizar la interfaz
	Username string
	Status   protocol.Status
}

func NewFooterModel(keymaps map[string]string, username string, status protocol.Status) FooterModel {
	return FooterModel{
		Keymaps:  keymaps,
		Username: username,
		Status:   status,
	}
}

func (m FooterModel) Init() tea.Cmd {
	return nil
}

func (m FooterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m FooterModel) View() string {
	availableWidth := m.width - 2
	footerStyle := lipgloss.NewStyle().Width(availableWidth).MarginLeft(1).MarginRight(1).Height(3)

	var sb strings.Builder
	for key, action := range m.Keymaps {
		sb.WriteString(styles.KeymapsStyle.Render("● " + key + " - " + action))
		sb.WriteString("    ")
	}

	left := sb.String()
	leftWidth := lipgloss.Width(left)

	statusWidth := max(0, availableWidth-leftWidth)
	right := styles.StatusStyle.Width(statusWidth).Align(lipgloss.Right).Foreground(defineColorByStatus(m.Status)).Render("◉ " + m.Username)

	return footerStyle.AlignVertical(lipgloss.Center).Render(left + right)
}

func defineColorByStatus(status protocol.Status) lipgloss.Color {
	switch status {
	case protocol.ACTIVE:
		return styles.PrimaryColor
	case protocol.BUSY:
		return styles.TerciaryColor
	case protocol.AWAY:
		return styles.SecondaryColor
	}
	return styles.PrimaryColor
}
