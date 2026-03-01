package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Rokkit-exe/deckctl/cdc"
	"github.com/Rokkit-exe/deckctl/cmd"
	"github.com/Rokkit-exe/deckctl/models"
	"go.bug.st/serial"
)

func HandleButtonPress(report *models.Report, cfg *models.Config) {
	for i, btn := range cfg.Buttons {
		mask := uint8(1 << i)
		isPressed := report.Buttons&mask != 0

		if isPressed {
			fmt.Printf("Button %d pressed → %s\n", btn.ID, btn.Action)
			go cmd.Exec(btn.Action)
		}
	}
}

func HandleSliderChange(report *models.Report, lastReport *models.Report, cfg *models.Config) {
	sliderValues := [3]uint8{report.Slider1, report.Slider2, report.Slider3}
	lastValues := [3]uint8{lastReport.Slider1, lastReport.Slider2, lastReport.Slider3}

	for i, sld := range cfg.Sliders {
		if sliderValues[i] != lastValues[i] {
			fmt.Printf("Slider %d changed to %d → %s\n", sld.ID, sliderValues[i], sld.Action)
			go cmd.Exec(fmt.Sprintf("%s %d", sld.Action, sliderValues[i]))
		}
	}
}

func main() {
	var port *serial.Port = nil
	cfg, err := models.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	port, err = cdc.Open("/dev/ttyACM1", 115200)
	if err != nil {
		log.Fatalf("Failed to open serial port: %v", err)
	}
	time.Sleep(2 * time.Second) // Wait for device to be ready
	defer cdc.Close(*port)

	if port == nil {
		log.Fatal("Serial port not available")
	}

	data := cfg.FormatData()
	_, err = cdc.Write(*port, data)
	if err != nil {
		log.Fatalf("Failed to write to serial port: %v", err)
	}
	fmt.Println("Listening for reports...")

	var lastReport models.Report

	for {
		packet, err := cdc.Read(*port)
		if err != nil {
			log.Printf("Failed to parse report: %v", err)
			continue
		}

		switch packet.ID {
		case models.ACK:
			fmt.Println("ESP32 ACK received")
			continue
		case models.Input:
			report, err := models.ParseReport(*packet)
			if err != nil {
				log.Printf("Failed to parse input report: %v", err)
				continue
			}
			fmt.Printf("Input report: %+v\n", report)
			HandleButtonPress(report, cfg)
			HandleSliderChange(report, &lastReport, cfg)
			lastReport = *report
		default:
			log.Printf("Unknown packet ID: %02x", packet.ID)
			continue
		}
	}
}
