package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address         string `env:"ADDRESS,required"`
	StoreInterval   int    `env:"STORE_INTERVAL,required"`
	FileStoragePath string `env:"FILE_STORAGE_PATH,required"`
	Restore         bool   `env:"RESTORE,required"`
}

func NewConfig() *Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Printf("Error parse configuration parameters:\n%s", err.Error())
		return nil
	}
	fmt.Println(cfg)
	return &cfg
}

type ConfigAgent struct {
	Address        string `env:"ADDRESS,required"`
	ReportInterval int    `env:"REPORT_INTERVAL,required"`
	PollInterval   int    `env:"POLL_INTERVAL,required"`
}

func NewConfigAgent() *ConfigAgent {
	var cfgAgent ConfigAgent
	err := env.Parse(&cfgAgent)
	if err != nil {
		log.Printf("Error parse configuration parameters:\n%s", err.Error())
		return nil
	}
	return &cfgAgent
}
