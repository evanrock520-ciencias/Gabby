package panels

import (
	"client/internal/domain"
	"client/internal/ui/styles"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

type FooterModel struct {
	Layout
	help     help.Model
	Username string
	Status   domain.Status
}

func NewFooterModel(username string, status domain.Status) FooterModel {
	h := help.New()
	h.ShortSeparator = "  "
	h.Styles.ShortKey = lipgloss.NewStyle().Foreground(styles.PrimaryColor).Bold(true)
	h.Styles.ShortDesc = styles.KeymapsStyle
	h.Styles.ShortSeparator = styles.KeymapsStyle
	return FooterModel{
		help:     h,
		Username: username,
		Status:   status,
	}
}

// ViewWithBindings renderiza el footer con los keybinds dados.
func (m FooterModel) ViewWithBindings(bindings []key.Binding) string {
	availableWidth := m.width - 2
	footerStyle := lipgloss.NewStyle().Width(availableWidth).MarginLeft(1).MarginRight(1).Height(3)

	m.help.Width = availableWidth
	left := m.help.ShortHelpView(bindings)
	leftWidth := lipgloss.Width(left)

	statusWidth := max(0, availableWidth-leftWidth)
	right := styles.StatusStyle.Width(statusWidth).Align(lipgloss.Right).Foreground(styles.ColorByStatus(m.Status)).Render("◉ " + m.Username)

	return footerStyle.AlignVertical(lipgloss.Center).Render(left + right)
}
