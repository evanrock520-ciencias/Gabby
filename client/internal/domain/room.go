package domain

type RoomKind int

const (
	RoomGlobal RoomKind = iota
	RoomChannel
	RoomDM
)

// ChatEntry es cualquier cosa que puede mostrarse en el historial del chat.
type ChatEntry interface {
	chatEntry()
}

// ChatMessage representa un mensaje de texto de un usuario.
type ChatMessage struct {
	Username string
	Message  string
}

func (ChatMessage) chatEntry() {}

// ChatEvent representa un evento del sistema.
type ChatEvent struct {
	Text string
}

func (ChatEvent) chatEntry() {}

// BaseRoom contiene los datos comunes de una sala.
type BaseRoom struct {
	Name    string
	Entries []ChatEntry
}

func (room *BaseRoom) AddEntry(entry ChatEntry) {
	room.Entries = append(room.Entries, entry)
}

// Room representa cualquier tipo de sala.
type Room struct {
	BaseRoom
	Users      []User
	kind       RoomKind
	targetUser string
}

// Kind da el tipo de la sala (DM, Channel, o Global)
func (room *Room) Kind() RoomKind {
	return room.kind
}

// TargetUser da el usuario al que se dirige un DM.
func (room *Room) TargetUser() string {
	return room.targetUser
}

func (room *Room) IsDM() bool {
	return room.kind == RoomDM
}

func (room *Room) CanLeave() bool {
	return room.kind == RoomChannel
}

func (room *Room) Title() string {
	switch room.kind {
	case RoomGlobal:
		return "Global"
	case RoomDM:
		return "@" + room.targetUser
	default:
		return "# " + room.Name
	}
}

func (room *Room) AddUser(user User) {
	room.Users = append(room.Users, user)
}

func CreateRoom(kind RoomKind, name string) *Room {
	switch kind {
	case RoomGlobal:
		return NewGlobalRoom()
	case RoomDM:
		return NewDMRoom(name)
	default:
		return NewChannelRoom(name)
	}
}

func NewGlobalRoom() *Room {
	return &Room{
		BaseRoom: BaseRoom{
			Name:    "Global",
			Entries: []ChatEntry{},
		},
		Users: []User{},
		kind:  RoomGlobal,
	}
}

func NewChannelRoom(name string) *Room {
	return &Room{
		BaseRoom: BaseRoom{
			Name:    name,
			Entries: []ChatEntry{},
		},
		Users: []User{},
		kind:  RoomChannel,
	}
}

func NewDMRoom(username string) *Room {
	return &Room{
		BaseRoom: BaseRoom{
			Name:    "@" + username,
			Entries: []ChatEntry{},
		},
		Users:      []User{},
		kind:       RoomDM,
		targetUser: username,
	}
}
