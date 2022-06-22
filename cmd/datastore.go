package cmd

import (
	"github.com/spf13/cobra"
)

var datastoreCmd = &cobra.Command{
	Use:   "datastore",
	Short: "Manage your Aptible datastores",
	Long:  ``,
}

func init() {
	rootCmd.AddCommand(datastoreCmd)
}
