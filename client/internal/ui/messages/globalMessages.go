package messages

import "client/internal/domain"

type GlobalResultMsg interface {
	isGlobalResult()
}

type ChangeStatusMsg struct {
	Status domain.Status
}

func (ChangeStatusMsg) isGlobalResult() {}
