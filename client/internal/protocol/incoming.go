package protocol

type TypeS2C string

const (
	RESPONSE         TypeS2C = "RESPONSE"
	NEW_USER         TypeS2C = "NEW_USER"
	NEW_STATUS       TypeS2C = "NEW_STATUS"
	USER_LIST        TypeS2C = "USER_LIST"
	TEXT_FROM        TypeS2C = "TEXT_FROM"
	PUBLIC_TEXT_FROM TypeS2C = "PUBLIC_TEXT_FROM"
	JOINED_ROOM      TypeS2C = "JOINED_ROOM"
	ROOM_USER_LIST   TypeS2C = "ROOM_USER_LIST"
	ROOM_TEXT_FROM   TypeS2C = "ROOM_TEXT_FROM"
	LEFT_ROOM        TypeS2C = "LEFT_ROOM"
	DISCONNECTED     TypeS2C = "DISCONNECTED"
	INVITATION       TypeS2C = "INVITATION"
)

type ServerMessage struct {
	Type      TypeS2C  `json:"type,omitempty"`
	Username  string   `json:"username,omitempty"`
	Room      string   `json:"room,omitempty"`
	Status    Status   `json:"status,omitempty"`
	Text      string   `json:"text,omitempty"`
	Usernames []string `json:"usernames,omitempty"`
}
