package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var envCmd = &cobra.Command{
	Use:     "environment",
	Short:   "The env subcommand helps manage your Aptible environments.",
	Long:    `The env subcommand helps manage your Aptible environments.`,
	Aliases: []string{"env", "e"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Echo: " + strings.Join(args, " "))
	},
}

var envCreateCmd = &cobra.Command{
	Use:     "create",
	Short:   "provision a new datastore.",
	Long:    `The environment create command will provision a new environment.`,
	Aliases: []string{"c"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Echo: " + strings.Join(args, " "))
	},
}

var envDestroyCmd = &cobra.Command{
	Use:     "destroy",
	Short:   "permentantly remove the environment.",
	Long:    `The datastore destroy command will permentantly remove the environment.`,
	Aliases: []string{"d", "delete", "rm", "remove"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Echo: " + strings.Join(args, " "))
	},
}

var envListCmd = &cobra.Command{
	Use:     "list",
	Short:   "list all environment within an organization.",
	Long:    `The environment list command will list all environment within an organization.`,
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Echo: " + strings.Join(args, " "))
	},
}

func init() {
	envCmd.AddCommand(envCreateCmd)
	envCmd.AddCommand(envDestroyCmd)
	envCmd.AddCommand(envListCmd)
}
