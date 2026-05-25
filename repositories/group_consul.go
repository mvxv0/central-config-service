package repositories

import (
	"central-config-service/model"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/hashicorp/consul/api"
)

type ConfigGroupConsulRepository struct {
	cli *api.Client
}

func NewConfigGroupConsulRepository() (model.ConfigGroupRepository, error) {
	config := api.DefaultConfig()
	consulAddr := os.Getenv("CONSUL_ADDRESS")
	if consulAddr == "" {
		consulAddr = "localhost:8500"
	}

	config.Address = consulAddr

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &ConfigGroupConsulRepository{cli: client}, nil
}

func (r *ConfigGroupConsulRepository) AddGroup(group model.ConfigGroup) error {
	key := fmt.Sprintf("groups/%s/%s", group.Name, group.Version)

	data, err := json.Marshal(group)
	if err != nil {
		return err
	}

	kv := &api.KVPair{
		Key:   key,
		Value: data,
	}

	_, err = r.cli.KV().Put(kv, nil)
	return err
}

func (r *ConfigGroupConsulRepository) GetGroup(name string, version string) (model.ConfigGroup, error) {
	key := fmt.Sprintf("groups/%s/%s", name, version)

	pair, _, err := r.cli.KV().Get(key, nil)
	if err != nil {
		return model.ConfigGroup{}, err
	}

	if pair == nil {
		return model.ConfigGroup{}, errors.New("grupa nije pronadjena u Consul-u")
	}

	var group model.ConfigGroup
	err = json.Unmarshal(pair.Value, &group)
	if err != nil {
		return model.ConfigGroup{}, err
	}

	return group, nil
}

func (r *ConfigGroupConsulRepository) GetAllGroups() ([]model.ConfigGroup, error) {
	pairs, _, err := r.cli.KV().List("groups/", nil)
	if err != nil {
		return nil, err
	}

	var allGroups []model.ConfigGroup

	for _, pair := range pairs {
		var group model.ConfigGroup
		err := json.Unmarshal(pair.Value, &group)
		if err != nil {
			continue
		}
		allGroups = append(allGroups, group)
	}

	return allGroups, nil
}

func (r *ConfigGroupConsulRepository) DeleteGroup(name string, version string) error {
	key := fmt.Sprintf("groups/%s/%s", name, version)

	pair, _, _ := r.cli.KV().Get(key, nil)
	if pair == nil {
		return errors.New("grupa nije pronadjena")
	}

	_, err := r.cli.KV().Delete(key, nil)
	return err
}

func (r *ConfigGroupConsulRepository) UpdateGroup(group model.ConfigGroup) error {
	// U Consulu je Put() komanda ista i za Create i za Update (prepisuje postojeće)
	return r.AddGroup(group)
}
