package teadebug

import (
	"fmt"
	"os"

	"github.com/aptible/cloud-cli/config"
	tea "github.com/charmbracelet/bubbletea"
)

func TeaDebug(config *config.CloudConfig) func() {
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
