package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	StoragePath string           `yaml:"storage_path"`
	AliasLength uint64           `yaml:"alias_length"`
	LogLevel    string           `yaml:"log_level"`
	LogFormat   string           `yaml:"log_format"`
	HttpServer  HttpServerConfig `yaml:"http_server"`
}

type HttpServerConfig struct {
	Address  string `yaml:"address"`
	Username string `yaml:"username" env-required:"true"`
	Password string `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

func MustLoad(path string) *Config {
	c := &Config{}
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		slog.Error(fmt.Sprintf("Couldn't load config at path %s", path))
		os.Exit(1)
	}
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		slog.Error(fmt.Sprintf("Couldn't read config at path %s", path))
		os.Exit(1)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		slog.Error(fmt.Sprintf("Couldn't unmarshal config at path %s", path))
		os.Exit(1)
	}
	return c
}
