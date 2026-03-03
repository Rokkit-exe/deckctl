package protocol

import (
	"encoding/hex"
	"strings"

	"github.com/Rokkit-exe/deckctl/config"
)

type ButtonData struct {
	Label [16]byte
	Color [3]byte
}

type SliderData struct {
	Label [16]byte
	Color [3]byte
	Value byte
}

type ConfigData struct {
	ButtonsData [8]ButtonData
	SlidersData [3]SliderData
}

func parseHexColor(s string) [3]byte {
	var color [3]byte

	s = strings.TrimPrefix(s, "#")

	if len(s) == 6 {
		bytes, err := hex.DecodeString(s)
		if err == nil && len(bytes) == 3 {
			copy(color[:], bytes)
		}
	}
	return color
}

func encodeButtons(buttons []config.Button) [8]ButtonData {
	var buttonsReport [8]ButtonData
	for i, btn := range buttons {
		btn.Name = strings.ToValidUTF8(btn.Name, "")
		var label [16]byte
		copy(label[:], btn.Name)
		color := parseHexColor(btn.Color)
		buttonsReport[i] = ButtonData{
			Label: label,
			Color: color,
		}
	}
	return buttonsReport
}

func encodeSliders(sliders []config.Slider) [3]SliderData {
	var slidersReport [3]SliderData
	for i, sld := range sliders {
		var label [16]byte
		val := sld.Value
		val = min(max(val, 0), 100) // Ensure value is between 0 and 100
		copy(label[:], sld.Name)
		color := parseHexColor(sld.Color)
		slidersReport[i] = SliderData{
			Label: label,
			Color: color,
			Value: byte(val),
		}
	}
	return slidersReport
}

func getConfigReport(cfg config.Config) ConfigData {
	return ConfigData{
		ButtonsData: encodeButtons(cfg.Buttons[:]),
		SlidersData: encodeSliders(cfg.Sliders[:]),
	}
}

func Encode(cfg *config.Config) []byte {
	report := getConfigReport(*cfg)

	payload := make([]byte, 0)

	for _, btn := range report.ButtonsData {
		payload = append(payload, btn.Label[:]...)
		payload = append(payload, btn.Color[:]...)
	}

	for _, sld := range report.SlidersData {
		payload = append(payload, sld.Label[:]...)
		payload = append(payload, sld.Color[:]...)
		payload = append(payload, sld.Value)
	}

	totalLen := len(payload)

	data := make([]byte, 0, totalLen+4)

	data = append(data, 0x10)              // report ID
	data = append(data, 0x01)              // version
	data = append(data, byte(totalLen))    // length LSB
	data = append(data, byte(totalLen>>8)) // length MSB
	data = append(data, payload...)

	return data
}
