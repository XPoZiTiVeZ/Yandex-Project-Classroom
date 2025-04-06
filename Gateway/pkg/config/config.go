package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Address string `mapstructure:"Address"`
		Port    int    `mapstructure:"Port"`
	} `mapstructure:"Gateway"`
	Auth struct {
		Address string `mapstructure:"Address"`
		Port    int    `mapstructure:"Port"`
	} `mapstructure:"Auth"`
	Courses struct {
		Address string `mapstructure:"Address"`
		Port    int    `mapstructure:"Port"`
	} `mapstructure:"Courses"`
	Lessons struct {
		Address string `mapstructure:"Address"`
		Port    int    `mapstructure:"Port"`
	} `mapstructure:"Lessons"`
}

func ReadConfig() (Config, error) {
	viper.SetDefault("Gateway.Address", "127.0.0.1")
	viper.SetDefault("Gateway.Port", 8000)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")

	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	var AppConfig Config
	if err := viper.Unmarshal(&AppConfig); err != nil {
		return Config{}, err
	}

	return AppConfig, err
}
