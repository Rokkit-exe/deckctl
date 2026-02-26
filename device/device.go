package device

import (
	"fmt"
	"github.com/karalabe/hid"
	"github.com/tarm/serial"
	"time"
)

func EnumerateDevices() {
	devices := hid.Enumerate(0, 0) // 0,0 = all devices
	for _, d := range devices {
		fmt.Printf("VID: %04x PID: %04x - %s\n", d.VendorID, d.ProductID, d.Product)
	}
}

func OpenDevice(vid, pid uint16) (*hid.Device, error) {
	devices := hid.Enumerate(vid, pid)
	if len(devices) == 0 {
		return nil, fmt.Errorf("device %04x:%04x not found", vid, pid)
	}
	dev, err := devices[0].Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open device: %w", err)
	}
	return dev, nil
}

func OpenSerial(port string, baud int) (*serial.Port, error) {
	c := &serial.Config{
		Name:        port,
		Baud:        baud,
		ReadTimeout: time.Millisecond * 500, // 500ms timeout
	}
	s, err := serial.OpenPort(c)
	if err != nil {
		return nil, fmt.Errorf("failed to open serial port: %w", err)
	}
	return s, nil
}

func ReadSerial(s *serial.Port, buf []byte) (int, error) {
	n, err := s.Read(buf)
	if err != nil {
		return n, fmt.Errorf("failed to read from serial port: %w", err)
	}
	return n, nil
}

func WriteSerial(s *serial.Port, data []byte) (int, error) {
	n, err := s.Write(data)
	if err != nil {
		return n, fmt.Errorf("failed to write to serial port: %w", err)
	}
	fmt.Printf("Wrote %d bytes to serial port\n", n)
	return n, nil
}

func CloseSerial(s *serial.Port) error {
	return s.Close()
}

func ReadLoop(s *serial.Port, reportChan chan<- []byte) {
	buf := make([]byte, 64)
	for {
		n, err := ReadSerial(s, buf)
		if err != nil {
			fmt.Printf("Serial read error: %v\n", err)
			continue
		}
		if n > 0 {
			reportChan <- buf[:n]
		}
	}
}

func WriteLoop(s *serial.Port, cmdChan <-chan []byte) {
	for cmd := range cmdChan {
		_, err := WriteSerial(s, cmd)
		if err != nil {
			fmt.Printf("Serial write error: %v\n", err)
		}
	}
}
