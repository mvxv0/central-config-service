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

func (s *ConfigGroupService) AddConfigToGroup(groupName string, groupVersion string, configDTO model.ConfigDTO) error {
	group, err := s.groupRepo.GetGroup(groupName, groupVersion)
	if err != nil {
		return err
	}

	group.Configs = append(group.Configs, configDTO)

	err = s.groupRepo.UpdateGroup(group)
	if err != nil {
		return err
	}

	globalConfig := model.Config{
		ID:      uuid.New().String(),
		Name:    configDTO.Name,
		Version: groupVersion,
		Params:  configDTO.Params,
	}

	_ = s.configRepo.Add(globalConfig)
	return nil
}

func (s *ConfigGroupService) AddExistingConfigToGroup(groupName string, groupVersion string, LabelsConfigDto model.LabelsConfigDto) error {
	globalConfig, err := s.configRepo.Get(LabelsConfigDto.Name, LabelsConfigDto.Version)
	if err != nil {
		return errors.New("globalna konfiguracija ne postoji")
	}

	group, err := s.groupRepo.GetGroup(groupName, groupVersion)
	if err != nil {
		return err
	}

	newGroupConfig := model.ConfigDTO{
		Name:   globalConfig.Name,
		Params: globalConfig.Params,
		Labels: LabelsConfigDto.Labels,
	}

	group.Configs = append(group.Configs, newGroupConfig)
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

func (s *ConfigGroupService) GetConfigsByLabels(groupName string, groupVersion string, searchLabels map[string]string) ([]model.ConfigDTO, error) {
	// hook group from repo
	group, err := s.groupRepo.GetGroup(groupName, groupVersion)
	if err != nil {
		return nil, err
	}

	var filteredConfigs []model.ConfigDTO

	// pass through all configs from group
	for _, cfg := range group.Configs {
		match := true

		// checking if labels are matched
		for searchKey, searchValue := range searchLabels {
			// key ?
			cfgValue, exists := cfg.Labels[searchKey]

			// value ?
			if !exists || cfgValue != searchValue {
				match = false
				break
			}
		}

		if match {
			filteredConfigs = append(filteredConfigs, cfg)
		}
	}

	return filteredConfigs, nil
}

func (s *ConfigGroupService) DeleteConfigsByLabels(groupName string, groupVersion string, searchLabels map[string]string) error {
	group, err := s.groupRepo.GetGroup(groupName, groupVersion)
	if err != nil {
		return err
	}

	var keptConfigs []model.ConfigDTO
	deletedCount := 0

	for _, cfg := range group.Configs {
		match := true

		for searchKey, searchValue := range searchLabels {
			cfgValue, exists := cfg.Labels[searchKey]
			if !exists || cfgValue != searchValue {
				match = false
				break
			}
		}

		if match {
			deletedCount++
		} else {
			keptConfigs = append(keptConfigs, cfg)
		}
	}

	if deletedCount == 0 {
		return errors.New("nijedna konfiguracija u grupi ne ispunjava uslove za brisanje")
	}

	group.Configs = keptConfigs
	return s.groupRepo.UpdateGroup(group)
}
