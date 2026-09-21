package ui

import (
	"client/internal/domain"
	"client/internal/protocol"
	"client/internal/ui/messages"
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

	case protocol.ROOM_TEXT_FROM:
		return m.handleRoomTextFrom(msg)

	case protocol.JOINED_ROOM:
		return m.handleJoinedRoom(msg)

	case protocol.LEFT_ROOM:
		return m.handleLeftRoom(msg)

	case protocol.ROOM_USER_LIST:
		return m.handleRoomUserList(msg)

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

	global, _ := m.session.GetRoom("Global")
	m.chatModel.AddEntry(global, domain.ChatEvent{Text: fmt.Sprintf("%s has joined the chat", msg.Username)})
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

	m.chatModel.AddEntry(room, chatMsg)
	return nil
}

// handleTextFrom maneja la llegada de DMs.
func (m *Model) handleTextFrom(msg protocol.ServerMessage) tea.Cmd {
	chatMsg := domain.ChatMessage{Username: msg.Username, Message: msg.Text}

	dmRoom := m.session.GetOrCreateDM(msg.Username)
	m.chatModel.AddEntry(dmRoom, chatMsg)
	return nil
}

// handleRoomTextFrom maneja la llegada de mensajes a salas.
func (m *Model) handleRoomTextFrom(msg protocol.ServerMessage) tea.Cmd {
	chatMsg := domain.ChatMessage{Username: msg.Username, Message: msg.Text}

	room, _ := m.session.GetRoom(msg.Roomname)
	m.chatModel.AddEntry(room, chatMsg)
	return nil
}

// handleJoinedRoom maneja la notificación de un nuevo usuario en la sala.
func (m *Model) handleJoinedRoom(msg protocol.ServerMessage) tea.Cmd {
	room, _ := m.session.GetRoom(msg.Roomname)
	m.chatModel.AddEntry(room, domain.ChatEvent{Text: fmt.Sprintf("%s has joined the room", msg.Username)})
	return nil
}

// handleLeftRoom maneja la notificación de un usuario que dejó la sala.
func (m *Model) handleLeftRoom(msg protocol.ServerMessage) tea.Cmd {
	room, _ := m.session.GetRoom(msg.Roomname)
	m.chatModel.AddEntry(room, domain.ChatEvent{Text: fmt.Sprintf("%s has left the room", msg.Username)})
	return nil
}

// handleRoomUserList notifica la apertura de un modal en el chat.
func (m *Model) handleRoomUserList(msg protocol.ServerMessage) tea.Cmd {
	users := make([]domain.User, 0, len(msg.Users))
	for username, status := range msg.Users {
		users = append(users, domain.User{Username: username, Status: status})
	}

	return func() tea.Msg {
		return messages.ShowRoomUserlist{Roomname: msg.Roomname, Users: users}
	}
}

// handleDisconnected maneja la desconexión de algún usuario del servidor.
func (m *Model) handleDisconnected(msg protocol.ServerMessage) tea.Cmd {
	global, _ := m.session.GetRoom("Global")
	m.chatModel.AddEntry(global, domain.ChatEvent{Text: fmt.Sprintf("%s has left the chat", msg.Username)})

	m.session.RemoveUser(msg.Username)
	m.usersModel.RemoveItem(msg.Username)
	return nil
}

// routeResponse enruta los mensajes del servidor al cliente
// que sean especificamente una respuesta.
func (m *Model) routeResponse(msg protocol.ServerMessage) tea.Cmd {
	switch msg.Operation {
	case protocol.IDENTIFY:
		return m.handleIdentifyResponse(msg)

	case protocol.NEW_ROOM:
		return m.handleNewRoomResponse(msg)

	case protocol.JOIN_ROOM:
		return m.handleJoinRoomResponse(msg)

	case protocol.TEXT:
		return m.handleTextResponse(msg)

	case protocol.INVITE:
		return m.handleInviteResponse(msg)

	case protocol.ROOM_USERS:
		return m.handleRoomUserResponse(msg)

	case protocol.ROOM_TEXT:
		return m.handleRoomTextResponse(msg)

	case protocol.LEAVE_ROOM:
		return m.handleLeaveRoomResponse(msg)

	case protocol.INVALID_OPERATION:
		return m.handleInvalidResponse(msg)
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

	if msg.Result == protocol.USER_ALREADY_EXISTS {
		return func() tea.Msg {
			return messages.ShowNotification{Prompt: "An error has ocurred", Notification: fmt.Sprintf("The user %s already exists", msg.Extra), IsFatal: true}
		}
	}

	return nil
}

// handleNewRoomResponse Verifica que se creó una sala exitosamente.
func (m *Model) handleNewRoomResponse(msg protocol.ServerMessage) tea.Cmd {
	if msg.Result == protocol.ROOM_ALREADY_EXISTS {
		return func() tea.Msg {
			return messages.ShowNotification{Prompt: "An error has ocurred", Notification: fmt.Sprintf("The room %s already exists", msg.Extra), IsFatal: false}
		}
	}

	if msg.Result == protocol.SUCCESS {
		m.closeModal()

		room := m.session.AddRoom(domain.RoomChannel, msg.Extra)
		m.roomsModel.AddItem(NewRoomItem(room))
	}

	return nil
}

// handleJoinRoomResponse Verifica la unión a una sala
func (m *Model) handleJoinRoomResponse(msg protocol.ServerMessage) tea.Cmd {
	if msg.Result == protocol.NO_SUCH_ROOM {
		return func() tea.Msg {
			return messages.ShowNotification{Prompt: "An error has ocurred", Notification: fmt.Sprintf("The room %s does not exist", msg.Extra), IsFatal: false}
		}
	}

	if msg.Result == protocol.NOT_INVITED {
		return func() tea.Msg {
			return messages.ShowNotification{Prompt: "An error has ocurred", Notification: fmt.Sprintf("You were not invited to the room %s", msg.Extra), IsFatal: false}
		}
	}

	if msg.Result == protocol.SUCCESS {
		room := m.session.AddRoom(domain.RoomChannel, msg.Extra)
		m.roomsModel.AddItem(NewRoomItem(room))
	}

	m.closeModal()

	return nil
}

// handleTextResponse notifica usuario inválido al mandar DM.
func (m *Model) handleTextResponse(msg protocol.ServerMessage) tea.Cmd {
	if msg.Result == protocol.NO_SUCH_USER {
		return func() tea.Msg {
			return messages.ShowNotification{Prompt: "An error has ocurred", Notification: fmt.Sprintf("The user %s does not exist", msg.Extra), IsFatal: false}
		}
	}

	return nil
}

// handleInviteResponse notifica sala o usuario inválido al invitar a salas.
func (m *Model) handleInviteResponse(msg protocol.ServerMessage) tea.Cmd {
	if msg.Result == protocol.NO_SUCH_ROOM {
		return func() tea.Msg {
			return messages.ShowNotification{Prompt: "An error has ocurred", Notification: fmt.Sprintf("The room %s does not exist", msg.Extra), IsFatal: false}
		}
	}

	if msg.Result == protocol.NO_SUCH_USER {
		return func() tea.Msg {
			return messages.ShowNotification{Prompt: "An error has ocurred", Notification: fmt.Sprintf("The user %s does not exist", msg.Extra), IsFatal: false}
		}
	}

	return nil
}

// handleRoomUserResponse notifica sala inválida o acceso inválido al pedir la lista de usuarios.
func (m *Model) handleRoomUserResponse(msg protocol.ServerMessage) tea.Cmd {
	if msg.Result == protocol.NO_SUCH_ROOM {
		return func() tea.Msg {
			return messages.ShowNotification{Prompt: "An error has ocurred", Notification: fmt.Sprintf("The room %s does not exist", msg.Extra), IsFatal: false}
		}
	}

	if msg.Result == protocol.NOT_JOINED {
		return func() tea.Msg {
			return messages.ShowNotification{Prompt: "An error has ocurred", Notification: fmt.Sprintf("You are not a member of %s", msg.Extra), IsFatal: false}
		}
	}

	return nil
}

// handleRoomTextResponse notifica sala inválida o acceso inválido al pedir al mandar mensajes a salas.
func (m *Model) handleRoomTextResponse(msg protocol.ServerMessage) tea.Cmd {
	if msg.Result == protocol.NO_SUCH_ROOM {
		return func() tea.Msg {
			return messages.ShowNotification{Prompt: "An error has ocurred", Notification: fmt.Sprintf("The room %s does not exist", msg.Extra), IsFatal: false}
		}
	}

	if msg.Result == protocol.NOT_JOINED {
		return func() tea.Msg {
			return messages.ShowNotification{Prompt: "An error has ocurred", Notification: fmt.Sprintf("You are not a member of %s", msg.Extra), IsFatal: false}
		}
	}

	return nil
}

// handleLeaveRoomResponse notifica sala inválida o acceso inválido al salir de salas.
func (m *Model) handleLeaveRoomResponse(msg protocol.ServerMessage) tea.Cmd {
	if msg.Result == protocol.NO_SUCH_ROOM {
		return func() tea.Msg {
			return messages.ShowNotification{Prompt: "An error has ocurred", Notification: fmt.Sprintf("The room %s does not exist", msg.Extra), IsFatal: false}
		}
	}

	if msg.Result == protocol.NOT_JOINED {
		return func() tea.Msg {
			return messages.ShowNotification{Prompt: "An error has ocurred", Notification: fmt.Sprintf("You are not a member of %s", msg.Extra), IsFatal: false}
		}
	}

	return nil
}

// handle notifica mensajes inválidos y manda señal de terminar el programa.
func (m *Model) handleInvalidResponse(msg protocol.ServerMessage) tea.Cmd {
	if msg.Result == protocol.NOT_IDENTIFIED {
		return func() tea.Msg {
			return messages.ShowNotification{Prompt: "Fatal Error", Notification: "You are not identified with the server", IsFatal: true}
		}
	}

	return func() tea.Msg {
		return messages.ShowNotification{Prompt: "Fatal Error", Notification: "Invalid message received by the server", IsFatal: true}
	}
}
