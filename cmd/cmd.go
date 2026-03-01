package cmd

import (
	"github.com/Rokkit-exe/deckctl/logger"
	"os/exec"
	"strings"
)

func Exec(action string) {
	parts := strings.Fields(action)
	if len(parts) == 0 {
		return
	}
	c := exec.Command(parts[0], parts[1:]...)
	if err := c.Start(); err != nil {
		logger.Error("Failed to start command: " + err.Error())
	}
}
