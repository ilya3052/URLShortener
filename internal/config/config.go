package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env          string `yaml:"env" env-default:"local"`
	Conn_str  string `yaml:"conn_str" env-required:"true"`
	HTTP_server  `yaml:"http_server"`
}

type HTTP_server struct {
	Address      string        `yaml:"address" env-default:"0.0.0.0:8080"`
	Timeout      time.Duration `yaml:"timeout" env-default:"5s"`
	Idle_timeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("Config file does not exist: ", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("Cannot read config: ", err)
	}

	return &cfg
}
