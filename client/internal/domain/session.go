package domain

import "slices"

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

func (s *SessionState) AddRoom(kind RoomKind, name string) *Room {
	if s.Rooms == nil {
		s.Rooms = make(map[string]*Room)
	}
	if room, exists := s.Rooms[name]; exists {
		return room
	}
	newRoom := CreateRoom(kind, name)
	s.Rooms[name] = newRoom
	return newRoom
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

func (s *SessionState) RemoveUser(username string) {
	i := slices.IndexFunc(s.Users, func(u User) bool {
		return u.Username == username
	})
	if i != -1 {
		s.Users = slices.Delete(s.Users, i, i+1)
	}
}

func (s *SessionState) GetOrCreateDM(username string) *Room {
	if s.DMs == nil {
		s.DMs = make(map[string]*Room)
	}
	if dm, exists := s.DMs[username]; exists {
		return dm
	}
	newDMRoom := NewDMRoom(username)
	s.DMs[username] = newDMRoom
	return newDMRoom
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

func (s *SessionState) SetUserStatus(username string, status Status) {
	for i, item := range s.Users {
		if item.Username == username {
			s.Users[i].Status = status
			return
		}
	}
}
