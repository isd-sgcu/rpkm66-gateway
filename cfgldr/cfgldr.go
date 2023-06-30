package cfgldr

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Service struct {
	Auth    string `mapstructure:"auth"`
	Backend string `mapstructure:"backend"`
	File    string `mapstructure:"file"`
}

type App struct {
	Port        int    `mapstructure:"port"`
	Debug       bool   `mapstructure:"debug"`
	Phase       string `mapstructure:"phase"`
	MaxFileSize int    `mapstructure:"max_file_size"`
}

type Config struct {
	Service Service `mapstructure:"service"`
	App     App     `mapstructure:"app"`
}

func LoadConfig() (config *Config, err error) {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return nil, errors.Wrap(err, "error occurs while reading the config")
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, errors.Wrap(err, "error occurs while unmarshal the config")
	}

	return
}
