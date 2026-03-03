package protocol

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

type Report struct {
	Buttons uint8
	Slider1 uint8
	Slider2 uint8
	Slider3 uint8
}

func DecodePacket(buf []byte) (*Packet, error) {
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

func DecodeReport(packet Packet) (*Report, error) {
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
