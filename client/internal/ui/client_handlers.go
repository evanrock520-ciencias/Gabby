package ui

import (
	"client/internal/protocol"

	tea "github.com/charmbracelet/bubbletea"
)

// routeClientMessage enruta los mensajes del cliente al servidor.
func (m *Model) routeClientMessage(msg protocol.ClientMessage) tea.Cmd {
	switch msg.Type {
	case protocol.IDENTIFY:
		return m.handleIdentify(msg)

	case protocol.STATUS:
		return m.handleStatus(msg)

	case protocol.NEW_ROOM:
		return m.handleNewRoom(msg)
	}
	return nil
}

// handleIdentify maneja el envío de identificación.
func (m *Model) handleIdentify(msg protocol.ClientMessage) tea.Cmd {
	return m.sendToServer(msg)
}

// handleStatus maneja el envío de un nuevo estado.
func (m *Model) handleStatus(msg protocol.ClientMessage) tea.Cmd {
	m.SetStatus(msg.Status)
	return m.sendToServer(msg)
}

// handleNewRoom maneja el envío de una nueva sala.
func (m *Model) handleNewRoom(msg protocol.ClientMessage) tea.Cmd {
	return m.sendToServer(msg)
}
