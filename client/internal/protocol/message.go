package protocol

import "client/internal/domain"

type Message struct {
	Username  string        `json:"username,omitempty"`
	Roomname  string        `json:"roomname,omitempty"`
	Status    domain.Status `json:"status,omitempty"`
	Text      string        `json:"text,omitempty"`
	Usernames []string      `json:"usernames,omitempty"`
}
