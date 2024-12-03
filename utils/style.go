package utils

import "github.com/charmbracelet/lipgloss"

var (
	SuccessMessage = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#22c55e"))
	WarningMessage = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#f97316"))
	ErrorMessage = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#ff5733"))
)
