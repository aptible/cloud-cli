package assets

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/aptible/cloud-cli/internal/common"
	"github.com/aptible/cloud-cli/internal/ui/fetch"
)

func destroyAsset(_ *cobra.Command, args []string) error {
	config := common.NewCloudConfig(viper.GetViper())
	orgId := config.Vconfig.GetString("org")
	assetId := args[0]

	if env == "" {
		return fmt.Errorf("must provide env")
	}

	msg := fmt.Sprintf("destroying asset %s (v%s)", engine, engineVersion)
	model := fetch.NewModel(msg, func() (interface{}, int, error) {
		status, err := config.Cc.DestroyAsset(orgId, env, assetId)
		return nil, status, err
	})
	_, err := fetch.WithOutput(model)
	if err != nil {
		return err
	}

	fmt.Printf("destroying asset with id: %+v\n", assetId)
	return nil
}
