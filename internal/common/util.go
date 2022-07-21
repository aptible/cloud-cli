package common

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func TeaDebug(config *CloudConfig) {
	if config.Vconfig.GetBool("debug") {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}
}
