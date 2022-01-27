package config

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	DB string `env:"DB" envDefault:"mongodb"`
	//DB             string `env:"DB" envDefault:"postgres"`
	User           string `env:"USER" envDefault:"egormelnikov"`
	Password       string `env:"PASSWORD" envDefault:"54236305"`
	Host           string `env:"HOST" envDefault:"postgres"`
	HostMongo      string `env:"HOST" envDefault:"mongo"`
	PortPostgres   int    `env:"PORT" envDefault:"5432"`
	PortMongo      int    `env:"PORT" envDefault:"27017"`
	DBNameMongo    string `env:"DBNAMEMONGO" envDefault:"egormelnikovdb"`
	DBNamePostgres string `env:"DBNAMEPOSTGRES" envDefault:"egormelnikovdb"`
	DBURL          string `env:"DBURL" envDefault:""`
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
