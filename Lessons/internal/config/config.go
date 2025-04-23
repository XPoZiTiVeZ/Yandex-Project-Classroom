package config

import (
	"flag"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port        int    `mapstructure:"port"`
	PostgresURL string `mapstructure:"postgres_url"`
	KafkaBroker string `mapstructure:"kafka_broker"`
}

func MustNew() *Config {
	configPath := flag.String("config", "./config/config.yaml", "path to config file")
	flag.Parse()

	v := viper.New()

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.BindEnv("postgres_url")
	v.BindEnv("kafka_broker")
	v.BindEnv("port")

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
