package repositories

import (
	"central-config-service/model"
	"errors"
	"fmt"
)

type ConfigGroupInMemRepository struct {
	groups map[string]model.ConfigGroup
}

func NewConfigGroupInMemRepository() model.ConfigGroupRepository {
	return &ConfigGroupInMemRepository{
		groups: make(map[string]model.ConfigGroup),
	}
}

func (r *ConfigGroupInMemRepository) AddGroup(group model.ConfigGroup) error {
	key := fmt.Sprintf("%s/%s", group.Name, group.Version)
	if _, ok := r.groups[key]; ok {
		return errors.New("grupa sa tim imenom i verzijom vec postoji")
	}
	r.groups[key] = group
	return nil
}

func (r *ConfigGroupInMemRepository) GetGroup(name string, version string) (model.ConfigGroup, error) {
	key := fmt.Sprintf("%s/%s", name, version)
	group, ok := r.groups[key]
	if !ok {
		return model.ConfigGroup{}, errors.New("grupa nije pronadjena")
	}
	return group, nil
}

func (r *ConfigGroupInMemRepository) DeleteGroup(name string, version string) error {
	key := fmt.Sprintf("%s/%s", name, version)
	if _, ok := r.groups[key]; !ok {
		return errors.New("grupa nije pronadjena")
	}
	delete(r.groups, key)
	return nil
}

func (r *ConfigGroupInMemRepository) AddConfigToGroup(name string, version string, config model.Config) error {
	key := fmt.Sprintf("%s/%s", name, version)
	group, ok := r.groups[key]
	if !ok {
		return errors.New("grupa nije pronadjena")
	}
	group.Configs = append(group.Configs, config)
	r.groups[key] = group
	return nil
}
