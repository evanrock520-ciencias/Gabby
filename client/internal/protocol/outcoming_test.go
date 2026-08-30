package protocol

import (
	"fmt"
	"testing"
)

func TestIdentifyMessage(t *testing.T) {
	name := "Dale"
	msg, err := IdentifyMessage(name)
	if err != nil {
		fmt.Println(err)
	}
	if msg.Type != IDENTIFY {
		t.Errorf("Se esperaba el tipo %s, pero se obtuvo %s", IDENTIFY, msg.Type)
	}
	if msg.Username != "Dale" {
		t.Errorf("Se esperaba el nombre %s, pero se obtuvo %s", name, msg.Username)
	}
}

func TestIdentifyMessageError(t *testing.T) {
	name := "Laura Palmer"
	_, err := IdentifyMessage(name)
	if err == nil {
		t.Errorf("Se esperaba un error por longitud de nombre de usuario mayor a 8 carácteres")
	}
}

func TestStatusMessage(t *testing.T) {
	status := ACTIVE
	msg, err := StatusMessage(status)
	if err != nil {
		fmt.Println(err)
	}
	if msg.Type != STATUS {
		t.Errorf("Se esperaba el tipo %s, pero se obtuvo %s", STATUS, msg.Type)
	}
	if msg.Status != status {
		t.Errorf("Se esperaba el status %s, pero se obtuvo %s", status, msg.Status)
	}
}

func TestUsersMessage(t *testing.T) {
	msg, err := UsersMessage()
	if err != nil {
		fmt.Println(err)
	}
	if msg.Type != USERS {
		t.Errorf("Se esperaba el tipo %s, pero se obtuvo %s", USERS, msg.Type)
	}
}

func TestTextMessage(t *testing.T) {
	name := "Cooper"
	text := "Good damn cup of coffee"
	msg, err := TextMessage(name, text)
	if err != nil {
		fmt.Println(err)
	}
	if msg.Type != TEXT {
		t.Errorf("Se esperaba el tipo %s, pero se obtuvo %s", TEXT, msg.Type)
	}
	if msg.Username != name {
		t.Errorf("Se esperaba el nombre %s, pero se nombre %s", name, msg.Username)
	}
	if msg.Text != text {
		t.Errorf("Se esperaba el texto %s, pero se obtuvo %s", text, msg.Text)
	}
}

func TestTextMessageError(t *testing.T) {
	name := "Edward Louis Severson III"
	text := "Hola"
	_, err := TextMessage(name, text)
	if err == nil {
		t.Errorf("Se esperaba un error por longitud de nombre de usuario mayor a 8 carácteres")
	}
}

func TestPublicTextMessage(t *testing.T) {
	text := "Fire Walk With Me"
	msg, err := PublicTextMessage(text)
	if err != nil {
		fmt.Println(err)
	}
	if msg.Type != PUBLIC_TEXT {
		t.Errorf("Se esperaba el tipo %s, pero se obtuvo %s", PUBLIC_TEXT, msg.Type)
	}
	if msg.Text != text {
		t.Errorf("Se esperaba el texto %s, pero se obtuvo %s", text, msg.Text)
	}
}

func TestNewRoomMessage(t *testing.T) {
	roomname := "La Logia"
	msg, err := NewRoomMessage(roomname)
	if err != nil {
		fmt.Println(err)
	}
	if msg.Type != NEW_ROOM {
		t.Errorf("Se esperaba el tipo %s, pero se obtuvo %s", NEW_ROOM, msg.Type)
	}
	if msg.Roomname != roomname {
		t.Errorf("Se esperaba el nombre de sala %s, pero se obtuvo %s", roomname, msg.Roomname)
	}
}

func TestNewRoomMessageError(t *testing.T) {
	roomname := "La corte de los búhos"
	_, err := NewRoomMessage(roomname)
	if err == nil {
		t.Errorf("Se esperaba un error por longitud de nombre de sala mayor a 16 carácteres")
	}
}

func TestInviteMessage(t *testing.T) {
	roomname := "La Logia"
	usernames := []string{"Dale", "Laura", "Donna"}
	msg, err := InviteMessage(roomname, usernames)
	if err != nil {
		fmt.Println(err)
	}
	if msg.Type != INVITE {
		t.Errorf("Se esperaba el tipo %s, pero se obtuvo %s", INVITE, msg.Type)
	}
	if msg.Roomname != roomname {
		t.Errorf("Se esperaba el nombre de sala %s, pero se obtuvo %s", roomname, msg.Roomname)
	}
	if len(msg.Usernames) != len(usernames) {
		t.Errorf("Se esperaba %d nombres de usuario, pero se obtuvo %d", len(usernames), len(msg.Usernames))
	}
	for i, u := range usernames {
		if msg.Usernames[i] != u {
			t.Errorf("Se esperaba el nombre de usuario %s, pero se obtuvo %s", u, msg.Usernames[i])
		}
	}
}

func TestInviteMessageError(t *testing.T) {
	roomname := "La corte de los búhos"
	usernames := []string{"Dale"}
	_, err := InviteMessage(roomname, usernames)
	if err == nil {
		t.Errorf("Se esperaba un error por longitud de nombre de sala mayor a 16 carácteres")
	}
}

func TestJoinRoomMessage(t *testing.T) {
	roomname := "La Logia"
	msg, err := JoinRoomMessage(roomname)
	if err != nil {
		fmt.Println(err)
	}
	if msg.Type != JOIN_ROOM {
		t.Errorf("Se esperaba el tipo %s, pero se obtuvo %s", JOIN_ROOM, msg.Type)
	}
	if msg.Roomname != roomname {
		t.Errorf("Se esperaba el nombre de sala %s, pero se obtuvo %s", roomname, msg.Roomname)
	}
}

func TestJoinRoomMessageError(t *testing.T) {
	roomname := "La corte de los búhos"
	_, err := JoinRoomMessage(roomname)
	if err == nil {
		t.Errorf("Se esperaba un error por longitud de nombre de sala mayor a 16 carácteres")
	}
}

func TestRoomUsersMessage(t *testing.T) {
	roomname := "La Logia"
	msg, err := RoomUsersMessage(roomname)
	if err != nil {
		fmt.Println(err)
	}
	if msg.Type != ROOM_USERS {
		t.Errorf("Se esperaba el tipo %s, pero se obtuvo %s", ROOM_USERS, msg.Type)
	}
	if msg.Roomname != roomname {
		t.Errorf("Se esperaba el nombre de sala %s, pero se obtuvo %s", roomname, msg.Roomname)
	}
}

func TestRoomUsersMessageError(t *testing.T) {
	roomname := "La corte de los búhos"
	_, err := RoomUsersMessage(roomname)
	if err == nil {
		t.Errorf("Se esperaba un error por longitud de nombre de sala mayor a 16 carácteres")
	}
}

func TestRoomTextMessage(t *testing.T) {
	roomname := "La Logia"
	text := "Está pasando de nuevo"
	msg, err := RoomTextMessage(roomname, text)
	if err != nil {
		fmt.Println(err)
	}
	if msg.Type != ROOM_TEXT {
		t.Errorf("Se esperaba el tipo %s, pero se obtuvo %s", ROOM_TEXT, msg.Type)
	}
	if msg.Roomname != roomname {
		t.Errorf("Se esperaba el nombre de sala %s, pero se obtuvo %s", roomname, msg.Roomname)
	}
	if msg.Text != text {
		t.Errorf("Se esperaba el texto %s, pero se obtuvo %s", text, msg.Text)
	}
}

func TestRoomTextMessageError(t *testing.T) {
	roomname := "La corte de los búhos"
	text := "Está pasando de nuevo"
	_, err := RoomTextMessage(roomname, text)
	if err == nil {
		t.Errorf("Se esperaba un error por longitud de nombre de sala mayor a 16 carácteres")
	}
}

func TestLeaveRoomMessage(t *testing.T) {
	roomname := "La Logia"
	msg, err := LeaveRoomMessage(roomname)
	if err != nil {
		fmt.Println(err)
	}
	if msg.Type != LEAVE_ROOM {
		t.Errorf("Se esperaba el tipo %s, pero se obtuvo %s", LEAVE_ROOM, msg.Type)
	}
	if msg.Roomname != roomname {
		t.Errorf("Se esperaba el nombre de sala %s, pero se obtuvo %s", roomname, msg.Roomname)
	}
}

func TestLeaveRoomMessageError(t *testing.T) {
	roomname := "La corte de los búhos"
	_, err := LeaveRoomMessage(roomname)
	if err == nil {
		t.Errorf("Se esperaba un error por longitud de nombre de sala mayor a 16 carácteres")
	}
}

func TestDisconnectMessage(t *testing.T) {
	msg, err := DisconnectMessage()
	if err != nil {
		fmt.Println(err)
	}
	if msg.Type != DISCONNECT {
		t.Errorf("Se esperaba el tipo %s, pero se obtuvo %s", DISCONNECT, msg.Type)
	}
}
