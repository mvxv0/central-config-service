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

func (r *ConfigGroupInMemRepository) GetAllGroups() ([]model.ConfigGroup, error) {
	var allGroups []model.ConfigGroup
	for _, group := range r.groups {
		allGroups = append(allGroups, group)
	}
	return allGroups, nil
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

	dto := model.ConfigDTO{
		Name:   config.Name,
		Params: config.Params,
	}

	group.Configs = append(group.Configs, dto)

	r.groups[key] = group
	return nil
}

func (r *ConfigGroupInMemRepository) DeleteConfigFromGroup(name string, version string, config model.Config) error {
	key := fmt.Sprintf("%s/%s", name, version)

	group, ok := r.groups[key]
	if !ok {
		return errors.New("grupa nije pronadjena")
	}

	found := false
	var newConfigs []model.ConfigDTO
	for _, c := range group.Configs {
		if c.Name == config.Name {
			found = true
			continue
		}
		newConfigs = append(newConfigs, c)
	}

	if !found {
		return errors.New("config nije pronadjen u grupi")
	}

	group.Configs = newConfigs
	r.groups[key] = group
	return nil
}
