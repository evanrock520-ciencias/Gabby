package messages

import "client/internal/domain"

type ChangeStatusMsg struct {
	Status domain.Status
}
