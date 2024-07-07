package main

import (
	"errors"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Database DatabaseConfig `yaml:"database"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

func LoadConfigFromFiles(fileNames ...string) (*Config, error) {
	if len(fileNames) == 0 {
		return nil, errors.New("no files provided")
	}

	var cfg Config
	var errs []error
	for _, fileName := range fileNames {
		err := loadConfigFromFile(fileName, &cfg)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == len(fileNames) {
		return nil, errors.New("failed to load config from all files")
	}

	return &cfg, errors.Join(errs...)
}

func loadConfigFromFile(fileName string, cfg *Config) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	return decoder.Decode(cfg)
}
