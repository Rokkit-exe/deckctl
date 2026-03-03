package main

import (
	"os"

	"github.com/Rokkit-exe/deckctl/cmd"
)

func main() {
	var App = cmd.CLI{
		Commands: []cmd.Command{
			cmd.DaemonCommand,
			cmd.FlashCommand,
		},
	}
	App.Execute(os.Args)
}
