package app

import (
	"os"
	"path"
	"strings"

	"github.com/spf13/viper"
)

const (
	envRelativeConfigPath = "RELATIVE_CONFIG_PATH"
	envOverrideConfigName = "OVERRIDE_CONFIG_NAME"

	defaultConfigName         = "common"
	defaultOverrideConfigName = "override.config"
)

var (
	config *viper.Viper
)

func Config() *viper.Viper {
	return config
}

func init() {
	cwd, _ := os.Getwd()
	relativeConfigPath := os.Getenv(envRelativeConfigPath)
	viper.AddConfigPath(path.Join(cwd, "configs", relativeConfigPath))
	viper.SetConfigName(defaultConfigName)
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(err)
		}
	}
	overrideConfigName := os.Getenv(envOverrideConfigName)
	if len(overrideConfigName) == 0 {
		overrideConfigName = defaultOverrideConfigName
	}
	viper.SetConfigName(overrideConfigName)
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
