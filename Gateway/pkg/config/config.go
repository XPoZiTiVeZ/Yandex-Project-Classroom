package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Address string `mapstructure:"address"`
		Port    int    `mapstructure:"port"`
	} `mapstructure:"gateway"`
	Auth struct {
		Address string `mapstructure:"address"`
		Port    int    `mapstructure:"port"`
	} `mapstructure:"auth"`
	Courses struct {
		Address string `mapstructure:"Address"`
		Port    int    `mapstructure:"Port"`
	} `mapstructure:"Courses"`
	Lessons struct {
		Address string `mapstructure:"Address"`
		Port    int    `mapstructure:"Port"`
	} `mapstructure:"Lessons"`

	Env struct {
		AuthJWTSecret string `mapstructure:"AUTH_JWT_SECRET"`
	}
}

func ReadConfig() (Config, error) {
	viper.SetDefault("Gateway.Address", "127.0.0.1")
	viper.SetDefault("Gateway.Port", 8000)

	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	var AppConfig Config

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}
	if err := viper.Unmarshal(&AppConfig); err != nil {
		return Config{}, err
	}

	viper.SetConfigType("env")
	viper.SetConfigFile("../.env")

	err = viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}
	if err := viper.Unmarshal(&(AppConfig.Env)); err != nil {
		return Config{}, err
	}

	return AppConfig, err
}
