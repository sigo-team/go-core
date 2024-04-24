package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env       string        `yaml:"env" env-required:"true"`
	Host      string        `yaml:"host" env-required:"true"`
	Port      int           `yaml:"port" env-required:"true"`
	JWTSecret string        `yaml:"JWTSecret" env-required:"true"`
	JWTMaxAge time.Duration `yaml:"JWTMaxAge" env-required:"true"`
}

func MustLoad() Config {
	path := "./configs/local.yml"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("config file does not exist in the path %v", path)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Fatalf("cannot read config %s", err)
	}

	return cfg
}
