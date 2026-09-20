package ui

import (
	"client/internal/domain"
	"client/internal/protocol"
	"client/internal/ui/panels"
	"fmt"

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

	case protocol.NEW_STATUS:
		return m.handleNewStatus(msg)

	case protocol.INVITATION:
		return m.handleInvitation(msg)

	case protocol.PUBLIC_TEXT_FROM:
		return m.handlePublicTextFrom(msg)

	case protocol.TEXT_FROM:
		return m.handleTextFrom(msg)

	case protocol.DISCONNECTED:
		return m.handleDisconnected(msg)
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

// handleNewStatus maneja la notificación de un nuevo estado de usuario en el servidor.
func (m *Model) handleNewStatus(msg protocol.ServerMessage) tea.Cmd {
	m.session.SetUserStatus(msg.Username, msg.Status)
	m.usersModel.UpdateItem(msg.Username, NewUserItem(domain.User{Username: msg.Username, Status: msg.Status}))
	return nil
}

// handleInvitation maneja la llegada de una invitación a una sala.
func (m *Model) handleInvitation(msg protocol.ServerMessage) tea.Cmd {
	return m.openModal(panels.NewConfirmModal(fmt.Sprintf("Join the room %s", msg.Roomname), []string{"Join", "Reject"}, func() tea.Msg {
		join, _ := protocol.JoinRoomMessage(msg.Roomname)
		return join
	}, func() tea.Msg {
		return nil
	}))
}

// handlePublicTextFrom maneja la llegada de mensajes globales.
func (m *Model) handlePublicTextFrom(msg protocol.ServerMessage) tea.Cmd {
	chatMsg := domain.ChatMessage{Username: msg.Username, Message: msg.Text}
	room, _ := m.session.GetRoom("Global")

	m.chatModel.AddMessage(room, chatMsg)
	return nil
}

// handleTextFrom maneja la llegada de DMs.
func (m *Model) handleTextFrom(msg protocol.ServerMessage) tea.Cmd {
	chatMsg := domain.ChatMessage{Username: msg.Username, Message: msg.Text}

	dmRoom := m.session.GetOrCreateDM(msg.Username)
	m.chatModel.AddMessage(dmRoom, chatMsg)
	return nil
}

// handleDisconnected maneja la desconexión de algún usuario del servidor.
func (m *Model) handleDisconnected(msg protocol.ServerMessage) tea.Cmd {
	// TODO: Muestra notificación al chat global.
	m.session.RemoveUser(msg.Username)
	m.usersModel.RemoveItem(msg.Username)
	return nil
}

// routeResponse enruta los mensajes del servidor al cliente
// que sean especificamente una respuesta.
func (m *Model) routeResponse(msg protocol.ServerMessage) tea.Cmd {
	switch msg.Operation {
	case protocol.IDENTIFY:
		m.handleIdentifyResponse(msg)

	case protocol.NEW_ROOM:
		m.handleNewRoomResponse(msg)

	case protocol.JOIN_ROOM:
		m.handleJoinRoomResponse(msg)
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

// handleNewRoomResponse Verifica que se creó una sala exitosamente.
func (m *Model) handleNewRoomResponse(msg protocol.ServerMessage) tea.Cmd {
	if msg.Result == protocol.SUCCESS {
		m.closeModal()

		room := m.session.AddRoom(domain.RoomChannel, msg.Extra)
		m.roomsModel.AddItem(NewRoomItem(room))
	}

	// TODO: Quizás añadir notificaciones modales.
	return nil
}

// handleJoinRoomResponse Verifica la unión a una sala
func (m *Model) handleJoinRoomResponse(msg protocol.ServerMessage) tea.Cmd {
	if msg.Result == protocol.SUCCESS {
		room := m.session.AddRoom(domain.RoomChannel, msg.Extra)
		m.roomsModel.AddItem(NewRoomItem(room))
	}
	m.closeModal()
	return nil
}
