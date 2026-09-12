package messages

type ModalResultMsg interface {
	isModalResult()
}

type CreateRoomMsg struct {
	Roomname string
}

func (CreateRoomMsg) isModalResult() {}

type InvitateMsg struct {
	Roomname string
}

func (InvitateMsg) isModalResult() {}

type SetUserMsg struct {
	Username string
}

func (SetUserMsg) isModalResult() {}

type LeftChat struct {
}

func (LeftChat) isModalResult() {}
