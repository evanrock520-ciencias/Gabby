package domain

type Status string

const (
	ACTIVE Status = "ACTIVE"
	AWAY   Status = "AWAY"
	BUSY   Status = "BUSY"
)

func (s Status) Next() Status {
	switch s {
	case ACTIVE:
		return AWAY
	case AWAY:
		return BUSY
	case BUSY:
		return ACTIVE
	default:
		return ACTIVE
	}
}

type User struct {
	Username string
	Status   Status
}
