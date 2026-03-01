package models

import "fmt"

type PacketType byte

const (
	Input PacketType = 0x01
	ACK   PacketType = 0x10
)

type Packet struct {
	ID         PacketType
	Version    byte
	PayloadLen int
	Payload    []byte
}

func ParsePacket(buf []byte) (*Packet, error) {
	if len(buf) < 4 {
		return nil, fmt.Errorf("buffer too short")
	}
	p := &Packet{
		ID:         PacketType(buf[0]),
		Version:    buf[1],
		PayloadLen: int(buf[2]) | int(buf[3])<<8,
	}

	if len(buf) < 4+p.PayloadLen {
		return nil, fmt.Errorf("buffer shorter than payload length")
	}

	p.Payload = buf[4 : 4+p.PayloadLen]
	return p, nil
}
