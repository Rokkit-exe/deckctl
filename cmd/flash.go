package cmd

import (
	"fmt"

	"github.com/Rokkit-exe/deckctl/internal/ipc"
)

var FlashCommand = Command{
	Name:        "flash",
	Description: "Flash the firmware to the device. This will erase all data on the device and reset it to factory settings.",
	Params: []Params{
		{
			Name:        "file",
			Description: "Path to the config to flash",
			Short:       "f",
		},
		{
			Name:        "help",
			Description: "flash a config file to the device. This will erase all data on the device and reset it to factory settings. Usage: flash -f <config.yaml>",
			Short:       "h",
		},
	},
	Handler: handleFlash,
}

func handleFlash(args []string, cmd *Command) {
	if len(args) == 0 {
		fmt.Println("Flash command requires parameters:")
		cmd.PrintUsage()
		return
	}
	switch args[0] {
	case "-f", "--file":
		if len(args) < 2 {
			println("Error: Missing config file path")
			cmd.PrintUsage()
			return
		}
		configFile := args[1]
		println("Flashing config file:", configFile)
		sendFlashCommand("flash", configFile)
	case "-h", "--help":
		println(cmd.Description)
	default:
		println("Unknown parameter:", args[0])
		cmd.PrintUsage()
	}
}

func sendFlashCommand(cmdStr string, file string) {
	req := ipc.Request{
		Command: cmdStr,
		File:    file,
	}

	resp, err := ipc.Send(req)
	if err != nil {
		fmt.Println("Failed to communicate with daemon:", err)
		return
	}

	fmt.Println(resp.Message)
}
