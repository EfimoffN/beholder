package config

import (
	"github.com/EfimoffN/beholder/types"
	"github.com/kelseyhightower/envconfig"
)

func CreateConfig() (*types.ConfigApp, error) {
	cfg := types.ConfigApp{}

	err := envconfig.Process("BEHOLDER", &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
