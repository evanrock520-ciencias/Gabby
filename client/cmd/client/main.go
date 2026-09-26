package main

import (
	"client/internal/domain"
	"client/internal/network"
	"client/internal/ui"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	f, err := tea.LogToFile("debug.log", "client")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	port := "localhost:9090"
	if len(os.Args) >= 2 {
		port = os.Args[1]
	}

	conn := network.ConnectionManager{}
	conn.Dial(port)
	go conn.Listen()

	session := domain.NewSessionState(domain.User{})

	p := tea.NewProgram(ui.NewModel(session, &conn), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
