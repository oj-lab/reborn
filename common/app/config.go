package app

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

const (
	EnvMode           = "MODE"
	defaultConfigName = "default"
)

var (
	config *viper.Viper
)

func Config() *viper.Viper {
	return config
}

func initConfig(configPath string) {
	viper.AddConfigPath(configPath)
	viper.SetConfigName(defaultConfigName)
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(err)
		}
	}

	viper.SetConfigName(os.Getenv(EnvMode))
	err = viper.MergeInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(err)
		}
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	config = viper.GetViper()
}
