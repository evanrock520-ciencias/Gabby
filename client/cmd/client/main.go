package main

import (
	"client/internal/domain"
	"client/internal/network"
	"client/internal/ui"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	f, err := tea.LogToFile("debug.log", "client")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	conn := network.ConnectionManager{}
	conn.Dial("localhost:9090")
	go conn.Listen()

	session := domain.SessionState{
		Users: []domain.User{},
		Rooms: make(map[string]*domain.Room),
		DMs:   make(map[string]*domain.Room),
	}

	p := tea.NewProgram(ui.NewModel(session, &conn), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
