package panels

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type CloseMsg struct{}

// Modelo para mostrar notificaciones del chat.
type NotificationModal struct {
	Layout
	prompt       string
	notification string
	active       bool
	fatal        bool
}

func NewNotificationModal(prompt string, notification string, isFatal bool) *NotificationModal {
	return &NotificationModal{
		prompt:       prompt,
		notification: notification,
		active:       true,
		fatal:        isFatal,
	}
}

// Init inicia un tick para cerrar la notificación tras 3 segundos.
func (m *NotificationModal) Init() tea.Cmd {
	return tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
		return CloseMsg{}
	})
}

func (m *NotificationModal) IsCapturingInput() bool {
	return m.active
}

func (m *NotificationModal) Focus() tea.Cmd {
	return nil
}

func (m *NotificationModal) Update(msg tea.Msg) (Modal, tea.Cmd) {
	switch msg := msg.(type) {
	case CloseMsg:
		return m.closeOrExit()

	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "esc", "q":
			return m.closeOrExit()
		}
	}
	return m, nil
}

func (m *NotificationModal) View() string {
	return renderModalCard(m.width, m.height, m.prompt, m.notification)
}

func (m *NotificationModal) closeOrExit() (Modal, tea.Cmd) {
	m.active = false

	if m.fatal {
		return m, tea.Quit
	}

	return m, nil
}
