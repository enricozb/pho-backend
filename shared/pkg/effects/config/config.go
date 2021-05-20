package config

import (
	"bytes"
	_ "embed"
	"fmt"

	"github.com/spf13/viper"
)

type PhoConfig struct {
	*viper.Viper

	DataPath string `mapstructure:"data_path"`
}

var Config = &PhoConfig{Viper: viper.New()}

//go:embed config_defaults.json
var defaultConfig []byte

func init() {
	Config.SetConfigType("json")
	Config.ReadConfig(bytes.NewReader(defaultConfig))

	// user-defined config
	Config.SetConfigName("config")
	Config.AddConfigPath("$HOME/.config/pho/")
	Config.AddConfigPath("/etc/pho/")

	// panic if an error other than ConfigFileNotFound occurs
	err := Config.MergeInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
		panic(fmt.Errorf("merge in config: %s", err))
	}

	if err := Config.Unmarshal(&Config); err != nil {
		panic(fmt.Errorf("unmarshal config: %v", err))
	}
}
