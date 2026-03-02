package serial

import (
	"context"
	"fmt"
	"github.com/Rokkit-exe/deckctl/config"
	"go.bug.st/serial"
	"io"
	"time"
)

var ACKPacket = []byte{0x10, 0x01, 0x00, 0x00}

type Manager struct {
	Cfg  *config.Config
	Port io.ReadWriteCloser

	RxChan chan []byte
	TxChan chan []byte

	ctx    context.Context
	cancel context.CancelFunc
}

func NewManager(cfg *config.Config) *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	return &Manager{
		Cfg:    cfg,
		RxChan: make(chan []byte, 100),
		TxChan: make(chan []byte, 100),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (m *Manager) connect() error {
	mode := &serial.Mode{BaudRate: m.Cfg.BaudRate}

	port, err := serial.Open(m.Cfg.Port, mode)
	if err != nil {
		return err
	}

	err = port.SetDTR(true) // REQUIRED!
	if err != nil {
		err = port.Close()
		if err != nil {
			return err
		}
		return err
	}
	err = port.SetRTS(true)
	if err != nil {
		err = port.Close()
		if err != nil {
			return err
		}
		return err
	}

	time.Sleep(100 * time.Millisecond) // Wait for the device to reset
	_, err = Write(port, ACKPacket)    // Send ACK to indicate we're ready
	if err != nil {
		err = port.Close()
		if err != nil {
			return err
		}
		return err
	}

	fmt.Printf("Connected to %s at %d baud\n", m.Cfg.Port, m.Cfg.BaudRate)
	m.Port = port
	return nil
}

func (m *Manager) readLoop() {
	buf := make([]byte, 1024)

	for {
		select {
		case <-m.ctx.Done():
			return
		default:
			n, err := m.Port.Read(buf)
			if err != nil {
				err = m.connect()
				if err != nil {
					time.Sleep(time.Second)
					continue
				}
				return
			}

			data := make([]byte, n)
			copy(data, buf[:n])

			m.RxChan <- data
			fmt.Printf("Read %d bytes: %x\n", n, data)
		}
	}
}

func (m *Manager) writeLoop() {
	for {
		select {
		case <-m.ctx.Done():
			return
		case msg := <-m.TxChan:
			_, err := m.Port.Write(msg)
			if err != nil {
				return
			}
			fmt.Printf("Wrote %d bytes: %x\n", len(msg), msg)
		}
	}
}

func (m *Manager) Run() {
	for {
		err := m.connect()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		go m.readLoop()
		go m.writeLoop()

		<-m.ctx.Done()
		return
	}
}
