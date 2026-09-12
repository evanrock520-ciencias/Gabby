package panels

import (
	"client/internal/ui/messages"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type ListModal struct {
	Layout
	title    string
	items    []Item
	cursor   int
	offset   int
	multi    bool
	selected map[string]struct{}
	onDone   func(values []string) messages.ModalResultMsg
	active   bool
}

func NewListModal(title string, items []Item, multi bool, done func(values []string) messages.ModalResultMsg) *ListModal {
	return &ListModal{
		title:    title,
		items:    items,
		multi:    multi,
		onDone:   done,
		selected: make(map[string]struct{}),
		active:   true,
	}
}

func (m *ListModal) Init() tea.Cmd {
	return nil
}

func (m *ListModal) IsCapturingInput() bool {
	return m.active
}

func (m *ListModal) Focus() tea.Cmd {
	return nil
}

func (m *ListModal) Update(msg tea.Msg) (Modal, tea.Cmd) {
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
			if m.cursor < len(m.items)-1 {
				m.cursor++
				m.scrollToCursor()
			}
		case " ":
			if len(m.items) == 0 || m.cursor < 0 || m.cursor >= len(m.items) {
				return m, nil
			}
			value := m.items[m.cursor].Value()
			// log.Printf("Value selected %s", value)
			if m.multi {
				if _, ok := m.selected[value]; ok {
					delete(m.selected, value)
				} else {
					m.selected[value] = struct{}{}
				}
			} else {
				m.active = false
				return m, func() tea.Msg { return m.onDone([]string{value}) }
			}
		case "enter":
			if m.multi {
				m.active = false
				values := make([]string, 0, len(m.selected))
				for v := range m.selected {
					values = append(values, v)
				}
				return m, func() tea.Msg { return m.onDone(values) }
			}
		case "esc":
			m.active = false
		}
	}
	return m, nil
}

func (m *ListModal) View() string {
	cardWidth := 50
	if m.width > 0 && cardWidth > m.width-4 {
		cardWidth = max(20, m.width-4)
	}
	innerWidth := max(10, cardWidth-6)

	start, end := m.computeRanges()
	lines := make([]string, 0, end-start)
	for i := start; i < end; i++ {
		state := m.currentState(i)
		lines = append(lines, m.items[i].Render(state, innerWidth))
	}
	body := strings.Join(lines, "\n")

	return renderModalCard(m.width, m.height, m.title, body)
}

func (m ListModal) SelectedItem() (string, bool) {
	if m.cursor >= 0 && m.cursor < len(m.items) {
		return m.items[m.cursor].Value(), true
	}
	return "", false
}

func (m ListModal) visibleItems() int {
	visible := m.height - 2
	if visible < 1 {
		return 1
	}

	return visible
}

func (m *ListModal) scrollToCursor() {
	visible := m.visibleItems()

	if m.cursor < m.offset {
		m.offset = m.cursor
	}

	if m.cursor >= m.offset+visible {
		m.offset = m.cursor - visible + 1
	}

	maxOffset := max(0, len(m.items)-visible)
	if m.offset > maxOffset {
		m.offset = maxOffset
	}

	if m.offset < 0 {
		m.offset = 0
	}
}

func (m *ListModal) computeRanges() (int, int) {
	visible := m.visibleItems()
	end := min(len(m.items), m.offset+visible)
	start := min(m.offset, end)
	return start, end
}

func (m *ListModal) currentState(i int) ItemState {
	value := m.items[i].Value()
	_, isChosen := m.selected[value]
	isCursor := i == m.cursor
	var state ItemState
	switch {
	case isCursor && isChosen:
		state = ItemBoth
	case isChosen:
		state = ItemChosen
	case isCursor:
		state = ItemSelected
	default:
		state = ItemNormal
	}

	return state
}
