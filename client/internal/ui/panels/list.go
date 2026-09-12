package panels

import (
	"client/internal/ui/messages"
	"client/internal/ui/styles"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Item interface {
	Value() string
	Render(selected bool, width int) string
}

type ListModel struct {
	Panel
	Title   string
	Items   []Item
	OnEnter func(value string) messages.ListResultMsg

	offset int
}

func NewListModel(title string, items []Item, onEnter func(value string) messages.ListResultMsg) ListModel {
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

	case tea.WindowSizeMsg:
		m.scrollToCursor()

	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.cursor > 0 {
				m.cursor--
				m.scrollToCursor()
			}
		case "down":
			if m.cursor < len(m.Items)-1 {
				m.cursor++
				m.scrollToCursor()
			}
		case "enter":
			if len(m.Items) == 0 || m.cursor < 0 || m.cursor >= len(m.Items) {
				return m, nil
			}
			return m, func() tea.Msg {
				return m.OnEnter(m.Items[m.cursor].Value())
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

	innerWidth := max(0, m.width-2)
	headerText := styles.HeaderTitleStyle.Render(m.Title)
	divider := styles.DividerStyle.Render(strings.Repeat("─", innerWidth))
	lines := []string{headerText, divider}

	start, end := m.computeRanges()

	for i := start; i < end; i++ {
		selected := i == m.cursor && m.focus
		lines = append(lines, m.Items[i].Render(selected, m.width-2))

	}

	return boxStyle.Render(strings.Join(lines, "\n"))
}

func (m ListModel) SelectedItem() (string, bool) {
	if m.cursor >= 0 && m.cursor < len(m.Items) {
		return m.Items[m.cursor].Value(), true
	}
	return "", false
}

func (m *ListModel) AddItem(item Item) {
	m.Items = append(m.Items, item)
}

func (m *ListModel) RemoveItem(value string) {
	i := slices.IndexFunc(m.Items, func(item Item) bool {
		return item.Value() == value
	})
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

func (m ListModel) visibleItems() int {
	visible := m.height - 2
	if visible < 1 {
		return 1
	}

	return visible
}

func (m *ListModel) scrollToCursor() {
	visible := m.visibleItems()

	if m.cursor < m.offset {
		m.offset = m.cursor
	}

	if m.cursor >= m.offset+visible {
		m.offset = m.cursor - visible + 1
	}

	maxOffset := max(0, len(m.Items)-visible)
	if m.offset > maxOffset {
		m.offset = maxOffset
	}

	if m.offset < 0 {
		m.offset = 0
	}
}

func (m *ListModel) computeRanges() (int, int) {
	visible := m.visibleItems()
	end := min(len(m.Items), m.offset+visible)
	start := min(m.offset, end)
	return start, end
}
