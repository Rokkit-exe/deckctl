package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Rokkit-exe/deckctl/config"
	"github.com/Rokkit-exe/deckctl/internal/daemon"
	"github.com/Rokkit-exe/deckctl/internal/ipc"
)

var DaemonCommand = Command{
	Name:        "daemon",
	Description: "Run the program in daemon mode, connecting to the serial port and listening for input.",
	SubCommands: []SubCommand{
		{
			Name:        "start",
			Description: "Start the daemon",
		},
		{
			Name:        "stop",
			Description: "Stop the daemon",
		},
		{
			Name:        "restart",
			Description: "Restart the daemon",
		},
		{
			Name:        "run",
			Description: "Run the daemon in the foreground (for debugging)",
		},
	},
	Handler: handleDaemon,
}

func handleDaemon(args []string, cmd *Command) {
	if len(args) == 0 {
		fmt.Println("Daemon mode requires a subcommand:")
		cmd.PrintUsage()
		return
	}
	switch args[0] {
	case "start":
		fmt.Println("Starting daemon...")
		startDaemon()
	case "stop":
		fmt.Println("Stopping daemon...")
		sendDaemonCommand("stop")
	case "restart":
		fmt.Println("Restarting daemon...")
		sendDaemonCommand("stop")
		startDaemon()
	case "run":
		fmt.Println("Running daemon in foreground...")
		runDaemon()
	default:
		println("Unknown subcommand:", args[0])
		cmd.PrintUsage()
		os.Exit(1)
	}
}

func runDaemon() {
	cfgPath := filepath.Join(os.Getenv("HOME"), ".config/deckctl/config.yaml")
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Println("Config loaded, starting daemon...")

	err = daemon.Start(cfg)
	if err != nil {
		log.Fatalf("Failed to start daemon: %v", err)
	}
}

func startDaemon() {
	cmd := exec.Command("systemctl", "start", "deckctl")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Failed to start daemon via systemd:", err)
	}
}

func sendDaemonCommand(cmdStr string) {
	req := ipc.Request{
		Command: cmdStr,
	}

	resp, err := ipc.Send(req)
	if err != nil {
		fmt.Println("Failed to communicate with daemon:", err)
		return
	}

	fmt.Println(resp.Message)
}
