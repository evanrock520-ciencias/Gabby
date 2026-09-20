package messages

type InternalMsg interface {
	isInternal()
}

type EnterRoom struct {
	Roomname string
}

func (msg EnterRoom) isInternal() {}

type EnterDM struct {
	Username string
}

func (msg EnterDM) isInternal() {}

type LeaveRoomPetition struct {
	Roomname string
}

func (msg LeaveRoomPetition) isInternal() {}
