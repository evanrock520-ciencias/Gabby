package messages

type SendMessageMsg struct {
	Roomname string
	Username string
	Message  string
}

type LeftRoomMsg struct {
	Roomname string
}
