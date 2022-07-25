package form

import (
	"github.com/aptible/cloud-cli/internal/common"
	"github.com/charmbracelet/bubbles/list"
)

type FormResult struct {
	Org string
	Env string
}

type FormOption struct {
	Value string
	Label string
}

func (i FormOption) Title() string       { return i.Label }
func (i FormOption) Description() string { return i.Value }
func (i FormOption) FilterValue() string { return i.Label }

type LoadOptionsFn func(config *common.CloudConfig) ([]list.Item, error)

type SubSchema struct {
	Title       string
	Type        string
	LoadOptions LoadOptionsFn
	Err         error
}
