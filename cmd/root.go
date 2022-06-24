package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile    string
	token      string
	authDomain string
	orgID      string
)

var desc = `aptible is a command line interface to the Aptible.com platform.

It allows users to manage authentication, application launch,
deployment, logging, and more with just the one command.

* Provision an app with the app create command
* Provision a datastore with the datastore create command
* View a deployed web application with the open command
* View detailed information about an app or datastore with the info command

To read more, use the docs command to view Aptible's help on the web.`

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "aptible",
		Short: "aptible is a command line interface to the Aptible.com platform.",
		Long:  desc,
	}

	cobra.OnInitialize(initConfig())

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.aptible.yaml)")
	rootCmd.PersistentFlags().StringVar(&token, "token", "", "jwt token")
	rootCmd.PersistentFlags().StringVar(&authDomain, "auth-domain", "https://auth.aptible.com", "auth domain")
	rootCmd.PersistentFlags().StringVar(&orgID, "org", "", "organization id")

	envCmd := NewEnvCmd()
	dsCmd := NewDatastoreCmd()

	rootCmd.AddCommand(
		dsCmd,
		envCmd,
	)

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(root *cobra.Command) {
	err := root.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func initConfig() func() {
	return func() {
		vconfig := viper.GetViper()
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		if cfgFile != "" {
			// Use config file from the flag.
			vconfig.SetConfigFile(cfgFile)
		} else {
			vconfig.AddConfigPath(home)
			vconfig.SetConfigName(".aptible")
			vconfig.SetConfigType("yaml")
		}

		vconfig.AutomaticEnv()

		if token == "" {
			token = findToken(home, authDomain)
		}
		vconfig.Set("token", token)

		if err := vconfig.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", vconfig.ConfigFileUsed())
		}
	}
}
