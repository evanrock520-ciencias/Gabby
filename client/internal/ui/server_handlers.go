package ui

import (
	"client/internal/domain"
	"client/internal/protocol"

	tea "github.com/charmbracelet/bubbletea"
)

// routeServerMessage enruta los mensajes del servidor al cliente.
func (m *Model) routeServerMessage(msg protocol.ServerMessage) tea.Cmd {
	switch msg.Type {
	case protocol.RESPONSE:
		return m.routeResponse(msg)
	case protocol.USER_LIST:
		return m.handleUserList(msg)
	case protocol.NEW_USER:
		return m.handleNewUser(msg)
	}

	return nil
}

// handleUserList maneja la lógica de obtener la lista de usuarios.
func (m *Model) handleUserList(msg protocol.ServerMessage) tea.Cmd {
	for username, status := range msg.Users {
		user := domain.User{
			Username: username,
			Status:   status,
		}
		m.session.Users = append(m.session.Users, user)
		m.usersModel.AddItem(NewUserItem(user))
	}

	return nil
}

// handleNewUser maneja la notificación de un nuevo usuario en el servidor.
func (m *Model) handleNewUser(msg protocol.ServerMessage) tea.Cmd {
	user := domain.User{Username: msg.Username, Status: domain.ACTIVE}

	m.session.Users = append(m.session.Users, user)
	m.usersModel.AddItem(NewUserItem(user))

	return nil
}

// routeResponse enruta los mensajes del servidor al cliente
// que sean especificamente una respuesta.
func (m *Model) routeResponse(msg protocol.ServerMessage) tea.Cmd {
	switch msg.Operation {
	case protocol.IDENTIFY:
		m.handleIdentifyResponse(msg)

	}
	return nil
}

// handleIdentifyResponse Verifica que hubo una identificación exitosa.
func (m *Model) handleIdentifyResponse(msg protocol.ServerMessage) tea.Cmd {
	if msg.Result == protocol.SUCCESS {
		m.closeModal()
		m.SetUsername(msg.Extra)
		return nil
	}

	return nil
}
