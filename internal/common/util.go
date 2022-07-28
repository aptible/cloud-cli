package common

import (
	"fmt"
	"os"
	"strings"

	cloudapiclient "github.com/aptible/cloud-api-clients/clients/go"
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

func FilterAssetsByType(assets []cloudapiclient.AssetOutput, types []string) []cloudapiclient.AssetOutput {
	filteredResults := make([]cloudapiclient.AssetOutput, 0)
	for _, result := range assets {
		for _, _type := range types {
			if strings.Contains(result.Asset, _type) {
				filteredResults = append(filteredResults, result)
			}
		}
	}
	return filteredResults
}
