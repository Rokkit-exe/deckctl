package main

import (
	"fmt"
	"github.com/Rokkit-exe/deckctl/device"
	"github.com/Rokkit-exe/deckctl/models"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"os/exec"
	"strings"
)

func execCommand(action string) {
	parts := strings.Fields(action)
	if len(parts) == 0 {
		return
	}
	cmd := exec.Command(parts[0], parts[1:]...)
	if err := cmd.Start(); err != nil {
		log.Printf("Command error: %v", err)
	}
}

func parseReport(data []byte) models.Report {
	// data[0] = report ID, data[1..4] = payload
	return models.Report{
		Buttons: data[1],
		Slider1: data[2],
		Slider2: data[3],
		Slider3: data[4],
	}
}

func LoadConfig(path string) (*models.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg models.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func main() {
	cfg, err := LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	device.EnumerateDevices()

	dev, err := device.OpenDevice(cfg.VID, cfg.PID)
	if err != nil {
		log.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	fmt.Println("Listening for reports...")

	var lastReport models.Report
	buf := make([]byte, 65) // 64 bytes + report ID

	for {
		n, err := dev.Read(buf)
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}
		if n < 5 {
			continue
		}

		report := parseReport(buf)

		// Handle button presses
		for i, btn := range cfg.Buttons {
			mask := uint8(1 << i)
			wasPressed := lastReport.Buttons&mask != 0
			isPressed := report.Buttons&mask != 0

			if isPressed && !wasPressed {
				fmt.Printf("Button %d pressed → %s\n", btn.ID, btn.Action)
				go execCommand(btn.Action)
			}
		}

		// Handle slider changes
		sliderValues := [3]uint8{report.Slider1, report.Slider2, report.Slider3}
		lastValues := [3]uint8{lastReport.Slider1, lastReport.Slider2, lastReport.Slider3}

		for i, sld := range cfg.Sliders {
			if sliderValues[i] != lastValues[i] {
				fmt.Printf("Slider %d changed to %d → %s\n", sld.ID, sliderValues[i], sld.Action)
				go execCommand(fmt.Sprintf("%s %d", sld.Action, sliderValues[i]))
			}
		}

		lastReport = report
	}
}
