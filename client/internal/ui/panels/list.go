package panels

import (
	"client/internal/ui/styles"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type ListModel struct {
	Title   string
	Width   int
	Height  int
	Focused bool
	Cursor  int
	Items   []string
}

func NewListModel(title string, items []string) ListModel {
	return ListModel{
		Title: title,
		Items: items,
	}
}

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m ListModel) Update(msg tea.Msg) (ListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down":
			if m.Cursor < len(m.Items)-1 {
				m.Cursor++
			}
		}
	}
	return m, nil
}

func (m ListModel) View() string {
	boxStyle := styles.BoxStyle.Width(m.Width).Height(m.Height)

	if m.Focused {
		boxStyle = boxStyle.BorderForeground(styles.PrimaryColor)
	}

	var sb strings.Builder
	for i, item := range m.Items {
		itemStyle := styles.TextStyle
		if m.Focused && i == m.Cursor {
			itemStyle = styles.SelectedStyle
		}
		sb.WriteString(itemStyle.Render(item))
		sb.WriteString("\n")
	}

	innerWidth := max(0, m.Width-2)
	headerText := styles.HeaderTitleStyle.Render(m.Title)
	divider := styles.DividerStyle.Render(strings.Repeat("─", innerWidth))

	return boxStyle.Render(headerText + "\n" + divider + "\n" + sb.String())
}

func (m ListModel) SelectedItem() (string, bool) {
	if m.Cursor >= 0 && m.Cursor < len(m.Items) {
		return m.Items[m.Cursor], true
	}
	return "", false
}
