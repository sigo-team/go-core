package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env string `yaml:"env" env-required:"true"`
	//PlayersAmount uint   `yaml:"playersAmount" env-required:"true"`
	Host      string        `env:"HOST" env-required:"true"`
	Port      int           `env:"PORT" env-required:"true"`
	JWTSecret string        `env:"JWTSECRET" env-required:"true"`
	JWTMaxAge time.Duration `env:"JWT_MAX_AGE" env-required:"true"`
}

func MustLoad(path string) Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("config file does not exist in the path %v", path)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Fatalf("cannot read config %s", err)
	}

	return cfg
}
