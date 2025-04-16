package config

import (
	"flag"

	"github.com/spf13/viper"
)

type Config struct {
	Host struct {
		Address string `mapstructure:"address"`
		Port    int    `mapstructure:"port"`
	} `mapstructure:"host"`
	Auth struct {
		Address string `mapstructure:"address"`
		Port    int    `mapstructure:"port"`
		Enabled bool   `mapstructure:"enabled"`
	} `mapstructure:"auth"`
	Courses struct {
		Address string `mapstructure:"address"`
		Port    int    `mapstructure:"port"`
		Enabled bool   `mapstructure:"enabled"`
	} `mapstructure:"courses"`
	Lessons struct {
		Address string `mapstructure:"address"`
		Port    int    `mapstructure:"port"`
		Enabled bool   `mapstructure:"enabled"`
	} `mapstructure:"lessons"`
	Tasks struct {
		Address string `mapstructure:"address"`
		Port    int    `mapstructure:"port"`
		Enabled bool   `mapstructure:"enabled"`
	} `mapstructure:"tasks"`

	Env struct {
		RedisURL      string `mapstructure:"REDIS_URL"`
		KafkaURL      string `mapstructure:"KAFKA_URL"`
		AuthJWTSecret string `mapstructure:"AUTH_JWT_SECRET"`
	}
}

func MustReadConfig() Config {
	configPath := flag.String("config", "./config/config.yaml", "path to config file")
	flag.Parse()

	viper.SetDefault("Gateway.Address", "127.0.0.1")
	viper.SetDefault("Gateway.Port", 8000)

	var AppConfig Config

	viper.SetConfigFile(*configPath)

	err := viper.ReadInConfig()
	if err != nil {
		return AppConfig
	}
	if err := viper.Unmarshal(&AppConfig); err != nil {
		return AppConfig
	}

	viper.AutomaticEnv()
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		return AppConfig
	}
	if err := viper.Unmarshal(&(AppConfig.Env)); err != nil {
		return AppConfig
	}

	return AppConfig
}
