package ui

import (
	"client/internal/domain"
	"client/internal/ui/panels"
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

func (user UserItem) Render(state panels.ItemState, width int) string {
	badgeColor := styles.ColorByStatus(user.Status)

	switch state {
	case panels.ItemSelected:
		badge := styles.SelectedStyle.Foreground(badgeColor).Render(" ● ")
		nameWidth := max(0, width-lipgloss.Width(badge))
		name := styles.SelectedStyle.Width(nameWidth).Render(user.Username)
		return badge + name
	case panels.ItemChosen:
		badge := styles.ChosenStyle.Foreground(badgeColor).Render(" ● ")
		nameWidth := max(0, width-lipgloss.Width(badge))
		name := styles.ChosenStyle.Width(nameWidth).Render(user.Username)
		return badge + name
	case panels.ItemBoth:
		badge := styles.ChosenStyle.Foreground(badgeColor).Bold(true).Render(" ● ")
		nameWidth := max(0, width-lipgloss.Width(badge))
		name := styles.ChosenStyle.Bold(true).Width(nameWidth).Render(user.Username)
		return badge + name
	default:
		badge := lipgloss.NewStyle().Foreground(badgeColor).Render(" ● ")
		return badge + styles.TextStyle.Width(width-lipgloss.Width(badge)).Render(user.Username)
	}
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

func (room RoomItem) Render(state panels.ItemState, width int) string {
	switch state {
	case panels.ItemSelected:
		return styles.SelectedStyle.Width(width).Render(" # " + room.Value())
	case panels.ItemChosen:
		return styles.ChosenStyle.Width(width).Render(" # " + room.Value())
	case panels.ItemBoth:
		return styles.ChosenStyle.Bold(true).Width(width).Render(" # " + room.Value())
	default:
		return styles.TextStyle.Width(width).Render(" # " + room.Value())
	}
}
