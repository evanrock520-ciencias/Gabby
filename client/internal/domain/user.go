package domain

import "client/internal/protocol"

type User struct {
	Username string
	Status   protocol.Status
}
