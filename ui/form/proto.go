package form

import (
	"github.com/aptible/cloud-cli/config"
	"github.com/charmbracelet/bubbles/list"
)

type FormResult struct {
	Org           string
	Env           string
	AssetType     string
	AssetName     string
	VpcName       string
	Engine        string
	EngineVersion string
}

type FormFn func(*config.CloudConfig, *FormResult) error

type FormOption struct {
	Value string
	Label string
}

func (i FormOption) Title() string       { return i.Label }
func (i FormOption) Description() string { return i.Value }
func (i FormOption) FilterValue() string { return i.Label }

type LoadOptionsFn func(cfg *config.CloudConfig) ([]list.Item, error)

type SubSchema struct {
	Title       string
	Type        string
	LoadOptions LoadOptionsFn
	Err         error
}
