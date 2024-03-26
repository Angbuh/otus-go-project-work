package main

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	// создаем конфиг
	LogLevel     string `env:"LOG_LEVEL" env-default:"debug"`
	BindIP       string `env:"BIND_IP" env-default:"0.0.0.0"`
	Port         string `env:"PORT" env-default:"8000"`
	DatabasePath string `env:"DATABASE_PATH" env-required:"true"`
}

func GetConfig() (Config, error) {
	//настраиваем конфиг
	var config Config

	// Считывает переменные среды в структуру Конфиг
	if err := cleanenv.ReadEnv(&config); err != nil {
		return config, err
	}

	return config, nil
}
