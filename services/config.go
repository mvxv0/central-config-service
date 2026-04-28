package services

import "central-config-service/model"

type ConfigService struct {
	repo model.ConfigRepository
}

func NewConfigService(repo model.ConfigRepository) *ConfigService {
	return &ConfigService{repo: repo}
}

func (s *ConfigService) Add(c model.Config) error {
	return s.repo.Add(c)
}

func (s *ConfigService) Get(name string, version string) (model.Config, error) {
	return s.repo.Get(name, version)
}
