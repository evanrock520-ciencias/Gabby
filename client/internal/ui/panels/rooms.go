package panels

import (
	"client/internal/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
)

type RoomsModel struct {
	Width   int
	Height  int
	Focused bool
}

func (m RoomsModel) Init() tea.Cmd {
	return nil
}

func (m RoomsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m RoomsModel) View() string {
	boxStyle := styles.BoxStyle.Width(m.Width).Height(m.Height)

	if m.Focused {
		boxStyle = boxStyle.BorderForeground(styles.PrimaryColor)
	}

	headerText := styles.HeaderTitleStyle.Render("[1] Rooms")

	return boxStyle.Render(headerText)
}
