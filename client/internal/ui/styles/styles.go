package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	BoxStyle         = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1)
	TextStyle        = lipgloss.NewStyle().Foreground(TextColor)
	SelectedStyle    = lipgloss.NewStyle().Foreground(TextColor).Background(AccentColor)
	HeaderTitleStyle = lipgloss.NewStyle().Foreground(PrimaryColor).Bold(true)
	DividerStyle     = lipgloss.NewStyle().Foreground(MutedColor)
	TextInputStyle   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
	ChatNameStyle    = lipgloss.NewStyle().Foreground(TextColor).Background(AccentColor).Align(lipgloss.Center)
	StatusStyle      = lipgloss.NewStyle().Bold(true)
	KeymapsStyle     = lipgloss.NewStyle().Foreground(MutedColor)
	MessageStyle     = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, false, true).MarginTop(1)
	//TODO: Añadir estilo para destacar mensajes propios
)
