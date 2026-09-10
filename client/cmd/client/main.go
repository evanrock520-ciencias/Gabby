package main

import (
	"client/internal/domain"
	"client/internal/ui"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	//TODO: Obtener los datos del servidor
	users := []domain.User{{Username: "Evan", Status: domain.ACTIVE}, {Username: "Derek", Status: domain.AWAY}}
	session := domain.SessionState{
		Users: users,
		Rooms: map[string]*domain.Room{
			"Global": {Name: "Global", Users: users, Messages: []domain.ChatMessage{}},
		},
		DMs: make(map[string]*domain.Room),
	}
	p := tea.NewProgram(ui.NewModel(session), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
	}
}
