package styles

import (
	"hash/fnv"

	"client/internal/domain"

	"github.com/charmbracelet/lipgloss"
)

const (
	PrimaryColor   = lipgloss.Color("#4D87EB")
	AccentColor    = lipgloss.Color("#624DEB")
	TextColor      = lipgloss.Color("#DDDDE0")
	MutedColor     = lipgloss.Color("240")
	SecondaryColor = lipgloss.Color("#D6EB4D")
	TerciaryColor  = lipgloss.Color("#EB624D")
	ActiveColor    = lipgloss.Color("#4DEB62")
)

func ColorByStatus(status domain.Status) lipgloss.Color {
	switch status {
	case domain.ACTIVE:
		return ActiveColor
	case domain.BUSY:
		return TerciaryColor
	case domain.AWAY:
		return SecondaryColor
	}
	return PrimaryColor
}

var UserPalette = []lipgloss.Color{
	lipgloss.Color("#4D87EB"),
	lipgloss.Color("#624DEB"),
	lipgloss.Color("#B14DEB"),
	lipgloss.Color("#EBB14D"),
	lipgloss.Color("#87EB4D"),
}

func ColorByUsername(username string) lipgloss.Color {
	hash := fnv.New32a()
	hash.Write([]byte(username))
	return UserPalette[hash.Sum32()%uint32(len(UserPalette))]
}
