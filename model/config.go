package model

type Config struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Version string            `json:"version"`
	Params  map[string]string `json:"params"`
}

type ConfigGroup struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Version string   `json:"version"`
	Configs []Config `json:"configs"`
}

type ConfigRepository interface {
	Add(config Config) error
	Get(name string, version string) (Config, error)
}

type ConfigGroupRepository interface {
	AddGroup(group ConfigGroup) error
	GetGroup(name string, version string) (ConfigGroup, error)
	DeleteGroup(name string, version string) error
	AddConfigToGroup(name string, version string, config Config) error
}
