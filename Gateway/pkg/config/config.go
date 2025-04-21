package config

import (
	"flag"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Host struct {
		Address string `mapstructure:"address"`
		Port    int    `mapstructure:"port"`
	} `mapstructure:"gateway"`
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
	Chat struct {
		Address string `mapstructure:"address"`
		Port    int    `mapstructure:"port"`
		Enabled bool   `mapstructure:"enabled"`
	} `mapstructure:"chat"`
	Notifications struct {
		Enabled bool `mapstructure:"enabled"`
	} `mapstructure:"notifications"`

	Common struct {
		Timeout       time.Duration `mapstructure:"timeout"`
		MaxRetries    int           `mapstructure:"max_retries"`
		AuthJWTSecret string        `mapstructure:"auth_jwt_secret"`
		RedisURL      string        `mapstructure:"redis_url"`
		KafkaURL      string        `mapstructure:"kafka_url"`
	} `mapstructure:"common"`
}

const (
	DefaultGatewayPort int           = 8080
	DefaultGatewayAddr string        = "0.0.0.0"
	DefaultTimeout     time.Duration = 10 * time.Second
	MaxRetries         int           = 5
)

var (
	DefaultConfig Config = func() Config {
		var Config Config
		Config.Host.Port = DefaultGatewayPort

		return Config
	}()
)

func MustReadConfig() *Config {
	configPath := flag.String("config", "./configs/config.yaml", "path to config file")
	flag.Parse()

	v := viper.New()
	v.SetDefault("gateway.port", DefaultGatewayPort)
	v.SetDefault("gateway.address", DefaultGatewayAddr)
	v.SetDefault("common.timeout", DefaultTimeout)
	v.SetDefault("common.max_retries", MaxRetries)

	var AppConfig Config

	viper.SetConfigFile(*configPath)

	err := viper.ReadInConfig()
	if err != nil {
		slog.Debug("Couldn't read config", slog.Any("error", err))
		return &DefaultConfig
	}
	if err := viper.Unmarshal(&AppConfig); err != nil {
		slog.Debug("Couldn't unmarshal config", slog.Any("error", err))
		return &DefaultConfig
	}

	v.AutomaticEnv()
	// v.SetConfigFile(".env")
	// v.SetConfigType("env")

	// if err := v.ReadInConfig(); err != nil {
	// 	slog.Debug("Couldn't read .env", slog.Any("error", err))
	// 	return &DefaultConfig
	// }

	envToMapstructure := map[string]string{
		"GATEWAY_PORT":    "gateway.port",
		"GATEWAY_ADDRESS": "gateway.address",
		"TIMEOUT":         "common.timeout",
		"MAX_RETRIES":     "common.max_retries",
		"AUTH_JWT_SECRET": "common.auth_jwt_secret",
		"REDIS_URL":       "common.redis_url",
		"KAFKA_URL":       "common.kafka_url",
	}

	serviceList := []string{
		"auth",
		"courses",
		"lessons",
		"tasks",
		"chat",
	}

	for _, service := range serviceList {
		for _, field := range []string{"port", "address", "enabled"} {
			key := fmt.Sprintf("%s_%s", strings.ToUpper(service), strings.ToUpper(field))
			value := fmt.Sprintf("%s.%s", service, field)
			envToMapstructure[key] = value
		}
	}

	for envKey, mapKey := range envToMapstructure {
		if v.IsSet(envKey) {
			v.Set(mapKey, v.Get(envKey))
		}
	}

	if err := v.Unmarshal(&AppConfig); err != nil {
		slog.Debug("Couldn't unmarshal config", slog.Any("error", err))
		return &DefaultConfig
	}

	return &AppConfig
}
