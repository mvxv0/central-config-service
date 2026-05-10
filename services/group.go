package services

import "central-config-service/model"

type ConfigGroupService struct {
	repo model.ConfigGroupRepository
}

func NewConfigGroupService(repo model.ConfigGroupRepository) *ConfigGroupService {
	return &ConfigGroupService{repo: repo}
}

func (s *ConfigGroupService) AddGroup(group model.ConfigGroup) error {
	return s.repo.AddGroup(group)
}

func (s *ConfigGroupService) GetGroup(name string, version string) (model.ConfigGroup, error) {
	return s.repo.GetGroup(name, version)
}

func (s *ConfigGroupService) DeleteGroup(name string, version string) error {
	return s.repo.DeleteGroup(name, version)
}

func (s *ConfigGroupService) AddConfigToGroup(name string, version string, config model.Config) error {
	return s.repo.AddConfigToGroup(name, version, config)
}
