package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	Address     string `env:"RUN_ADDRESS"`
	LogLevel    string `env:"LOG_LEVEL"`
	DatabaseDSN string `env:"DATABASE_URI"`
}

func ParseConfig() *Config {
	var config Config
	flag.StringVar(&config.Address, "a", ":8080", "address and port to run server")
	flag.StringVar(&config.LogLevel, "l", "info", "log level")
	flag.StringVar(&config.DatabaseDSN, "d", "", "database connection string")
	flag.Parse()

	if err := env.Parse(&config); err != nil {
		log.Fatal(err)
	}

	return &config
}
