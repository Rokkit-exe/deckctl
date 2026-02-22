package device

import (
	"fmt"
	"github.com/karalabe/hid"
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
