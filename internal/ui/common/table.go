package common

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

type Model struct {
	simpleTable table.Model
}

func DefaultRowStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(white).
		Align(lipgloss.Center)
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
