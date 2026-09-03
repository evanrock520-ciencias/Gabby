package protocol

import (
	"encoding/json"
	"fmt"
	"reflect"
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
	data, _ := Serialize(msg)
	if !contains(data, "type", string(IDENTIFY)) {
		t.Fatal(typeError("type", data))
	}
	if !contains(data, "username", "Evan") {
		t.Fatal(typeError("username", data))
	}
}

func TestSerializeStatusMessage(t *testing.T) {
	msg, _ := StatusMessage(AWAY)
	data, _ := Serialize(msg)
	if !contains(data, "type", string(STATUS)) {
		t.Fatal(typeError("type", data))
	}
	if !contains(data, "status", string(AWAY)) {
		t.Fatal(typeError("status", data))
	}
}

func TestSerializeUsersMessage(t *testing.T) {
	msg, _ := UsersMessage()
	data, _ := Serialize(msg)
	if !contains(data, "type", string(USERS)) {
		t.Fatal(typeError("type", data))
	}
}

func TestSerializeTextMessage(t *testing.T) {
	msg, _ := TextMessage("Evan", "Hola")
	data, _ := Serialize(msg)
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
	data, _ := Serialize(msg)
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
	data, _ := Serialize(msg)
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
	data, _ := Serialize(msg)
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
	data, _ := Serialize(msg)
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
	data, _ := Serialize(msg)
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
	data, _ := Serialize(msg)
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
	data, _ := Serialize(msg)
	if !contains(data, "type", string(LEAVE_ROOM)) {
		t.Fatal(typeError("type", data))
	}
	if !contains(data, "roomname", room) {
		t.Fatal(typeError("roomname", data))
	}
}

func TestSerializeDisconnectMessage(t *testing.T) {
	msg, _ := DisconnectMessage()
	data, _ := Serialize(msg)
	if !contains(data, "type", string(DISCONNECT)) {
		t.Fatal(typeError("type", data))
	}
}

// Test de Deserialización. Para esto usaré los mensajes de servidor cliente.

func mismatchField(field string, value string, expected string) string {
	return fmt.Sprintf("Field value mismatch %s, got %s, but expected %s", field, value, expected)
}

func TestDeserializeNewUserMessage(t *testing.T) {
	data := `{"type":"NEW_USER","username":"Evan"}`
	msg, _ := Deserialize(data)
	if msg.Type != NEW_USER {
		t.Fatal(mismatchField("type", string(msg.Type), "NEW_USER"))
	}
	if msg.Username != "Evan" {
		t.Fatal(mismatchField("username", msg.Username, "Evan"))
	}
}

func TestDeserializeNewStatusMessage(t *testing.T) {
	data := `{"type":"NEW_STATUS","username":"Evan","status":"BUSY"}`
	msg, _ := Deserialize(data)
	if msg.Type != NEW_STATUS {
		t.Fatal(mismatchField("type", string(msg.Type), "NEW_STATUS"))
	}
	if msg.Username != "Evan" {
		t.Fatal(mismatchField("username", msg.Username, "Evan"))
	}
	if msg.Status != "BUSY" {
		t.Fatal(mismatchField("status", string(msg.Status), "BUSY"))
	}
}

func TestDeserializeUserListsMessage(t *testing.T) {
	data := `{"type":"USER_LIST","usernames":["Derek","Yahel","Luis"]}`
	msg, _ := Deserialize(data)
	if msg.Type != USER_LIST {
		t.Fatal(mismatchField("type", string(msg.Type), "USER_LIST"))
	}
	expectedUsernames := []string{"Derek", "Yahel", "Luis"}
	if !reflect.DeepEqual(msg.Usernames, expectedUsernames) {
		t.Fatal(mismatchField("usernames", fmt.Sprintf("%v", msg.Usernames), fmt.Sprintf("%v", expectedUsernames)))
	}
}

func TestDeserializeTextFromMessage(t *testing.T) {
	data := `{"type":"TEXT_FROM","username":"Evan","text":"Hola amigos"}`
	msg, _ := Deserialize(data)
	if msg.Type != TEXT_FROM {
		t.Fatal(mismatchField("type", string(msg.Type), "TEXT_FROM"))
	}
	if msg.Username != "Evan" {
		t.Fatal(mismatchField("username", msg.Username, "Evan"))
	}
	if msg.Text != "Hola amigos" {
		t.Fatal(mismatchField("text", msg.Text, "Hola amigos"))
	}
}

func TestDeserializePublicTextFromMessage(t *testing.T) {
	data := `{"type":"PUBLIC_TEXT_FROM","username":"Evan","text":"Hola amigos"}`
	msg, _ := Deserialize(data)
	if msg.Type != PUBLIC_TEXT_FROM {
		t.Fatal(mismatchField("type", string(msg.Type), "PUBLIC_TEXT_FROM"))
	}
	if msg.Username != "Evan" {
		t.Fatal(mismatchField("username", msg.Username, "Evan"))
	}
	if msg.Text != "Hola amigos" {
		t.Fatal(mismatchField("text", msg.Text, "Hola amigos"))
	}
}

func TestDeserializeJoinedRoomMessage(t *testing.T) {
	data := `{"type":"JOINED_ROOM","username":"Evan","roomname":"Sala 1"}`
	msg, _ := Deserialize(data)
	if msg.Type != JOINED_ROOM {
		t.Fatal(mismatchField("type", string(msg.Type), "JOINED_ROOM"))
	}
	if msg.Username != "Evan" {
		t.Fatal(mismatchField("username", msg.Username, "Evan"))
	}
	if msg.Room != "Sala 1" {
		t.Fatal(mismatchField("room", msg.Room, "Sala 1"))
	}
}

func TestDeserializeRoomTextFromMessage(t *testing.T) {
	data := `{"type":"ROOM_TEXT_FROM","username":"Evan","roomname":"Sala 1","text":"Hola amigos"}`
	msg, _ := Deserialize(data)
	if msg.Type != ROOM_TEXT_FROM {
		t.Fatal(mismatchField("type", string(msg.Type), "ROOM_TEXT_FROM"))
	}
	if msg.Username != "Evan" {
		t.Fatal(mismatchField("username", msg.Username, "Evan"))
	}
	if msg.Room != "Sala 1" {
		t.Fatal(mismatchField("room", msg.Room, "Sala 1"))
	}
	if msg.Text != "Hola amigos" {
		t.Fatal(mismatchField("text", msg.Text, "Hola amigos"))
	}
}

func TestDeserializeRoomUserListsMessage(t *testing.T) {
	data := `{"type":"ROOM_USER_LIST","roomname":"Sala 1","usernames":["Derek","Yahel","Luis"]}`
	msg, _ := Deserialize(data)
	if msg.Type != ROOM_USER_LIST {
		t.Fatal(mismatchField("type", string(msg.Type), "ROOM_USER_LIST"))
	}
	expectedUsernames := []string{"Derek", "Yahel", "Luis"}
	if !reflect.DeepEqual(msg.Usernames, expectedUsernames) {
		t.Fatal(mismatchField("usernames", fmt.Sprintf("%v", msg.Usernames), fmt.Sprintf("%v", expectedUsernames)))
	}
	if msg.Room != "Sala 1" {
		t.Fatal(mismatchField("room", msg.Room, "Sala 1"))
	}
}

func TestDeserializeLeftRoomMessage(t *testing.T) {
	data := `{"type":"LEFT_ROOM","username":"Evan","roomname":"Sala 1"}`
	msg, _ := Deserialize(data)
	if msg.Type != LEFT_ROOM {
		t.Fatal(mismatchField("type", string(msg.Type), "LEFT_ROOM"))
	}
	if msg.Username != "Evan" {
		t.Fatal(mismatchField("username", msg.Username, "Evan"))
	}
	if msg.Room != "Sala 1" {
		t.Fatal(mismatchField("room", msg.Room, "Sala 1"))
	}
}

func TestDeserializeDisconnectedMessage(t *testing.T) {
	data := `{"type":"DISCONNECTED","username":"Evan"}`
	msg, _ := Deserialize(data)
	if msg.Type != DISCONNECTED {
		t.Fatal(mismatchField("type", string(msg.Type), "DISCONNECTED"))
	}
	if msg.Username != "Evan" {
		t.Fatal(mismatchField("username", msg.Username, "Evan"))
	}
}
