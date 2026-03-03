package ipc

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

// Request from CLI
type Request struct {
	Command string `json:"command"`
	File    string `json:"file,omitempty"`
}

// Response to CLI
type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

const SocketPath = "/tmp/deckctl.sock"

// HandlerFunc processes a request and returns a response
type HandlerFunc func(Request) Response

// StartServer listens on the Unix socket and handles incoming connections
func StartServer(handler HandlerFunc) error {
	// Remove old socket if it exists
	if _, err := os.Stat(SocketPath); err == nil {
		os.Remove(SocketPath)
	}

	listener, err := net.Listen("unix", SocketPath)
	if err != nil {
		return err
	}
	defer listener.Close()

	os.Chmod(SocketPath, 0666) // Allow all users to connect

	fmt.Println("Daemon IPC server listening on", SocketPath)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}

		// Handle each connection in a new goroutine
		go func(c net.Conn) {
			defer c.Close()

			var req Request
			err := json.NewDecoder(c).Decode(&req)
			if err != nil {
				fmt.Println("Decode error:", err)
				return
			}

			resp := handler(req)

			err = json.NewEncoder(c).Encode(resp)
			if err != nil {
				fmt.Println("Encode error:", err)
				return
			}

		}(conn)
	}
}

// Send connects to the daemon socket, sends a request, and waits for a response
func Send(req Request) (*Response, error) {
	conn, err := net.Dial("unix", SocketPath)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Send the request
	if err := json.NewEncoder(conn).Encode(req); err != nil {
		return nil, err
	}

	// Wait for the response
	var resp Response
	if err := json.NewDecoder(conn).Decode(&resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

