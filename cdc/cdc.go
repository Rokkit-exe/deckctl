package cdc

import (
	"fmt"
	"github.com/Rokkit-exe/deckctl/models"
	"go.bug.st/serial"
	"time"
)

var ACKPacket = []byte{0x10, 0x01, 0x00, 0x00}
var maxPacketSize = 64
var packetHeaderLen = 4
var rxBuf []byte

func Open(portName string, baud int) (*serial.Port, error) {
	mode := &serial.Mode{BaudRate: baud}

	port, err := serial.Open(portName, mode)
	if err != nil {
		return nil, err
	}

	err = port.SetDTR(true) // REQUIRED!
	if err != nil {
		port.Close()
		return nil, err
	}
	err = port.SetRTS(true)
	if err != nil {
		port.Close()
		return nil, err
	}

	time.Sleep(100 * time.Millisecond) // Wait for the device to reset
	_, err = Write(port, ACKPacket)    // Send ACK to indicate we're ready
	if err != nil {
		port.Close()
		return nil, err
	}

	return &port, nil
}

func Close(port serial.Port) error {
	if port == nil {
		return nil
	}
	return port.Close()
}

func Write(port serial.Port, data []byte) (int, error) {
	n, err := port.Write(data)
	if err != nil {
		return n, err
	}
	fmt.Printf("Wrote %d bytes: %x\n", n, data[:n])
	return n, nil
}

func HandlePacket(p *models.Packet) (*models.Report, error) {
	switch p.ID {
	case models.ACK:
		fmt.Println("ESP32 ACK received")
		return nil, nil
	case models.Input:
		report, err := models.ParseReport(*p)
		if err != nil {
			return nil, err
		}
		fmt.Printf("Input report: %+v\n", report)
		return report, nil
	default:
		return nil, fmt.Errorf("unknown report ID: %02x", p.ID)
	}
}

func Read(port serial.Port) (*models.Packet, error) {
	tmp := make([]byte, maxPacketSize)
	n, err := port.Read(tmp)
	if err != nil {
		return nil, err
	}
	rxBuf = append(rxBuf, tmp[:n]...)

	for len(rxBuf) >= packetHeaderLen {
		payloadLen := int(rxBuf[2]) | int(rxBuf[3])<<8
		fullLen := packetHeaderLen + payloadLen
		if len(rxBuf) < fullLen {
			break // wait for more bytes
		}

		packet, err := models.ParsePacket(rxBuf[:fullLen])
		if err != nil {
			rxBuf = rxBuf[fullLen:]
			continue
		}

		rxBuf = rxBuf[fullLen:]
		return packet, nil
	}
	return nil, nil
}
