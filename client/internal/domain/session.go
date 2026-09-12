package domain

type SessionState struct {
	Users       []User
	Rooms       map[string]*Room
	DMs         map[string]*Room
	CurrentUser User
}

func NewSessionState(currentUser User) SessionState {
	return SessionState{
		Users:       []User{},
		Rooms:       make(map[string]*Room),
		DMs:         make(map[string]*Room),
		CurrentUser: currentUser,
	}
}

func (s *SessionState) AddRoom(name string) *Room {
	if s.Rooms == nil {
		s.Rooms = make(map[string]*Room)
	}
	if room, exists := s.Rooms[name]; exists {
		return room
	}
	newRoom := NewRoom(name)
	s.Rooms[name] = &newRoom
	return &newRoom
}

func (s *SessionState) GetRoom(name string) (*Room, bool) {
	if s.Rooms == nil {
		return nil, false
	}
	room, ok := s.Rooms[name]
	return room, ok
}

func (s *SessionState) RemoveRoom(name string) {
	if s.Rooms != nil {
		delete(s.Rooms, name)
	}
}

func (s *SessionState) GetOrCreateDM(username string) *Room {
	if s.DMs == nil {
		s.DMs = make(map[string]*Room)
	}
	if dm, exists := s.DMs[username]; exists {
		return dm
	}
	newDMRoom := NewRoom("@" + username)
	s.DMs[username] = &newDMRoom
	return &newDMRoom
}

func (s *SessionState) Usernames() []string {
	names := make([]string, 0, len(s.Users))
	for _, user := range s.Users {
		names = append(names, user.Username)
	}
	return names
}

func (s *SessionState) RoomNames() []string {
	names := make([]string, 0, len(s.Rooms))
	for _, room := range s.Rooms {
		names = append(names, room.Name)
	}
	return names
}

func (s *SessionState) SetCurrentUser(username string, status Status) {
	s.CurrentUser.Username = username
	s.CurrentUser.Status = status
}

func (s *SessionState) SetStatus(status Status) {
	s.CurrentUser.Status = status
}
