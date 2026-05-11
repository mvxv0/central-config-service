package services

import "central-config-service/model"

type ConfigGroupService struct {
	groupRepo  model.ConfigGroupRepository
	configRepo model.ConfigRepository
}

func NewConfigGroupService(gr model.ConfigGroupRepository, cr model.ConfigRepository) *ConfigGroupService {
	return &ConfigGroupService{groupRepo: gr, configRepo: cr}
}

func (s *ConfigGroupService) AddGroup(group model.ConfigGroup) error {
	return s.groupRepo.AddGroup(group)
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
	err := s.groupRepo.AddConfigToGroup(groupName, groupVersion, config)
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

	err = s.groupRepo.AddConfigToGroup(groupName, groupVersion, existingConfig)
	if err != nil {
		return err
	}

	return nil
}

func (s *ConfigGroupService) DeleteConfigFromGroup(groupName string, groupVersion string, config model.Config) error {
	return s.groupRepo.DeleteConfigFromGroup(groupName, groupVersion, config)
}
