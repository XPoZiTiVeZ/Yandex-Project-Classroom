package config

import (
	"flag"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	PostgresURL string `mapstructure:"postgres_url"`
	KafkaBroker string `mapstructure:"kafka_broker"`
	SMTP        SMTP   `mapstructure:"smtp"`
}

type SMTP struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	User string `mapstructure:"user"`
	Pass string `mapstructure:"password"`
}

func MustNew() *Config {
	configPath := flag.String("config", "./config/config.yaml", "path to config file")
	flag.Parse()

	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.BindEnv("kafka_broker")
	v.BindEnv("postgres_url")
	v.BindEnv("smtp.host")
	v.BindEnv("smtp.port")
	v.BindEnv("smtp.user")
	v.BindEnv("smtp.password")

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
