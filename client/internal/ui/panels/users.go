package panels

import (
	"client/internal/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
)

type UsersModel struct {
	Width   int
	Height  int
	Focused bool
}

func (m UsersModel) Init() tea.Cmd {
	return nil
}

func (m UsersModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m UsersModel) View() string {
	boxStyle := styles.BoxStyle.Width(m.Width).Height(m.Height)

	if m.Focused {
		boxStyle = boxStyle.BorderForeground(styles.PrimaryColor)
	}

	headerText := styles.HeaderTitleStyle.Render("[2] Users")

	return boxStyle.Render(headerText)
}
