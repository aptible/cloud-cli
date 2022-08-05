package libconn

import (
	"fmt"

	"github.com/aptible/cloud-cli/config"
	libasset "github.com/aptible/cloud-cli/lib/asset"
	libenv "github.com/aptible/cloud-cli/lib/env"
	"github.com/aptible/cloud-cli/ui/form"
)

func OutAssetForm(cfg *config.CloudConfig, results *form.FormResult) error {
	if results.OutAsset != "" {
		return nil
	}

	prop := libasset.NewAssetProp(results.Org, results.Env)
	prop.Title = "Select an outgoing asset (from)"
	result, err := form.Run(form.NewModel(cfg, prop))
	if err != nil {
		return err
	}
	if result == "" {
		return fmt.Errorf("You must select an outgoing asset in order to create a connection")
	}
	results.OutAsset = result

	return nil
}

func InAssetForm(cfg *config.CloudConfig, results *form.FormResult) error {
	if results.InAsset != "" {
		return nil
	}

	prop := libasset.NewAssetProp(results.Org, results.Env)
	prop.Title = "Select an incoming asset (to)"
	result, err := form.Run(form.NewModel(cfg, prop))
	if err != nil {
		return err
	}
	if result == "" {
		return fmt.Errorf("You must select an incoming asset in order to create a connection")
	}
	results.InAsset = result

	return nil
}

func NewDescProp() *form.SubSchema {
	return &form.SubSchema{
		Type:  "input",
		Title: "Describe the asset connection",
	}
}

func DescForm(cfg *config.CloudConfig, results *form.FormResult) error {
	if results.Description != "" {
		return nil
	}

	prop := NewDescProp()
	result, err := form.Run(form.NewModel(cfg, prop))
	if err != nil {
		return err
	}
	if result == "" {
		return fmt.Errorf("You must enter a description for the asset connection")
	}
	results.Description = result

	return nil
}

func ConnCreateForm(cfg *config.CloudConfig, results *form.FormResult) error {
	forms := []form.FormFn{
		libenv.EnvForm,
		OutAssetForm,
		InAssetForm,
		DescForm,
	}

	for _, formFn := range forms {
		err := formFn(cfg, results)
		if err != nil {
			return err
		}
	}

	return nil
}
