package daemon

import (
	"fmt"
	"os"

	"github.com/Rokkit-exe/deckctl/config"
	"github.com/Rokkit-exe/deckctl/internal/ctl"
	"github.com/Rokkit-exe/deckctl/internal/ipc"
)

func Start(cfg *config.Config) error {
	ctrl := ctl.NewController(cfg)

	err := ctrl.Start()
	if err != nil {
		return err
	}

	fmt.Println("Daemon started")

	// Start IPC server and process connections
	return ipc.StartServer(func(req ipc.Request) ipc.Response {
		switch req.Command {
		case "flash":
			ctrl.Flash(req.File)
			return ipc.Response{Status: "ok", Message: "Flashed"}
		case "stop":
			ctrl.Stop()
			go func() {
				os.Remove(ipc.SocketPath)
				os.Exit(0)
			}()
			return ipc.Response{Status: "ok", Message: "Shutting down"}
		default:
			return ipc.Response{Status: "error", Message: "unknown command"}
		}
	})
}
