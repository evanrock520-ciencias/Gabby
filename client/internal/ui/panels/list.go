package panels

import (
	"client/internal/ui/messages"
	"client/internal/ui/styles"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type ListModel struct {
	Panel
	Title   string
	Items   []string
	OnEnter func(value string) messages.ListResultMsg
}

func NewListModel(title string, items []string, onEnter func(value string) messages.ListResultMsg) ListModel {
	return ListModel{
		Title:   title,
		Items:   items,
		OnEnter: onEnter,
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
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.Items)-1 {
				m.cursor++
			}
		case "enter":
			if len(m.Items) == 0 || m.cursor < 0 || m.cursor >= len(m.Items) {
				return m, nil
			}
			return m, func() tea.Msg {
				return m.OnEnter(m.Items[m.cursor])
			}
		}
	}
	return m, nil
}

func (m ListModel) IsCapturingInput() bool {
	return false
}

func (m ListModel) View() string {
	boxStyle := styles.BoxStyle.Width(m.width).Height(m.height)

	if m.IsFocused() {
		boxStyle = boxStyle.BorderForeground(styles.PrimaryColor)
	}

	var sb strings.Builder
	for i, item := range m.Items {
		itemStyle := styles.TextStyle
		if m.IsFocused() && i == m.cursor {
			itemStyle = styles.SelectedStyle.Width(m.width - 2)
		}
		sb.WriteString(itemStyle.Render(item))
		sb.WriteString("\n")
	}

	innerWidth := max(0, m.width-2)
	headerText := styles.HeaderTitleStyle.Render(m.Title)
	divider := styles.DividerStyle.Render(strings.Repeat("─", innerWidth))

	return boxStyle.Render(headerText + "\n" + divider + "\n" + sb.String())
}

func (m ListModel) SelectedItem() (string, bool) {
	if m.cursor >= 0 && m.cursor < len(m.Items) {
		return m.Items[m.cursor], true
	}
	return "", false
}

func (m *ListModel) AddItem(item string) {
	m.Items = append(m.Items, item)
}

func (m *ListModel) RemoveItem(item string) {
	i := slices.Index(m.Items, item)
	if i == -1 {
		return
	}
	m.Items = slices.Delete(m.Items, i, i+1)
	if len(m.Items) == 0 {
		m.cursor = 0
	} else if m.cursor >= len(m.Items) {
		m.cursor = len(m.Items) - 1
	}
}
