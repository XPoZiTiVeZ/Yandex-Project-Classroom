package config

import (
	"flag"
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Port        int    `mapstructure:"port"`
	PostgresURL string `mapstructure:"postgres_url"`
	RedisURL    string `mapstructure:"redis_url"`
	Auth        Auth   `mapstructure:"auth"`
}

type Auth struct {
	JwtSecret  string        `mapstructure:"jwt_secret"`
	AccessTTL  time.Duration `mapstructure:"access_ttl"`
	RefreshTTL time.Duration `mapstructure:"refresh_ttl"`
}

func MustNew() *Config {
	configPath := flag.String("config", "./config/config.yaml", "path to config file")
	flag.Parse()

	v := viper.New()

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.SetConfigFile(*configPath)

	if err := v.ReadInConfig(); err != nil {
		log.Printf("no config file found: %v (continuing with env vars only)", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("unable to decode config: %v", err)
	}

	return &cfg
}
