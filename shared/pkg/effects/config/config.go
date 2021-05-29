package config

import (
	"bytes"
	_ "embed"
	"fmt"
	"reflect"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"github.com/enricozb/pho/shared/pkg/lib/file"
)

type PhoConfig struct {
	*viper.Viper

	// DBDir is the directory containing the pho db.
	DBDir string `config:"dir" mapstructure:"db_dir"`

	// DataDir is the directory containing all media data.
	DataDir string `config:"dir" mapstructure:"data_dir"`
}

var Config = &PhoConfig{Viper: viper.New()}

//go:embed config_defaults.json
var defaultConfig []byte

func init() {
	Config.SetConfigType("json")
	if err := Config.ReadConfig(bytes.NewReader(defaultConfig)); err != nil {
		panic("read default config")
	}

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

	if err := initDirs(); err != nil {
		panic(fmt.Errorf("init data path: %v", err))
	}
}

func initDirs() error {
	t := reflect.TypeOf(PhoConfig{})
	v := reflect.ValueOf(Config)

	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Tag.Get("config") == "dir" {
			dir, err := homedir.Expand(v.Elem().Field(i).String())
			if err != nil {
				return fmt.Errorf("expand: %v", err)
			}
			v.Elem().Field(i).SetString(dir)

			if err := file.MakeDirIfNotExist(dir); err != nil {
				return fmt.Errorf("make dir: %v", err)
			}
		}
	}

	return nil
}
