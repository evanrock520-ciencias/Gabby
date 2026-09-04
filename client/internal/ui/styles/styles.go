package styles

import "github.com/charmbracelet/lipgloss"

var (
	BoxStyle         = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1)
	TextStyle        = lipgloss.NewStyle().Foreground(TextColor)
	SelectedStyle    = lipgloss.NewStyle().Foreground(TextColor).Background(AccentColor)
	HeaderTitleStyle = lipgloss.NewStyle().Foreground(PrimaryColor).Bold(true)
)
