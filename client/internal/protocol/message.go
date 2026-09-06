package protocol

type Message struct {
	Username  string   `json:"username,omitempty"`
	Roomname  string   `json:"roomname,omitempty"`
	Status    Status   `json:"status,omitempty"`
	Text      string   `json:"text,omitempty"`
	Usernames []string `json:"usernames,omitempty"`
}
