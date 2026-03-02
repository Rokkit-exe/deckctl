package main

import (
	"fmt"
	"github.com/Rokkit-exe/deckctl/config"
	"github.com/Rokkit-exe/deckctl/ctl"
	"github.com/Rokkit-exe/deckctl/serial"
	"log"
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	manager := serial.NewManager(cfg)
	ctl := ctl.NewController(manager)

	ctl.SendConfig()
	fmt.Println("Connecting to serial port...")
	go manager.Run()
	fmt.Println("Connected. Listening for input...")
	ctl.Handle()
}
