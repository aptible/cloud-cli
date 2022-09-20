package common

import (
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
}

func DefaultRowStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(white).
		Align(lipgloss.Center)
}

func LeftRowStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(white).
		Align(lipgloss.Left)
}

func DisabledRowStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(red).
		Faint(true).
		Align(lipgloss.Center)
}

func ActiveRowStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(green).
		Align(lipgloss.Center)
}

func PendingRowStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(yellowGreen).
		Align(lipgloss.Center)
}
