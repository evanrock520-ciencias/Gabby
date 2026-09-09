package messages

type ListType int

type ListResultMsg interface {
	isListResult()
}

type EnterRoomMsg struct {
	Roomname string
}

func (EnterRoomMsg) isListResult() {}

type EnterDMMsg struct {
	Username string
}

func (EnterDMMsg) isListResult() {}
