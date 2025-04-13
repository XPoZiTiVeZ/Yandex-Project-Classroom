package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Host struct {
		Address string `mapstructure:"Address"`
		Port    int    `mapstructure:"Port"`
	} `mapstructure:"Host"`
	Auth struct {
		Address string `mapstructure:"Address"`
		Port    int    `mapstructure:"Port"`
		Enabled bool   `mapstructure:"Enabled"`
	} `mapstructure:"Auth"`
	Courses struct {
		Address string `mapstructure:"Address"`
		Port    int    `mapstructure:"Port"`
		Enabled bool   `mapstructure:"Enabled"`
	} `mapstructure:"Courses"`
	Lessons struct {
		Address string `mapstructure:"Address"`
		Port    int    `mapstructure:"Port"`
		Enabled bool   `mapstructure:"Enabled"`
	} `mapstructure:"Lessons"`
	Tasks struct {
		Address string `mapstructure:"Address"`
		Port    int    `mapstructure:"Port"`
		Enabled bool   `mapstructure:"Enabled"`
	} `mapstructure:"Tasks"`

	Env struct {
		AuthJWTSecret string `mapstructure:"AUTH_JWT_SECRET"`
	}
}

func ReadConfig() (Config, error) {
	viper.SetDefault("Gateway.Address", "127.0.0.1")
	viper.SetDefault("Gateway.Port", 8000)
	
	viper.AddConfigPath("./config")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("..")
	
	var AppConfig Config
	
	viper.SetConfigFile("./config/config.yaml")
	viper.SetConfigFile("config.yaml")
	
	err := viper.ReadInConfig()
	parser := viper.Sub("API-Gateway")
	if err != nil {
		return Config{}, err
	}
	if err := parser.Unmarshal(&AppConfig); err != nil {
		return Config{}, err
	}

	viper.SetConfigType("env")
	viper.SetConfigFile(".env")
	// viper.SetConfigFile("../.env")

	err = viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}
	if err := viper.Unmarshal(&(AppConfig.Env)); err != nil {
		return Config{}, err
	}

	return AppConfig, err
}
