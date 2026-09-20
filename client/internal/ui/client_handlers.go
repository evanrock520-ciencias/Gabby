package ui

import (
	"client/internal/protocol"

	tea "github.com/charmbracelet/bubbletea"
)

// routeClientMessage enruta los mensajes del cliente al servidor.
func (m *Model) routeClientMessage(msg protocol.ClientMessage) tea.Cmd {
	switch msg.Type {
	case protocol.DISCONNECT:
		return tea.Sequence(m.sendToServer(msg), tea.Quit)

	case protocol.IDENTIFY:
		return m.handleIdentify(msg)

	case protocol.STATUS:
		return m.handleStatus(msg)

	case protocol.NEW_ROOM:
		return m.handleNewRoom(msg)

	case protocol.INVITE:
		return m.handleInvite(msg)

	case protocol.JOIN_ROOM:
		return m.handleJoinRoom(msg)
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

// handleInvite maneja el envío de invitaciones.
func (m *Model) handleInvite(msg protocol.ClientMessage) tea.Cmd {
	return m.sendToServer(msg)
}

// handleJoinRoom maneja la unión a salas tras una invitación.
func (m *Model) handleJoinRoom(msg protocol.ClientMessage) tea.Cmd {
	return m.sendToServer(msg)
}
