package ui

import "github.com/charmbracelet/bubbles/key"

var (
	modalKeymaps = []key.Binding{
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "Confirm")),
		key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "Cancel")),
	}

	chatInputKeymaps = []key.Binding{
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "Send")),
		key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "Stop Writing")),
	}

	commonNavKeymaps = []key.Binding{
		key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "Exit")),
		key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "Next Panel")),
		key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "Create Room")),
		key.NewBinding(key.WithKeys("i"), key.WithHelp("i", "Invite")),
		key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "Change Status")),
	}

	roomsKeymaps   []key.Binding
	usersKeymaps   []key.Binding
	chatNavKeymaps []key.Binding
)

func init() {
	roomsKeymaps = append(append([]key.Binding{}, commonNavKeymaps...),
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "Enter Room")),
	)

	usersKeymaps = append(append([]key.Binding{}, commonNavKeymaps...),
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "Open DM")),
	)

	chatNavKeymaps = append(append([]key.Binding{}, commonNavKeymaps...),
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "Write")),
		key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "Leave Room")),
		key.NewBinding(key.WithKeys("l"), key.WithHelp("l", "Room Users")),
	)
}

func (m Model) CurrentKeymaps() []key.Binding {
	// TODO: Embeber keybinds de paneles concretos, no solo la interfaz modal
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
