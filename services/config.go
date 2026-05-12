package services

import (
	"central-config-service/model"
	"errors"
)

type ConfigService struct {
	repo model.ConfigRepository
}

func NewConfigService(r model.ConfigRepository) *ConfigService {
	return &ConfigService{repo: r}
}

func (s *ConfigService) Add(c model.Config) error {
	_, err := s.repo.Get(c.Name, c.Version)
	if err == nil {
		return errors.New("konfiguracija vec postoji")
	}
	return s.repo.Add(c)
}

func (s *ConfigService) Get(name string, version string) (model.Config, error) {
	return s.repo.Get(name, version)
}

func (s *ConfigService) Delete(name string, version string) error {

	return s.repo.Delete(name, version)
}

func (s *ConfigService) GetAll() ([]model.Config, error) {
	return s.repo.GetAll()
}
