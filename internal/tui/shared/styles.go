package shared

import (
	"charm.land/lipgloss/v2"
)

var (
	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7DC4E4")).
		MarginBottom(1)

	Subtitle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9399B2"))

	Label = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A6DA95"))

	ActiveItem = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A6DA95")).
		Bold(true)

	InactiveItem = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#CAD3F8"))

	Success = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A6DA95")).
		Bold(true)

	Error = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ED8796"))

	Dim = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6E738D"))

	Box = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#494D64")).
		Padding(1, 2)

	Prompt = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#CAD3F8"))
)
