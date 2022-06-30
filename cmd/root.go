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
	apiDomain  string
	orgID      string
	debug      bool
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
	rootCmd.PersistentFlags().StringVar(&authDomain, "auth-domain", "auth.aptible.com", "auth domain")
	rootCmd.PersistentFlags().StringVar(&apiDomain, "api-domain", "cloud-api.aptible.com", "api domain")
	rootCmd.PersistentFlags().StringVar(&orgID, "org", "", "organization id")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug logging")

	viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag("auth-domain", rootCmd.PersistentFlags().Lookup("auth-domain"))
	viper.BindPFlag("api-domain", rootCmd.PersistentFlags().Lookup("api-domain"))
	viper.BindPFlag("org", rootCmd.PersistentFlags().Lookup("org"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	envCmd := NewEnvCmd()
	dsCmd := NewDatastoreCmd()
	orgCmd := NewOrgCmd()

	rootCmd.AddCommand(
		dsCmd,
		envCmd,
		orgCmd,
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
			token = findToken(home, fmt.Sprintf("https://%s", authDomain))
		}
		vconfig.Set("token", token)

		if err := vconfig.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", vconfig.ConfigFileUsed())
		}
	}
}