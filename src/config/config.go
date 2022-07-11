package config

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
	Port        int  `mapstructure:"port"`
	Debug       bool `mapstructure:"debug"`
	MaxFileSize int  `mapstructure:"max_file_size"`
}

type Guard struct {
	Phase string `mapstructure:"phase"`
}

type Config struct {
	Service Service `mapstructure:"service"`
	App     App     `mapstructure:"app"`
	Guard   Guard   `mapstructure:"guard"`
}

func LoadConfig() (config *Config, err error) {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

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
