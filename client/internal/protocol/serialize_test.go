package protocol

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func jsonValueKey(value string, key string) string {
	return fmt.Sprintf("\"%s\":\"%s\"", value, key)
}

func jsonArrayKey(value string, array []string) string {
	arrayJson, _ := json.Marshal(array)
	return fmt.Sprintf("\"%s\":%s", value, string(arrayJson))
}

func typeError(msgType string, data string) string {
	return fmt.Sprintf("No se encontró el campo `%s` o contiene una llave diferente: %s", msgType, data)
}

func contains(data string, value string, key string) bool {
	return strings.Contains(data, jsonValueKey(value, key))
}

func containsArray(data string, value string, array []string) bool {
	return strings.Contains(data, jsonArrayKey(value, array))
}

func TestSerializeIdentifyMessage(t *testing.T) {
	msg, _ := IdentifyMessage("Evan")
	data, _ := serialize(msg)
	if !contains(data, "type", string(IDENTIFY)) {
		t.Fatal(typeError("type", data))
	}
	if !contains(data, "username", "Evan") {
		t.Fatal(typeError("username", data))
	}
}

func TestSerializeStatusMessage(t *testing.T) {
	msg, _ := StatusMessage(AWAY)
	data, _ := serialize(msg)
	if !contains(data, "type", string(STATUS)) {
		t.Fatal(typeError("type", data))
	}
	if !contains(data, "status", string(AWAY)) {
		t.Fatal(typeError("status", data))
	}
}

func TestSerializeUsersMessage(t *testing.T) {
	msg, _ := UsersMessage()
	data, _ := serialize(msg)
	if !contains(data, "type", string(USERS)) {
		t.Fatal(typeError("type", data))
	}
}

func TestSerializeTextMessage(t *testing.T) {
	msg, _ := TextMessage("Evan", "Hola")
	data, _ := serialize(msg)
	if !contains(data, "type", string(TEXT)) {
		t.Fatal(typeError("type", data))
	}
	if !contains(data, "username", "Evan") {
		t.Fatal(typeError("username", data))
	}
	if !contains(data, "text", "Hola") {
		t.Fatal(typeError("text", data))
	}
}

func TestSerializePublicTextMessage(t *testing.T) {
	text := "Hola camaradas."
	msg, _ := PublicTextMessage(text)
	data, _ := serialize(msg)
	if !contains(data, "type", string(PUBLIC_TEXT)) {
		t.Fatal(typeError("type", data))
	}
	if !contains(data, "text", text) {
		t.Fatal(typeError("text", data))
	}
}

func TestSerializeNewRoomMessage(t *testing.T) {
	room := "Sala 1"
	msg, _ := NewRoomMessage(room)
	data, _ := serialize(msg)
	if !contains(data, "type", string(NEW_ROOM)) {
		t.Fatal(typeError("type", data))
	}
	if !contains(data, "roomname", room) {
		t.Fatal(typeError("roomname", data))
	}
}

func TestSerializeInviteMessage(t *testing.T) {
	room := "Sala 1"
	usernames := []string{"Luis", "Antonio", "Fernando"}
	msg, _ := InviteMessage(room, usernames)
	data, _ := serialize(msg)
	if !contains(data, "type", string(INVITE)) {
		t.Fatal(typeError("type", data))
	}
	if !contains(data, "roomname", room) {
		t.Fatal(typeError("roomname", data))
	}
	if !containsArray(data, "usernames", usernames) {
		t.Fatal(typeError("usernames", data))
	}
}

func TestSerializeJoinRoomMessage(t *testing.T) {
	room := "Sala 1"
	msg, _ := JoinRoomMessage(room)
	data, _ := serialize(msg)
	if !contains(data, "type", string(JOIN_ROOM)) {
		t.Fatal(typeError("type", data))
	}
	if !contains(data, "roomname", room) {
		t.Fatal(typeError("roomname", data))
	}
}

func TestSerializeRoomUsersMessage(t *testing.T) {
	room := "Sala 1"
	msg, _ := RoomUsersMessage(room)
	data, _ := serialize(msg)
	if !contains(data, "type", string(ROOM_USERS)) {
		t.Fatal(typeError("type", data))
	}
	if !contains(data, "roomname", room) {
		t.Fatal(typeError("roomname", data))
	}
}

func TestSerializeRoomTextMessage(t *testing.T) {
	room := "Sala 1"
	text := "¡Hola sala 1!"
	msg, _ := RoomTextMessage(room, text)
	data, _ := serialize(msg)
	if !contains(data, "type", string(ROOM_TEXT)) {
		t.Fatal(typeError("type", data))
	}
	if !contains(data, "roomname", room) {
		t.Fatal(typeError("roomname", data))
	}
	if !contains(data, "text", text) {
		t.Fatal(typeError("text", data))
	}
}

func TestSerializeLeaveRoomMessage(t *testing.T) {
	room := "Sala 1"
	msg, _ := LeaveRoomMessage(room)
	data, _ := serialize(msg)
	if !contains(data, "type", string(LEAVE_ROOM)) {
		t.Fatal(typeError("type", data))
	}
	if !contains(data, "roomname", room) {
		t.Fatal(typeError("roomname", data))
	}
}

func TestSerializeDisconnectMessage(t *testing.T) {
	msg, _ := DisconnectMessage()
	data, _ := serialize(msg)
	if !contains(data, "type", string(DISCONNECT)) {
		t.Fatal(typeError("type", data))
	}
}
