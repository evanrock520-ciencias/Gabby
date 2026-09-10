package panels

import (
	"client/internal/domain"
	"client/internal/ui/styles"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Keymap struct {
	key    string
	action string
}

func NewKeymap(key string, action string) Keymap {
	return Keymap{
		key:    key,
		action: action,
	}
}

type FooterModel struct {
	Layout
	Keymaps []Keymap
	// Probablemente convenga utilizar una referencia a un tipo usuario para actualizar la interfaz
	Username string
	Status   domain.Status
}

func NewFooterModel(keymaps []Keymap, username string, status domain.Status) FooterModel {
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
	for _, keymap := range m.Keymaps {
		sb.WriteString(styles.KeymapsStyle.Render("● " + keymap.key + " - " + keymap.action))
		sb.WriteString("    ")
	}

	left := sb.String()
	leftWidth := lipgloss.Width(left)

	statusWidth := max(0, availableWidth-leftWidth)
	right := styles.StatusStyle.Width(statusWidth).Align(lipgloss.Right).Foreground(defineColorByStatus(m.Status)).Render("◉ " + m.Username)

	return footerStyle.AlignVertical(lipgloss.Center).Render(left + right)
}

func (m *FooterModel) SetKeys(keys []Keymap) {
	m.Keymaps = keys
}

func defineColorByStatus(status domain.Status) lipgloss.Color {
	switch status {
	case domain.ACTIVE:
		return styles.PrimaryColor
	case domain.BUSY:
		return styles.TerciaryColor
	case domain.AWAY:
		return styles.SecondaryColor
	}
	return styles.PrimaryColor
}
