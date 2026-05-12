package model

type Config struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Version string            `json:"version"`
	Params  map[string]string `json:"params"`
}
type ConfigDTO struct {
	Name   string            `json:"name"`
	Params map[string]string `json:"params"`
}

type ConfigGroup struct {
	ID      string      `json:"id"`
	Name    string      `json:"name"`
	Version string      `json:"version"`
	Configs []ConfigDTO `json:"configs"`
	// Labels  map[string]string `json:"labels"`
}

type ConfigRepository interface {
	Add(config Config) error
	Get(name string, version string) (Config, error)
	Delete(name string, version string) error
	GetAll() ([]Config, error)
}

type ConfigGroupRepository interface {
	AddGroup(group ConfigGroup) error
	GetGroup(name string, version string) (ConfigGroup, error)
	GetAllGroups() ([]ConfigGroup, error)
	DeleteGroup(name string, version string) error
	UpdateGroup(group ConfigGroup) error
}
