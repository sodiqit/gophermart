package config

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	Address      string `env:"RUN_ADDRESS"`
	LogLevel     string `env:"LOG_LEVEL"`
	DatabaseDSN  string `env:"DATABASE_URI"`
	JWTSecretKey string `env:"JWT_SECRET"`
	JWTTimeExp   time.Duration

	JWTTimeExpInMinutes int `env:"JWT_TIME_EXP"`
}

func ParseConfig() *Config {
	var config Config
	flag.StringVar(&config.Address, "a", ":8080", "address and port to run server")
	flag.StringVar(&config.LogLevel, "l", "info", "log level")
	flag.StringVar(&config.DatabaseDSN, "d", "", "database connection string")
	flag.StringVar(&config.JWTSecretKey, "k", "", "jwt secret key")
	flag.IntVar(&config.JWTTimeExpInMinutes, "t", 10, "jwt time exp in minutes")
	flag.Parse()

	if err := env.Parse(&config); err != nil {
		log.Fatal(err)
	}

	config.JWTTimeExp = time.Minute * time.Duration(config.JWTTimeExpInMinutes)

	return &config
}
