package serial

import (
	"context"
	"fmt"
	"github.com/Rokkit-exe/deckctl/config"
	"go.bug.st/serial"
	"io"
	"sync"
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

	wg sync.WaitGroup
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

	if err := port.SetDTR(true); err != nil {
		_ = port.Close()
		return err
	}

	if err := port.SetRTS(true); err != nil {
		_ = port.Close()
		return err
	}

	time.Sleep(100 * time.Millisecond)

	if _, err := Write(port, ACKPacket); err != nil {
		_ = port.Close()
		return err
	}

	fmt.Printf("Connected to %s at %d baud\n", m.Cfg.Port, m.Cfg.BaudRate)

	m.Port = port
	return nil
}

func (m *Manager) readLoop() {
	defer m.wg.Done()

	buf := make([]byte, 1024)

	for {
		n, err := m.Port.Read(buf)
		if err != nil {
			select {
			case <-m.ctx.Done():
				return
			default:
			}
			fmt.Println("Read error:", err)
			return
		}

		if n > 0 {
			data := make([]byte, n)
			copy(data, buf[:n])

			select {
			case m.RxChan <- data:
			case <-m.ctx.Done():
				return
			}
		}
	}
}

func (m *Manager) writeLoop() {
	defer m.wg.Done()

	for {
		select {
		case <-m.ctx.Done():
			return

		case msg, ok := <-m.TxChan:
			if !ok {
				return
			}

			if m.Port == nil {
				continue
			}

			_, err := m.Port.Write(msg)
			if err != nil {
				fmt.Println("Write error:", err)
				return
			}
		}
	}
}

func (m *Manager) Run() {
	for {
		select {
		case <-m.ctx.Done():
			return
		default:
		}

		if err := m.connect(); err != nil {
			fmt.Println("Connect error:", err)
			time.Sleep(time.Second)
			continue
		}

		m.wg.Add(2)
		go m.readLoop()
		go m.writeLoop()

		done := make(chan struct{})
		go func() {
			m.wg.Wait()
			close(done)
		}()

		select {
		case <-m.ctx.Done():
		case <-done:
		}

		// Cleanup
		if m.Port != nil {
			_ = m.Port.Close()
			m.Port = nil
		}

		// If shutting down, exit fully
		select {
		case <-m.ctx.Done():
			return
		default:
		}

		// Otherwise reconnect
		fmt.Println("Reconnecting...")
		time.Sleep(time.Second)
	}
}

func (m *Manager) Stop() {
	m.cancel()

	if m.Port != nil {
		_ = m.Port.Close()
	}

	m.wg.Wait()

	close(m.RxChan)
	close(m.TxChan)

	fmt.Println("Serial manager stopped")
}
