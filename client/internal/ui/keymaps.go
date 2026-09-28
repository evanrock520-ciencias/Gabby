package ui

import "github.com/charmbracelet/bubbles/key"

var (
	modalKeymaps = []key.Binding{
		key.NewBinding(key.WithKeys("↵"), key.WithHelp("↵", "Confirm")),
		key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "Cancel")),
		key.NewBinding(key.WithKeys("space"), key.WithHelp("␣", "Select")),
	}

	chatInputKeymaps = []key.Binding{
		key.NewBinding(key.WithKeys("↵"), key.WithHelp("↵", "Send")),
		key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "Stop Writing")),
	}

	commonNavKeymaps = []key.Binding{
		key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "Exit")),
		key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "New")),
		key.NewBinding(key.WithKeys("i"), key.WithHelp("i", "Invite")),
		key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "Status")),
	}

	roomsKeymaps   []key.Binding
	usersKeymaps   []key.Binding
	chatNavKeymaps []key.Binding
)

func init() {
	roomsKeymaps = append(append([]key.Binding{}, commonNavKeymaps...),
		key.NewBinding(key.WithKeys("↵"), key.WithHelp("↵", "Enter")),
	)

	usersKeymaps = append(append([]key.Binding{}, commonNavKeymaps...),
		key.NewBinding(key.WithKeys("↵"), key.WithHelp("↵", "Enter")),
	)

	chatNavKeymaps = append(append([]key.Binding{}, commonNavKeymaps...),
		key.NewBinding(key.WithKeys("↵"), key.WithHelp("↵", "Write")),
		key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "Leave")),
		key.NewBinding(key.WithKeys("l"), key.WithHelp("l", "Users")),
	)
}

func (m Model) CurrentKeymaps() []key.Binding {
	if m.activeModal != nil {
		return modalKeymaps
	}

	if active := m.activePanel(); active != nil && active.IsCapturingInput() {
		return chatInputKeymaps
	}

	switch m.focus {
	case Rooms:
		return roomsKeymaps
	case Users:
		return usersKeymaps
	case Chat:
		return chatNavKeymaps
	default:
		return roomsKeymaps
	}
}
