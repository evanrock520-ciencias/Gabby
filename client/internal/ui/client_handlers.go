package ui

import (
	"client/internal/protocol"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) routeClientMessage(msg protocol.ClientMessage) tea.Cmd {
	switch msg.Type {
	case protocol.IDENTIFY:
		return m.handleIdentify(msg)

	case protocol.STATUS:
		return m.handleStatus(msg)
	}
	return nil
}

func (m *Model) handleIdentify(msg protocol.ClientMessage) tea.Cmd {
	return m.sendToServer(msg)
}

func (m *Model) handleStatus(msg protocol.ClientMessage) tea.Cmd {
	m.SetStatus(msg.Status)
	return m.sendToServer(msg)
}
