package ui

import "client/internal/ui/panels"

var (
	modalKeymaps = []panels.Keymap{
		panels.NewKeymap("enter", "Confirm"),
		panels.NewKeymap("esc", "Cancel"),
	}

	chatInputKeymaps = []panels.Keymap{
		panels.NewKeymap("enter", "Send"),
		panels.NewKeymap("esc", "Stop Writing"),
	}

	commonNavKeymaps = []panels.Keymap{
		panels.NewKeymap("q", "Exit"),
		panels.NewKeymap("tab", "Next Panel"),
		panels.NewKeymap("c", "Create Room"),
		panels.NewKeymap("i", "Invite"),
		panels.NewKeymap("s", "Change Status"),
	}

	roomsKeymaps   []panels.Keymap
	usersKeymaps   []panels.Keymap
	chatNavKeymaps []panels.Keymap
)

func init() {
	roomsKeymaps = append(append([]panels.Keymap{}, commonNavKeymaps...),
		panels.NewKeymap("enter", "Enter Room"),
	)

	usersKeymaps = append(append([]panels.Keymap{}, commonNavKeymaps...),
		panels.NewKeymap("enter", "Open DM"),
	)

	chatNavKeymaps = append(append([]panels.Keymap{}, commonNavKeymaps...),
		panels.NewKeymap("enter", "Write"),
		panels.NewKeymap("e", "Leave Room"),
	)
}

func (m Model) CurrentKeymaps() []panels.Keymap {
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
