package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

func Load() (*Config, error) {
	cfg := new(Config)

	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	// Read base config
	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to read from base.yaml: %w", err)
	}

	// Validate loaded config values
	err = cfg.validate()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
