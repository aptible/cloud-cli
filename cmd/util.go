package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"

	apiclient "github.com/aptible/cloud-api-clients/clients/go"
	client "github.com/aptible/cloud-cli/client"
	proto "github.com/aptible/cloud-cli/proto"
	"github.com/aptible/cloud-cli/ui/fetch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ctxCloudConfig struct{}

func GetCloudConfig(cmd *cobra.Command) *proto.CloudConfig {
	val := cmd.Context().Value(ctxCloudConfig{})
	if val == nil {
		return nil
	}
	return val.(*proto.CloudConfig)
}

type RunE func(cmd *cobra.Command, args []string) error

func findToken(home string, domain string) string {
	var tokenObj map[string]string
	text, err := ioutil.ReadFile(path.Join(home, ".aptible", "tokens.json"))
	if err != nil {
		return ""
	}
	json.Unmarshal(text, &tokenObj)
	return string(tokenObj[domain])
}

func NewContext(token string) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, apiclient.ContextAccessToken, fmt.Sprintf("Bearer %s", token))
	return ctx
}

func NewCloudConfig(v *viper.Viper) *proto.CloudConfig {
	conf := apiclient.NewConfiguration()
	host := v.GetString("api-domain")
	conf.Host = host
	conf.Scheme = "http"
	apiClient := apiclient.NewAPIClient(conf)
	token := v.GetString("token")
	fmt.Println(token)
	ctx := NewContext(token)
	debug := v.GetBool("debug")
	cc := client.NewClient(ctx, apiClient, debug)

	return &proto.CloudConfig{
		Vconfig: v,
		Cc:      cc,
		Ctx:     ctx,
	}
}

func FetchWithOutput[T interface{}](model tea.Model, result *T) error {
	p := tea.NewProgram(model)
	m, err := p.StartReturningModel()
	if err != nil {
		return err
	}

	n := m.(fetch.Model)
	switch n.Result.(type) {
	case T:
		r := n.Result.(T)
		result = &r
	default:
		result = nil
	}

	return nil
}
