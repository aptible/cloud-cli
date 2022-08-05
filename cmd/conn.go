package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	cac "github.com/aptible/cloud-api-clients/clients/go"
	"github.com/aptible/cloud-cli/config"
	"github.com/aptible/cloud-cli/lib/conn"
	"github.com/aptible/cloud-cli/ui/fetch"
	"github.com/aptible/cloud-cli/ui/form"
)

type ConnOptions struct {
	OutAsset    string
	InAsset     string
	Description string
}

var connOptions = ConnOptions{}

func connCreateRun() config.CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		config := config.NewCloudConfig(viper.GetViper())
		formResult := form.FormResult{
			Org:         config.Vconfig.GetString("org"),
			Env:         config.Vconfig.GetString("env"),
			InAsset:     connOptions.InAsset,
			OutAsset:    connOptions.OutAsset,
			Description: connOptions.Description,
		}

		err := libconn.ConnCreateForm(config, &formResult)
		if err != nil {
			return nil
		}

		params := cac.ConnectionInput{
			Description:     &formResult.Description,
			OutgoingAssetId: formResult.OutAsset,
		}

		msg := fmt.Sprintf(
			"creating asset connection from (out %s) to (in %s)",
			formResult.OutAsset,
			formResult.InAsset,
		)
		model := fetch.NewModel(msg, func() (interface{}, error) {
			return config.Cc.CreateConnection(
				formResult.Org,
				formResult.Env,
				formResult.InAsset,
				params,
			)
		})
		data, err := fetch.WithOutput(model)
		if err != nil {
			return err
		}
		conn := data.Result.(*cac.ConnectionOutput)
		fmt.Println(conn)

		return nil
	}
}

// NewConnectionCmd - create a connection between two assets.
func NewConnectionCmd() *cobra.Command {
	connCmd := &cobra.Command{
		Use:     "connection",
		Short:   "The connection subcommand helps manage your connections between Aptible resources.",
		Long:    `The connection subcommand helps manage your connections between Aptible resources.`,
		Aliases: []string{"c", "conn"},
	}

	connCreateCmd := &cobra.Command{
		Use:     "create",
		Short:   "create connection between two assets.",
		Long:    `The connection create command links two assets together.  The outgoing asset is the asset that wants to connect to another asset.  The incoming asset is the asset that the outgoing asset connects to.`,
		Aliases: []string{"c", "deploy"},
		RunE:    connCreateRun(),
	}

	connCreateCmd.Flags().StringVarP(&connOptions.OutAsset, "outgoing-asset", "", "", "The outgoing asset is the asset that wants to connect to another asset.")
	connCreateCmd.Flags().StringVarP(&connOptions.InAsset, "incoming-asset", "", "", "The incoming asset is the asset that the outgoing asset connects to.")
	connCreateCmd.Flags().StringVarP(&connOptions.Description, "description", "", "", "Describe the connection")

	connCmd.AddCommand(connCreateCmd)

	return connCmd
}
