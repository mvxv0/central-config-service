package model

// osnovni model
type Config struct {
	ID      string            `json:"id"` // uuid
	Name    string            `json:"name"`
	Version string            `json:"version"`
	Params  map[string]string `json:"params"` // kljuc - vrednost
}

type ConfigGroup struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Version string            `json:"version"`
	Configs []Config          `json:"configs"`
	Labels  map[string]string `json:"labels"`
}

// def interfejsa sta treba da radi
type ConfigRepository interface {
	Add(config Config) error
	Get(name string, version string) (Config, error)
}
