package messages

type CreateRoomMsg struct {
	Roomname string
}

type InvitateMsg struct {
	Roomname string
	Users    []string
}

type SetUserMsg struct {
	Username string
}

type LeftChat struct {
}
