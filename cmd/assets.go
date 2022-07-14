package cmd

import (
	"fmt"

	"github.com/aptible/cloud-cli/ui/fetch"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func destroyAsset(_ *cobra.Command, args []string) error {
	config := NewCloudConfig(viper.GetViper())
	orgId := config.Vconfig.GetString("org")
	assetId := args[0]

	if env == "" {
		return fmt.Errorf("must provide env")
	}

	msg := fmt.Sprintf("destroying asset %s (v%s)", engine, engineVersion)
	model := fetch.NewModel(msg, func() (interface{}, error) {
		return nil, config.Cc.DestroyAsset(orgId, env, assetId)
	})
	_, err := fetch.FetchWithOutput(model)
	if err != nil {
		return err
	}

	fmt.Printf("destroying asset with id: %+v\n", assetId)
	return nil
}
