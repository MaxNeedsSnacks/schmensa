package utils

import (
	gloss "github.com/charmbracelet/lipgloss"
)

func AppStyle() gloss.Style {
	return gloss.NewStyle().Padding(1, 2)
}

func ErrorStyle() gloss.Style {
	return gloss.NewStyle().
		Foreground(gloss.Color("15")).
		Background(gloss.Color("9")).
		Padding(0, 1)
}

func ErrorMessageStyle() gloss.Style {
	return gloss.NewStyle().Foreground(gloss.Color("9"))
}
