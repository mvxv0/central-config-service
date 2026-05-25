package repositories

import (
	"central-config-service/model"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/hashicorp/consul/api"
)

// consul komunicira preko bajtova
type ConfigConsulRepository struct {
	cli *api.Client
}

func NewConfigConsulRepository() (model.ConfigRepository, error) {
	config := api.DefaultConfig()
	consulAddr := os.Getenv("CONSUL_ADDRESS")

	//ako je promenljiva prazna setuje se na port 8500 (ako se startuje lokalno sa go run)
	if consulAddr == "" {
		consulAddr = "localhost:8500"
	}

	config.Address = consulAddr

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &ConfigConsulRepository{cli: client}, nil
}

func (c *ConfigConsulRepository) Add(config model.Config) error {
	key := fmt.Sprintf("configs/%s/%s", config.Name, config.Version)

	data, err := json.Marshal(config) //json.Marshal pretvara go strukturu u json bajtove
	if err != nil {
		return err
	}

	kv := &api.KVPair{
		Key:   key,
		Value: data,
	}

	_, err = c.cli.KV().Put(kv, nil)
	return err
}

func (c *ConfigConsulRepository) Get(name string, version string) (model.Config, error) {
	key := fmt.Sprintf("configs/%s/%s", name, version)

	pair, _, err := c.cli.KV().Get(key, nil)
	if err != nil {
		return model.Config{}, err
	}

	if pair == nil {
		return model.Config{}, errors.New("konfiguracija nije pronadjena u Consul-u")
	}

	var config model.Config
	err = json.Unmarshal(pair.Value, &config) // json.Unmarshal vraca nazad json bajtove u go strukturu
	if err != nil {
		return model.Config{}, err
	}

	return config, nil
}

func (c *ConfigConsulRepository) Delete(name string, version string) error {
	key := fmt.Sprintf("configs/%s/%s", name, version)

	pair, _, _ := c.cli.KV().Get(key, nil)
	if pair == nil {
		return errors.New("konfiguracija ne postoji")
	}

	_, err := c.cli.KV().Delete(key, nil)
	return err
}

func (c *ConfigConsulRepository) GetAll() ([]model.Config, error) {
	pairs, _, err := c.cli.KV().List("configs/", nil)
	if err != nil {
		return nil, err
	}

	var allConfigs []model.Config

	for _, pair := range pairs {
		var config model.Config
		err := json.Unmarshal(pair.Value, &config)
		if err != nil {
			continue // preskacemo ostecene podatke
		}
		allConfigs = append(allConfigs, config)
	}

	return allConfigs, nil
}
