package panels

import (
	"client/internal/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
)

type ChatModel struct {
	Width   int
	Height  int
	Focused bool
}

func (m ChatModel) Init() tea.Cmd {
	return nil
}

func (m ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m ChatModel) View() string {
	boxStyle := styles.BoxStyle.Width(m.Width).Height(m.Height)

	if m.Focused {
		boxStyle = boxStyle.BorderForeground(styles.PrimaryColor)
	}

	return boxStyle.Render()
}
