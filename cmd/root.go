package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile    string
	token      string
	authDomain string
)

var rootCmd = &cobra.Command{
	Use:   "aptible",
	Short: "aptible is a command line interface to the Aptible.com platform.",
	Long: `aptible is a command line interface to the Aptible.com platform.

It allows users to manage authentication, application launch,
deployment, logging, and more with just the one command.

* Provision an app with the app create command
* Provision a datastore with the datastore create command
* View a deployed web application with the open command
* View detailed information about an app or datastore with the info command

To read more, use the docs command to view Aptible's help on the web.`,
}

var tokenObj map[string]string

func findToken(home string, domain string) string {
	text, err := ioutil.ReadFile(path.Join(home, ".aptible", "tokens.json"))
	if err != nil {
		return ""
	}
	json.Unmarshal(text, &tokenObj)
	return string(tokenObj[domain])
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.aptible.yaml)")
	rootCmd.PersistentFlags().StringVar(&token, "token", "", "jwt token")
	rootCmd.PersistentFlags().StringVar(&authDomain, "auth-domain", "https://auth.aptible.com", "auth domain")

	rootCmd.AddCommand(datastoreCmd)
	rootCmd.AddCommand(envCmd)
}

func initConfig() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	if token == "" {
		token = findToken(home, authDomain)
	}
	fmt.Println(token)

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(home)
		viper.SetConfigName(".aptible")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

}
