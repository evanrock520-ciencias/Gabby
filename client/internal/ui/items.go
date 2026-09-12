package ui

import (
	"client/internal/domain"
	"client/internal/ui/styles"

	"github.com/charmbracelet/lipgloss"
)

type UserItem struct {
	domain.User
}

func NewUserItem(user domain.User) UserItem {
	return UserItem{User: user}
}

func (user UserItem) Value() string {
	return user.Username
}

func (user UserItem) Render(selected bool, width int) string {
	badgeColor := styles.ColorByStatus(user.Status)

	if selected {
		badge := styles.SelectedStyle.Foreground(badgeColor).Render(" ● ")
		nameWidth := max(0, width-lipgloss.Width(badge))
		name := styles.SelectedStyle.Width(nameWidth).Render(user.Username)

		return badge + name
	}

	badge := lipgloss.NewStyle().Foreground(badgeColor).Render(" ● ")
	return badge + styles.TextStyle.Render(user.Username)
}

type RoomItem struct {
	*domain.Room
}

func NewRoomItem(room *domain.Room) RoomItem {
	return RoomItem{Room: room}
}

func (room RoomItem) Value() string {
	return room.Name
}

func (room RoomItem) Render(selected bool, width int) string {
	if selected {
		return styles.SelectedStyle.Width(width).Render(" # " + room.Value())
	}
	return styles.TextStyle.Render(" # " + room.Value())
}
