package messages

import (
	"client/internal/domain"
)

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

type ShowRoomUserlist struct {
	Roomname string
	Users    []domain.User
}

func (msg ShowRoomUserlist) isInternal() {}

type ShowNotification struct {
	Prompt       string
	Notification string
	IsFatal      bool
}

func (msg ShowNotification) isInternal() {}
