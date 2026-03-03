package serial

import (
	"fmt"

	"github.com/Rokkit-exe/deckctl/internal/protocol"
	"go.bug.st/serial"
)

var maxPacketSize = 64
var packetHeaderLen = 4
var maxAllowedPayload = 212
var rxBuf []byte

type Port interface {
	Read(serial.Port) (*protocol.Packet, error)
	Write(serial.Port, []byte) (int, error)
	Close(serial.Port) error
}

func Close(p serial.Port) error {
	if p == nil {
		return nil
	}
	fmt.Println("Closing port")
	return p.Close()
}

func Write(p serial.Port, data []byte) (int, error) {
	n, err := p.Write(data)
	if err != nil {
		return n, err
	}
	return n, nil
}

func Read(p serial.Port) ([]byte, error) {
	tmp := make([]byte, maxPacketSize)

	n, err := p.Read(tmp)
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

		if payloadLen > maxAllowedPayload {
			rxBuf = rxBuf[1:] // shift by one to resync
			continue
		}

		// Copy frame BEFORE consuming buffer
		frame := make([]byte, fullLen)
		copy(frame, rxBuf[:fullLen])

		// Consume from buffer
		rxBuf = rxBuf[fullLen:]

		return frame, nil
	}

	return nil, nil
}
