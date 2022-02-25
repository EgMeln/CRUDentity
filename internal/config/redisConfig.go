package config

import (
	"github.com/caarlos0/env/v6"
)

// Redis struct to redis config env
type Redis struct {
	Addr     string `env:"ADDR_REDIS" envDefault:"redis:6379"`
	Password string `env:"PASSWORD_REDIS" envDefault:""`
	DB       int    `env:"DB_REDIS" envDefault:"0"`
}

// NewRedis contract redis config
func NewRedis() (*Redis, error) {
	cfg := &Redis{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
