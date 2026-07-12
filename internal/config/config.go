package config

import (
	"log"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	address string `env:"ADDRESS"`
}

func NewConfig() *Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Printf("Error parse configuration parameters:\n%s", err.Error())
		return nil
	}
	return &cfg
}

func (c Config) Address() string {
	return c.address
}
