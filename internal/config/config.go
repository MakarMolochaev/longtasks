package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string        `yaml:"env"`
	Redis       RedisConfig   `yaml:"redis"`
	Http        HttpConfig    `yaml:"http"`
	TaskTimeout time.Duration `yaml:"task_timeout"`
}

type RedisConfig struct {
	Url string `yaml:"url"`
}

type HttpConfig struct {
	Url string `yaml:"url"`
}

func MustLoad() *Config {
	var path string

	flag.StringVar(&path, "config", "", "")
	flag.Parse()

	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}
