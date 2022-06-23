package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var datastoreCmd = &cobra.Command{
	Use:     "datastore",
	Short:   "The datastore subcommand helps manage your Aptible datastores.",
	Long:    `The datastore subcommand helps manage your Aptible datastores.`,
	Aliases: []string{"database", "ds", "db"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Echo: " + strings.Join(args, " "))
	},
}

var dsCreateCmd = &cobra.Command{
	Use:     "create",
	Short:   "provision a new datastore.",
	Long:    `The datastore create command will provision a new datastore.`,
	Aliases: []string{"c", "deploy"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Echo: " + strings.Join(args, " "))
	},
}

var dsDestroyCmd = &cobra.Command{
	Use:     "destroy",
	Short:   "permentantly remove the datastore.",
	Long:    `The datastore destroy command will permentantly remove the datastore.`,
	Aliases: []string{"d", "delete", "rm", "remove"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Echo: " + strings.Join(args, " "))
	},
}

var dsListCmd = &cobra.Command{
	Use:     "list",
	Short:   "list all datastores within an organization.",
	Long:    `The datastore list command will list all datastores within an organization.`,
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Echo: " + strings.Join(args, " "))
	},
}

var engine string
var engineVersion string
var env string

func init() {
	dsCreateCmd.Flags().StringVarP(&engine, "engine", "e", "", "the datastore engine, e.g. rds/postgres, rds/mysql, etc.")
	dsCreateCmd.Flags().StringVarP(&engineVersion, "version", "v", "", "the engine version, e.g. 14.2")
	dsListCmd.Flags().StringVar(&env, "env", "", "list datastores within an environment")

	datastoreCmd.AddCommand(dsCreateCmd)
	datastoreCmd.AddCommand(dsDestroyCmd)
	datastoreCmd.AddCommand(dsListCmd)
}
