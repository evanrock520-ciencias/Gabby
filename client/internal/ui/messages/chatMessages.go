package messages

type ChatResultMsg interface {
	isChatResult()
}

type SendMessageMsg struct {
	Roomname string
	Username string
	Message  string
}

func (SendMessageMsg) isChatResult() {}

type LeftRoomMsg struct {
	Roomname string
}

func (LeftRoomMsg) isChatResult() {}
