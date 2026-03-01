package parser

import (
	"github.com/Rokkit-exe/deckctl/models"
)

func ParseReport(data []byte) models.Report {
	// data[0] = report ID, data[1..4] = payload
	return models.Report{
		Buttons: data[1],
		Slider1: data[2],
		Slider2: data[3],
		Slider3: data[4],
	}
}
