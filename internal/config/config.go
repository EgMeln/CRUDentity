// Package config to env
package config

import (
	"github.com/caarlos0/env/v6"
)

// Config struct to config env
type Config struct {
	// DB string `env:"DB" envDefault:"mongodb"`
	DB             string `env:"DB" envDefault:"postgres"`
	User           string `env:"USER" envDefault:"egormelnikov"`
	Password       string `env:"PASSWORD" envDefault:"54236305"`
	Host           string `env:"HOST" envDefault:"postgresql"`
	HostMongo      string `env:"HOST" envDefault:"mongo"`
	PortPostgres   int    `env:"PORT" envDefault:"5432"`
	PortMongo      int    `env:"PORT" envDefault:"27017"`
	DBNameMongo    string `env:"DBNAMEMONGO" envDefault:"egormelnikovdb"`
	DBNamePostgres string `env:"DBNAMEPOSTGRES" envDefault:"egormelnikov"`
	DBURL          string `env:"DBURL" envDefault:""`

	AccessToken          string `env:"ACCESSTOKEN" envDefault:"access_key"`
	RefreshToken         string `env:"REFRESHTOKEN" envDefault:"refresh_key"`
	AccessTokenLifeTime  int    `env:"ACCESSLIFETIME" envDefault:"1800"`
	RefreshTokenLifeTime int    `env:"REFRESHLIFETOKEN" envDefault:"518400"`
	HashSalt             string `env:"HASHSALT" envDefault:"HASHSALT"`

	Server string `env:"SERVER" envDefault:"grpc"`
}

// New contract config
func New() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
