package models

type Button struct {
	ID     int    `yaml:"id"`
	Name   string `yaml:"name"`
	Type   string `yaml:"type"`
	State  bool   `yaml:"state"`
	Action string `yaml:"action"`
}

type Slider struct {
	ID     int    `yaml:"id"`
	Name   string `yaml:"name"`
	Type   string `yaml:"type"`
	Value  int    `yaml:"value"`
	Action string `yaml:"action"`
}

type Config struct {
	VID     uint16   `yaml:"VID"`
	PID     uint16   `yaml:"PID"`
	Buttons []Button `yaml:"buttons"`
	Sliders []Slider `yaml:"sliders"`
}
