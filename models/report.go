package models

import (
	"fmt"
)

type Report struct {
	Buttons uint8
	Slider1 uint8
	Slider2 uint8
	Slider3 uint8
}

func ParseReport(packet Packet) (*Report, error) {
	if packet.ID != Input {
		return nil, fmt.Errorf("unexpected report ID: %02x", packet.ID)
	}
	if packet.PayloadLen != 4 {
		return nil, fmt.Errorf("unexpected input report length: %d", packet.PayloadLen)
	}
	report := &Report{
		Buttons: packet.Payload[0],
		Slider1: packet.Payload[1],
		Slider2: packet.Payload[2],
		Slider3: packet.Payload[3],
	}
	return report, nil
}
