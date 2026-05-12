package services

import (
	"central-config-service/model"
	"errors"

	"github.com/google/uuid"
)

type ConfigGroupService struct {
	groupRepo  model.ConfigGroupRepository
	configRepo model.ConfigRepository
}

func NewConfigGroupService(gr model.ConfigGroupRepository, cr model.ConfigRepository) *ConfigGroupService {
	return &ConfigGroupService{groupRepo: gr, configRepo: cr}
}

func (s *ConfigGroupService) AddGroup(group model.ConfigGroup) error {
	_, err := s.groupRepo.GetGroup(group.Name, group.Version)
	if err == nil {
		return errors.New("grupa sa tim imenom i verzijom vec postoji")
	}

	err = s.groupRepo.AddGroup(group)
	if err != nil {
		return err
	}

	for _, cfgDTO := range group.Configs {
		cfg := model.Config{
			ID:      uuid.New().String(),
			Name:    cfgDTO.Name,
			Version: group.Version,
			Params:  cfgDTO.Params,
		}
		_ = s.configRepo.Add(cfg)
	}

	return nil
}

func (s *ConfigGroupService) GetGroup(name string, version string) (model.ConfigGroup, error) {
	return s.groupRepo.GetGroup(name, version)
}

func (s *ConfigGroupService) GetAllGroups() ([]model.ConfigGroup, error) {
	return s.groupRepo.GetAllGroups()
}

func (s *ConfigGroupService) DeleteGroup(name string, version string) error {
	return s.groupRepo.DeleteGroup(name, version)
}

func (s *ConfigGroupService) AddConfigToGroup(groupName string, groupVersion string, config model.Config) error {
	group, err := s.groupRepo.GetGroup(groupName, groupVersion)
	if err != nil {
		return err
	}

	dto := model.ConfigDTO{
		Name:   config.Name,
		Params: config.Params,
	}
	group.Configs = append(group.Configs, dto)

	err = s.groupRepo.UpdateGroup(group)
	if err != nil {
		return err
	}

	_ = s.configRepo.Add(config)
	return nil
}

func (s *ConfigGroupService) AddExistingConfigToGroup(groupName string, groupVersion string, configName string, configVersion string) error {
	existingConfig, err := s.configRepo.Get(configName, configVersion)
	if err != nil {
		return err
	}

	group, err := s.groupRepo.GetGroup(groupName, groupVersion)
	if err != nil {
		return err
	}

	dto := model.ConfigDTO{
		Name:   existingConfig.Name,
		Params: existingConfig.Params,
	}
	group.Configs = append(group.Configs, dto)

	return s.groupRepo.UpdateGroup(group)
}

func (s *ConfigGroupService) DeleteConfigFromGroup(groupName string, groupVersion string, config model.Config) error {
	group, err := s.groupRepo.GetGroup(groupName, groupVersion)
	if err != nil {
		return err
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
	return s.groupRepo.UpdateGroup(group)
}
