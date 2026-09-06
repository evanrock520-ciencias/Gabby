package protocol

import "errors"

type TypeC2S string

const (
	IDENTIFY    TypeC2S = "IDENTIFY"
	STATUS      TypeC2S = "STATUS"
	USERS       TypeC2S = "USERS"
	TEXT        TypeC2S = "TEXT"
	PUBLIC_TEXT TypeC2S = "PUBLIC_TEXT"
	NEW_ROOM    TypeC2S = "NEW_ROOM"
	INVITE      TypeC2S = "INVITE"
	JOIN_ROOM   TypeC2S = "JOIN_ROOM"
	ROOM_USERS  TypeC2S = "ROOM_USERS"
	ROOM_TEXT   TypeC2S = "ROOM_TEXT"
	LEAVE_ROOM  TypeC2S = "LEAVE_ROOM"
	DISCONNECT  TypeC2S = "DISCONNECT"
)

type ClientMessage struct {
	Type TypeC2S `json:"type,omitempty"`
	Message
}

func IdentifyMessage(username string) (ClientMessage, error) {
	if len(username) > 8 {
		return ClientMessage{}, errors.New("La longitud del nombre de usuario es mayor a 8 carácteres.")
	}
	return ClientMessage{Type: IDENTIFY, Message: Message{Username: username}}, nil
}

func StatusMessage(status Status) (ClientMessage, error) {
	return ClientMessage{Type: STATUS, Message: Message{Status: status}}, nil
}

func UsersMessage() (ClientMessage, error) {
	return ClientMessage{Type: USERS}, nil
}

func TextMessage(username string, text string) (ClientMessage, error) {
	if len(username) > 8 {
		return ClientMessage{}, errors.New("La longitud del nombre de usuario es mayor a 8 carácteres.")
	}
	return ClientMessage{Type: TEXT, Message: Message{Username: username, Text: text}}, nil
}

func PublicTextMessage(text string) (ClientMessage, error) {
	return ClientMessage{Type: PUBLIC_TEXT, Message: Message{Text: text}}, nil
}

func NewRoomMessage(roomname string) (ClientMessage, error) {
	if len(roomname) > 16 {
		return ClientMessage{}, errors.New("La longitud del nombre de la sala es mayor a 16 carácteres.")
	}
	return ClientMessage{Type: NEW_ROOM, Message: Message{Roomname: roomname}}, nil
}

func InviteMessage(roomname string, usernames []string) (ClientMessage, error) {
	if len(roomname) > 16 {
		return ClientMessage{}, errors.New("La longitud del nombre de la sala es mayor a 16 carácteres.")
	}
	return ClientMessage{Type: INVITE, Message: Message{Roomname: roomname, Usernames: usernames}}, nil
}

func JoinRoomMessage(roomname string) (ClientMessage, error) {
	if len(roomname) > 16 {
		return ClientMessage{}, errors.New("La longitud del nombre de la sala es mayor a 16 carácteres.")
	}
	return ClientMessage{Type: JOIN_ROOM, Message: Message{Roomname: roomname}}, nil
}

func RoomUsersMessage(roomname string) (ClientMessage, error) {
	if len(roomname) > 16 {
		return ClientMessage{}, errors.New("La longitud del nombre de la sala es mayor a 16 carácteres.")
	}
	return ClientMessage{Type: ROOM_USERS, Message: Message{Roomname: roomname}}, nil
}

func RoomTextMessage(roomname string, text string) (ClientMessage, error) {
	if len(roomname) > 16 {
		return ClientMessage{}, errors.New("La longitud del nombre de la sala es mayor a 16 carácteres.")
	}
	return ClientMessage{Type: ROOM_TEXT, Message: Message{Roomname: roomname, Text: text}}, nil
}

func LeaveRoomMessage(roomname string) (ClientMessage, error) {
	if len(roomname) > 16 {
		return ClientMessage{}, errors.New("La longitud del nombre de la sala es mayor a 16 carácteres.")
	}
	return ClientMessage{Type: LEAVE_ROOM, Message: Message{Roomname: roomname}}, nil
}

func DisconnectMessage() (ClientMessage, error) {
	return ClientMessage{Type: DISCONNECT}, nil
}
