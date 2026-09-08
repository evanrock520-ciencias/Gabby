package domain

type ChatMessage struct {
	Username string
	Message  string
}

type Room struct {
	Name     string
	Users    []string
	Messages []ChatMessage
}

func NewRoom(roomname string) Room {
	return Room{
		Name:     roomname,
		Users:    []string{},
		Messages: []ChatMessage{},
	}
}

func (room *Room) AddUser(username string) {
	room.Users = append(room.Users, username)
}

func (room *Room) AddMessage(msg ChatMessage) {
	room.Messages = append(room.Messages, msg)
}
