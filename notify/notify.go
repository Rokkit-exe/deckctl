package notify

import (
	"github.com/Rokkit-exe/deckctl/logger"
	"os/exec"
)

func Notify(title, message string) {
	cmd := exec.Command("notify-send", title, message)
	err := cmd.Run()
	if err != nil {
		logger.Error("Failed to send notification: " + err.Error())
	}
}
