package cmd

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"path"

	apiclient "github.com/aptible/cloud-api-clients/clients/go"
	client "github.com/aptible/cloud-cli/client"
	proto "github.com/aptible/cloud-cli/proto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ctxCloudConfig struct{}

func GetCloudConfig(cmd *cobra.Command) *proto.CloudConfig {
	return cmd.Context().Value(ctxCloudConfig{}) //.Value(ctxCloudConfig{}).(*proto.CloudConfig)
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
	ctx = context.WithValue(ctx, apiclient.ContextAccessToken, token)
	return ctx
}

func NewCloudConfig(v *viper.Viper) *proto.CloudConfig {
	conf := apiclient.NewConfiguration()
	apiClient := apiclient.NewAPIClient(conf)
	token := v.GetString("token")
	ctx := NewContext(token)
	cc := client.NewClient(ctx, apiClient)

	return &proto.CloudConfig{
		Vconfig: v,
		Cc:      cc,
		Ctx:     ctx,
	}
}
