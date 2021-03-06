package common

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func TeaDebug(config *CloudConfig) func() {
	if !config.Vconfig.GetBool("debug") {
		return func() {}
	}
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	return func() {
		defer f.Close()
	}
}
