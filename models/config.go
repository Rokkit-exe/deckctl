package models

import (
	"encoding/hex"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
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
	VID     uint16    `yaml:"VID"`
	PID     uint16    `yaml:"PID"`
	Buttons [8]Button `yaml:"buttons"`
	Sliders [3]Slider `yaml:"sliders"`
}

type ButtonReport struct {
	Label [16]byte
	Color [3]byte
}

type SliderReport struct {
	Label [16]byte
	Color [3]byte
	Value byte
}

type ConfigReport struct {
	ButtonReports [8]ButtonReport
	SliderReports [3]SliderReport
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

func parseHexColor(s string) [3]byte {
	var color [3]byte

	s = strings.TrimPrefix(s, "#")

	if len(s) == 6 {
		bytes, err := hex.DecodeString(s)
		if err == nil && len(bytes) == 3 {
			copy(color[:], bytes)
		}
	}
	return color
}

func (c *Config) getButtonReports() [8]ButtonReport {
	var buttonsReport [8]ButtonReport
	for i, btn := range c.Buttons {
		btn.Name = strings.ToValidUTF8(btn.Name, "")
		var label [16]byte
		copy(label[:], btn.Name)
		color := parseHexColor(btn.Color)
		buttonsReport[i] = ButtonReport{
			Label: label,
			Color: color,
		}
	}
	return buttonsReport
}

func (c *Config) getSliderReports() [3]SliderReport {
	var slidersReport [3]SliderReport
	for i, sld := range c.Sliders {
		var label [16]byte
		val := sld.Value
		val = min(max(val, 0), 100) // Ensure value is between 0 and 100
		copy(label[:], sld.Name)
		color := parseHexColor(sld.Color)
		slidersReport[i] = SliderReport{
			Label: label,
			Color: color,
			Value: byte(val),
		}
	}
	return slidersReport
}

func (c *Config) getConfigReport() ConfigReport {
	return ConfigReport{
		ButtonReports: c.getButtonReports(),
		SliderReports: c.getSliderReports(),
	}
}

func (c *Config) FormatData() []byte {
	report := c.getConfigReport()

	payload := make([]byte, 0)

	for _, btn := range report.ButtonReports {
		payload = append(payload, btn.Label[:]...)
		payload = append(payload, btn.Color[:]...)
	}

	for _, sld := range report.SliderReports {
		payload = append(payload, sld.Label[:]...)
		payload = append(payload, sld.Color[:]...)
		payload = append(payload, sld.Value)
	}

	totalLen := len(payload)

	data := make([]byte, 0, totalLen+4)

	data = append(data, 0x10)              // report ID
	data = append(data, 0x01)              // version
	data = append(data, byte(totalLen))    // length LSB
	data = append(data, byte(totalLen>>8)) // length MSB
	data = append(data, payload...)

	return data
}
