package config

import (
	"context"
	"encoding/json"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/aptible/cloud-cli/client"
)

/*
CloudConfig
Core common for the Cloud API
*/
type CloudConfig struct {
	Vconfig *viper.Viper
	Cc      client.CloudClient
	Ctx     context.Context
}

// CobraRunE - alias for Cobra's RunE
type CobraRunE func(cmd *cobra.Command, args []string) error

// FindToken - tries to find an aptible token in various paths
func FindToken(home string, domain string) (string, error) {
	var tokenObj map[string]string
	text, err := os.ReadFile(path.Join(home, ".aptible", "tokens.json"))
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(text, &tokenObj)
	if err != nil {
		panic(err)
	}

	return tokenObj[domain], nil
}

func NewCloudConfig(v *viper.Viper) *CloudConfig {
	host := v.GetString("api-domain")
	token := v.GetString("token")
	debug := v.GetBool("debug")
	cc := client.NewClient(debug, host, token)

	return &CloudConfig{
		Vconfig: v,
		Cc:      cc,
	}
}

func configCreateRun() CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		// TODO
		return nil
	}
}

func NewConfigCmd() *cobra.Command {
	configCmd := &cobra.Command{
		Use:     "common",
		Short:   "The common subcommand has assorted common utils associated with the CLI.",
		Long:    "The common subcommand has assorted common utils associated with the CLI.",
		Aliases: []string{"common", "c"},
	}

	configCreateCmd := &cobra.Command{
		Use:     "create common",
		Short:   "provision a fresh common file.",
		Long:    "provision a fresh common file.",
		Aliases: []string{"c"},
		Args:    cobra.MinimumNArgs(0),
		RunE:    configCreateRun(),
	}

	configCmd.AddCommand(configCreateCmd)

	return configCmd
}
