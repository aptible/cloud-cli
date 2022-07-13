package proto

import (
	"context"
	"github.com/spf13/viper"
)

type CloudConfig struct {
	Vconfig *viper.Viper
	Cc      CloudClient
	Ctx     context.Context
}
