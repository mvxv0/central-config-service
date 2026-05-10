package repositories

import (
	"central-config-service/model"
	"errors"
	"fmt"
)

type ConfigInMemRepository struct {
	configs map[string]model.Config
}

func NewConfigInMemRepository() model.ConfigRepository {
	return &ConfigInMemRepository{
		configs: make(map[string]model.Config),
	}
}

func (c *ConfigInMemRepository) Add(config model.Config) error {
	key := fmt.Sprintf("%s/%s", config.Name, config.Version)
	c.configs[key] = config
	return nil
}

func (c *ConfigInMemRepository) Get(name string, version string) (model.Config, error) {
	key := fmt.Sprintf("%s/%s", name, version)
	config, ok := c.configs[key]
	if !ok {
		return model.Config{}, errors.New("konfiguracija nije pronađena")
	}
	return config, nil
}
