package domain

type SessionState struct {
	Users       []User
	Rooms       map[string]*Room
	DMs         map[string]*Room
	CurrentUser User
}
