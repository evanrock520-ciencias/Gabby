package panels

import (
	"client/internal/ui/styles"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Modal interface {
	InputCapturer
	Init() tea.Cmd
	Update(msg tea.Msg) (Modal, tea.Cmd)
	View() string
	Focus() tea.Cmd
	SetSize(width, height int)
}

func renderModalCard(width, height int, title string, body string) string {
	cardWidth := 50
	if width > 0 && cardWidth > width-4 {
		cardWidth = max(20, width-4)
	}

	innerWidth := max(10, cardWidth-6)

	header := styles.HeaderTitleStyle.Width(innerWidth).Align(lipgloss.Center).Render(title)
	divider := styles.DividerStyle.Render(strings.Repeat("─", innerWidth))

	cardBody := lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		divider,
		"",
		body,
	)

	cardStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(styles.PrimaryColor).Padding(1, 2).Width(cardWidth)
	card := cardStyle.Render(cardBody)

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		card,
	)
}
