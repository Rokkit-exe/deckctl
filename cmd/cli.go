package cmd

import (
	"os"
)

type Command struct {
	Name        string
	Description string
	SubCommands []SubCommand
	Params      []Params
	Handler     HandlerFunc
}

type SubCommand struct {
	Name        string
	Description string
	Params      []Params
}

type Params struct {
	Name        string
	Description string
	Short       string
}

type CLI struct {
	ExecutableName string
	Commands       []Command
}

type HandlerFunc func(args []string, cmd *Command)

func (c *CLI) Execute(args []string) {
	c.ExecutableName = os.Args[0]
	if len(args) < 2 {
		c.PrintUsage()
		return
	}
	for _, cmd := range c.Commands {
		if args[1] == cmd.Name {
			if cmd.Handler == nil {
				println("No handlerFunc defined for command:", cmd.Name)
				cmd.PrintUsage()
				return
			}
			cmd.Handler(args[2:], &cmd)
			return
		}
	}
	c.PrintUsage()
}

func (c *CLI) PrintUsage() {
	println("Usage:")
	for _, cmd := range c.Commands {
		println("  " + cmd.Name + ": " + cmd.Description)
	}
	os.Exit(1)
}

func (cmd *Command) PrintUsage() {
	println("Usage: " + cmd.Name)
	println(cmd.Description)
	if len(cmd.SubCommands) > 0 {
		println("Subcommands:")
		for _, sub := range cmd.SubCommands {
			println("  " + sub.Name + ": " + sub.Description)
		}
	}
	if len(cmd.Params) > 0 {
		println("Parameters:")
		for _, param := range cmd.Params {
			println("  " + param.Short + ", " + param.Name + ": " + param.Description)
		}
	}
	os.Exit(1)
}
