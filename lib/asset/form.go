package libasset

import (
	"fmt"

	"github.com/aptible/cloud-cli/config"
	"github.com/aptible/cloud-cli/lib/env"
	"github.com/aptible/cloud-cli/ui/form"
	"github.com/charmbracelet/bubbles/list"
)

func NewAssetNameProp() *form.SubSchema {
	return &form.SubSchema{
		Type:  "input",
		Title: "What do you want to call the asset?",
	}
}

func CreateAssetTypeOptions(orgId, envId string) form.LoadOptionsFn {
	options := []list.Item{}
	return func(cfg *config.CloudConfig) ([]list.Item, error) {
		bundles, err := cfg.Cc.ListAssetBundles(orgId, envId)
		if err != nil {
			return options, err
		}
		for _, bundle := range bundles {
			options = append(options, form.FormOption{Label: bundle.Name, Value: bundle.Identifier})
		}
		return options, nil
	}
}

func NewAssetTypeProp(orgId, envId string) *form.SubSchema {
	return &form.SubSchema{
		Type:        "select",
		Title:       "Select an asset type",
		LoadOptions: CreateAssetTypeOptions(orgId, envId),
	}
}

func CreateVPCOptions(orgId, envId string) form.LoadOptionsFn {
	options := []list.Item{}
	return func(cfg *config.CloudConfig) ([]list.Item, error) {
		assets, err := cfg.Cc.ListAssets(orgId, envId)
		if err != nil {
			return options, err
		}
		filtered := FilterByType(assets, []string{"vpc"})
		for _, asset := range filtered {
			name := GetName(asset)
			options = append(options, form.FormOption{Label: name, Value: name})
		}
		return options, nil
	}
}

func NewVpcNameProp(orgId, envId string) *form.SubSchema {
	return &form.SubSchema{
		Type:        "select",
		Title:       "Select a VPC",
		LoadOptions: CreateVPCOptions(orgId, envId),
	}
}

func CreateEngineOptions() form.LoadOptionsFn {
	options := []list.Item{}
	return func(cfg *config.CloudConfig) ([]list.Item, error) {
		options = append(options, form.FormOption{Label: "postgres", Value: "postgres"})
		options = append(options, form.FormOption{Label: "mysql", Value: "mysql"})
		return options, nil
	}
}

func NewEngineProp() *form.SubSchema {
	return &form.SubSchema{
		Type:        "select",
		Title:       "Select an engine",
		LoadOptions: CreateEngineOptions(),
	}
}

func CreateEngineVersionOptions(engine string) form.LoadOptionsFn {
	engineMap := map[string][]string{
		"postgres": {"13", "14"},
		"mysql":    {"8.0", "5.7"},
	}
	var options []list.Item
	return func(cfg *config.CloudConfig) ([]list.Item, error) {
		versions := engineMap[engine]
		fmt.Println(versions)
		fmt.Println(engine)
		for _, option := range versions {
			options = append(options, form.FormOption{Label: option, Value: option})
		}
		return options, nil
	}
}

func NewEngineVersionProp(engine string) *form.SubSchema {
	return &form.SubSchema{
		Type:        "select",
		Title:       "Select an engine version",
		LoadOptions: CreateEngineVersionOptions(engine),
	}
}

func AssetNameForm(cfg *config.CloudConfig, results *form.FormResult) error {
	if results.AssetName != "" {
		return nil
	}

	prop := NewAssetNameProp()
	result, err := form.Run(form.NewModel(cfg, prop))
	if err != nil {
		return err
	}
	if result == "" {
		return fmt.Errorf("You must enter a name for your asset")
	}
	results.AssetName = result

	return nil
}

func AssetTypeForm(cfg *config.CloudConfig, results *form.FormResult) error {
	if results.AssetType != "" {
		return nil
	}

	prop := NewAssetTypeProp(results.Org, results.Env)
	result, err := form.Run(form.NewModel(cfg, prop))
	if err != nil {
		return err
	}
	if result == "" {
		return fmt.Errorf("You must enter a type for your asset")
	}
	results.AssetType = result

	return nil
}

func AssetVpcNameForm(cfg *config.CloudConfig, results *form.FormResult) error {
	if results.VpcName != "" {
		return nil
	}

	prop := NewVpcNameProp(results.Org, results.Env)
	result, err := form.Run(form.NewModel(cfg, prop))
	if err != nil {
		return err
	}
	if result == "" {
		return fmt.Errorf("You must enter a vpc for your asset")
	}
	results.VpcName = result

	return nil
}

func AssetEngineForm(cfg *config.CloudConfig, results *form.FormResult) error {
	if results.Engine != "" {
		return nil
	}

	prop := NewEngineProp()
	result, err := form.Run(form.NewModel(cfg, prop))
	if err != nil {
		return err
	}
	if result == "" {
		return fmt.Errorf("You must select an engine for your asset")
	}
	results.Engine = result

	return nil
}

func AssetEngineVersionForm(cfg *config.CloudConfig, results *form.FormResult) error {
	if results.EngineVersion != "" {
		return nil
	}

	prop := NewEngineVersionProp(results.Engine)
	result, err := form.Run(form.NewModel(cfg, prop))
	if err != nil {
		return err
	}
	if result == "" {
		return fmt.Errorf("You must select an engine version for your asset")
	}
	results.EngineVersion = result

	return nil
}

func CreateAssetOptions(orgId, envId string) form.LoadOptionsFn {
	options := []list.Item{}
	return func(cfg *config.CloudConfig) ([]list.Item, error) {
		assets, err := cfg.Cc.ListAssets(orgId, envId)
		if err != nil {
			return options, err
		}
		for _, asset := range assets {
			name := GetName(asset)
			options = append(options, form.FormOption{Label: name, Value: asset.Id})
		}
		return options, nil
	}
}

func NewAssetProp(orgId, envId string) *form.SubSchema {
	return &form.SubSchema{
		Type:        "select",
		Title:       "Select an asset",
		LoadOptions: CreateAssetOptions(orgId, envId),
	}
}

func AssetForm(cfg *config.CloudConfig, results *form.FormResult) error {
	if results.Asset != "" {
		return nil
	}

	prop := NewAssetProp(results.Org, results.Env)
	result, err := form.Run(form.NewModel(cfg, prop))
	if err != nil {
		return err
	}
	if result == "" {
		return fmt.Errorf("You must select an asset")
	}
	results.Asset = result

	return nil
}

func AssetDescribeForm(cfg *config.CloudConfig, results *form.FormResult) error {
	forms := []form.FormFn{
		libenv.EnvForm,
		AssetForm,
	}

	for _, formFn := range forms {
		err := formFn(cfg, results)
		if err != nil {
			return err
		}
	}

	return nil
}

func AssetCreateForm(cfg *config.CloudConfig, results *form.FormResult) error {
	forms := []form.FormFn{
		libenv.EnvForm,
		AssetVpcNameForm,
		AssetTypeForm,
		AssetEngineForm,
		AssetEngineVersionForm,
		AssetNameForm,
	}

	for _, formFn := range forms {
		err := formFn(cfg, results)
		if err != nil {
			return err
		}
	}

	return nil
}
