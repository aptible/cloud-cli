package form

import (
	"fmt"

	"github.com/aptible/cloud-cli/internal/common"
	uiAsset "github.com/aptible/cloud-cli/internal/ui/asset"
	"github.com/charmbracelet/bubbles/list"
)

func NewAssetNameProp() *SubSchema {
	return &SubSchema{
		Type:  "input",
		Title: "What do you want to call the asset?",
	}
}

func CreateAssetTypeOptions(orgId, envId string) LoadOptionsFn {
	options := []list.Item{}
	return func(config *common.CloudConfig) ([]list.Item, error) {
		bundles, err := config.Cc.ListAssetBundles(orgId, envId)
		if err != nil {
			return options, err
		}
		for _, bundle := range bundles {
			options = append(options, FormOption{Label: bundle.Name, Value: bundle.Identifier})
		}
		return options, nil
	}
}

func NewAssetTypeProp(orgId, envId string) *SubSchema {
	return &SubSchema{
		Type:        "select",
		Title:       "Select an asset type",
		LoadOptions: CreateAssetTypeOptions(orgId, envId),
	}
}

func CreateVPCOptions(orgId, envId string) LoadOptionsFn {
	options := []list.Item{}
	return func(config *common.CloudConfig) ([]list.Item, error) {
		assets, err := config.Cc.ListAssets(orgId, envId)
		if err != nil {
			return options, err
		}
		filtered := common.FilterAssetsByType(assets, []string{"vpc"})
		for _, asset := range filtered {
			name := uiAsset.GetAssetName(&asset)
			options = append(options, FormOption{Label: name, Value: name})
		}
		return options, nil
	}
}

func NewVpcNameProp(orgId, envId string) *SubSchema {
	return &SubSchema{
		Type:        "select",
		Title:       "Select a VPC",
		LoadOptions: CreateVPCOptions(orgId, envId),
	}
}

func CreateEngineOptions() LoadOptionsFn {
	options := []list.Item{}
	return func(config *common.CloudConfig) ([]list.Item, error) {
		options = append(options, FormOption{Label: "postgres", Value: "postgres"})
		options = append(options, FormOption{Label: "mysql", Value: "mysql"})
		return options, nil
	}
}

func NewEngineProp() *SubSchema {
	return &SubSchema{
		Type:        "select",
		Title:       "Select an engine",
		LoadOptions: CreateEngineOptions(),
	}
}

func CreateEngineVersionOptions(engine string) LoadOptionsFn {
	engineMap := map[string][]string{
		"postgres": {"13", "14"},
		"mysql":    {"8.0", "5.7"},
	}
	var options []list.Item
	return func(config *common.CloudConfig) ([]list.Item, error) {
		versions := engineMap[engine]
		fmt.Println(versions)
		fmt.Println(engine)
		for _, option := range versions {
			options = append(options, FormOption{Label: option, Value: option})
		}
		return options, nil
	}
}

func NewEngineVersionProp(engine string) *SubSchema {
	return &SubSchema{
		Type:        "select",
		Title:       "Select an engine version",
		LoadOptions: CreateEngineVersionOptions(engine),
	}
}

func AssetNameForm(config *common.CloudConfig, results *FormResult) error {
	if results.AssetName != "" {
		return nil
	}

	prop := NewAssetNameProp()
	result, err := Run(NewModel(config, prop))
	if err != nil {
		return err
	}
	if result == "" {
		return fmt.Errorf("You must enter a name for your asset")
	}
	results.AssetName = result

	return nil
}

func AssetTypeForm(config *common.CloudConfig, results *FormResult) error {
	if results.AssetType != "" {
		return nil
	}

	prop := NewAssetTypeProp(results.Org, results.Env)
	result, err := Run(NewModel(config, prop))
	if err != nil {
		return err
	}
	if result == "" {
		return fmt.Errorf("You must enter a type for your asset")
	}
	results.AssetType = result

	return nil
}

func AssetVpcNameForm(config *common.CloudConfig, results *FormResult) error {
	if results.VpcName != "" {
		return nil
	}

	prop := NewVpcNameProp(results.Org, results.Env)
	result, err := Run(NewModel(config, prop))
	if err != nil {
		return err
	}
	if result == "" {
		return fmt.Errorf("You must enter a vpc for your asset")
	}
	results.VpcName = result

	return nil
}

func AssetEngineForm(config *common.CloudConfig, results *FormResult) error {
	if results.Engine != "" {
		return nil
	}

	prop := NewEngineProp()
	result, err := Run(NewModel(config, prop))
	if err != nil {
		return err
	}
	if result == "" {
		return fmt.Errorf("You must select an engine for your asset")
	}
	results.Engine = result

	return nil
}

func AssetEngineVersionForm(config *common.CloudConfig, results *FormResult) error {
	if results.EngineVersion != "" {
		return nil
	}

	prop := NewEngineVersionProp(results.Engine)
	result, err := Run(NewModel(config, prop))
	if err != nil {
		return err
	}
	if result == "" {
		return fmt.Errorf("You must select an engine version for your asset")
	}
	results.EngineVersion = result

	return nil
}

func AssetCreateForm(config *common.CloudConfig, results *FormResult) error {
	forms := []FormFn{
		EnvForm,
		AssetVpcNameForm,
		AssetTypeForm,
		AssetEngineForm,
		AssetEngineVersionForm,
		AssetNameForm,
	}

	for _, formFn := range forms {
		err := formFn(config, results)
		if err != nil {
			return err
		}
	}

	return nil
}
