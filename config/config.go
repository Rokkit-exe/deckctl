package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Button struct {
	ID     int    `yaml:"id"`
	Name   string `yaml:"name"`
	Type   string `yaml:"type"`
	State  bool   `yaml:"state"`
	Color  string `yaml:"color"`
	Action string `yaml:"action"`
}

type Slider struct {
	ID     int    `yaml:"id"`
	Name   string `yaml:"name"`
	Type   string `yaml:"type"`
	Color  string `yaml:"color"`
	Value  int    `yaml:"value"`
	Action string `yaml:"action"`
}

type Config struct {
	VID      uint16    `yaml:"VID"`
	PID      uint16    `yaml:"PID"`
	BaudRate int       `yaml:"BaudRate"`
	Port     string    `yaml:"Port"`
	Buttons  [8]Button `yaml:"buttons"`
	Sliders  [3]Slider `yaml:"sliders"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (cfg *Config) Save(path string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
