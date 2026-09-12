package domain

type ChatMessage struct {
	Username string
	Message  string
}

type Room struct {
	Name     string
	Users    []User
	Messages []ChatMessage
}

func NewRoom(roomname string) Room {
	return Room{
		Name:     roomname,
		Users:    []User{},
		Messages: []ChatMessage{},
	}
}

func (room *Room) AddUser(user User) {
	room.Users = append(room.Users, user)
}

func (room *Room) AddMessage(msg ChatMessage) {
	room.Messages = append(room.Messages, msg)
}
