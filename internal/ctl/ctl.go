package ctl

import (
	"fmt"
	"github.com/Rokkit-exe/deckctl/config"
	"github.com/Rokkit-exe/deckctl/internal/cmd"
	"github.com/Rokkit-exe/deckctl/internal/protocol"
	"github.com/Rokkit-exe/deckctl/internal/serial"
)

type Controller struct {
	Manager    *serial.Manager
	LastReport protocol.Report
}

func NewController(cfg *config.Config) *Controller {
	return &Controller{
		Manager:    serial.NewManager(cfg),
		LastReport: protocol.Report{Buttons: 0, Slider1: 0, Slider2: 0, Slider3: 0},
	}
}

func (c *Controller) Start() error {
	go c.Manager.Run()
	fmt.Println("Serial manager started.")
	go c.Handle()
	fmt.Println("Controller handler started.")
	return nil
}

func (c *Controller) Stop() {
	c.Manager.Stop()
	fmt.Println("Serial manager stopped.")
}

func (c *Controller) Restart() {
	c.Manager.Stop()
	c.Manager.Run()
	fmt.Println("Serial manager restarted.")
}

func (c *Controller) Handle() {
	for data := range c.Manager.RxChan {
		pck, err := protocol.DecodePacket(data)
		if err != nil {
			continue
		}

		switch pck.ID {
		case protocol.ACK:
			fmt.Println("ESP32 ACK received")
		case protocol.Input:
			report, err := protocol.DecodeReport(*pck)
			if err != nil {
				fmt.Printf("Failed to decode report: %v\n", err)
				continue
			}
			HandleButtonPress(report, c.Manager.Cfg)
			HandleSliderChange(report, &c.LastReport, c.Manager.Cfg)
			c.LastReport = *report
		default:
			fmt.Printf("Unknown packet ID: %02x\n", pck.ID)
		}
	}
}

func (c *Controller) SendAck() {
	ack := []byte{0x10, 0x01, 0x00, 0x00}
	c.Manager.TxChan <- ack
}

func (c *Controller) Flash(file string) {
	cfg, err := config.LoadConfig(file)
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		return
	}
	c.Manager.Cfg = cfg
	data := protocol.Encode(c.Manager.Cfg)
	c.Manager.TxChan <- data
}

func HandleButtonPress(report *protocol.Report, cfg *config.Config) {
	for i, btn := range cfg.Buttons {
		mask := uint8(1 << i)
		isPressed := report.Buttons&mask != 0

		if isPressed {
			fmt.Printf("Button %d pressed → %s\n", btn.ID, btn.Action)
			go cmd.Exec(btn.Action)
		}
	}
}

func HandleSliderChange(report *protocol.Report, lastReport *protocol.Report, cfg *config.Config) {
	sliderValues := [3]uint8{report.Slider1, report.Slider2, report.Slider3}
	lastValues := [3]uint8{lastReport.Slider1, lastReport.Slider2, lastReport.Slider3}

	for i, sld := range cfg.Sliders {
		if sliderValues[i] != lastValues[i] {
			fmt.Printf("Slider %d changed to %d → %s\n", sld.ID, sliderValues[i], sld.Action)
			go cmd.Exec(fmt.Sprintf("%s %d", sld.Action, sliderValues[i]))
		}
	}
}
