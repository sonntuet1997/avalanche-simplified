package mocks

import (
	"context"
	"fmt"
	"gitlab.com/golibs-starter/golib/log"
	"net"
	"time"
)

type BroadcastListener struct {
	Port           int
	CancelFunction *context.CancelFunc
	InputChannel   chan string
}

func NewBroadcastListener(port int) *BroadcastListener {
	return &BroadcastListener{
		Port:         port,
		InputChannel: make(chan string, 0),
	}
}

const (
	protocol = "udp"
	portTmp  = ":%d"
)

func (p *BroadcastListener) ListenForBroadcasts(ctx context.Context) error {
	addr, err := net.ResolveUDPAddr(protocol, fmt.Sprintf(portTmp, p.Port))
	if err != nil {
		return fmt.Errorf("failed to ResolveUDPAddr with error: %w", err)
	}
	conn, err := net.ListenUDP(protocol, addr)
	if err != nil {
		return fmt.Errorf("failed to ListenUDP with error: %w", err)
	}
	ctx, cancelFunc := context.WithCancel(ctx)
	p.CancelFunction = &cancelFunc
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Debugf("[BroadcastListener] Received Done Signal")
				conn.Close()
				return
			default:
				buf := make([]byte, 1024)
				err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
				if err != nil {
					log.Errorf("Error SetReadDeadline: %+v", err)
					return
				}
				n, nodeAddr, err := conn.ReadFromUDP(buf)
				if err != nil {
					if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
						continue
					}
					log.Errorf("Error reading UDP message: %+v", err)
					continue
				}
				log.Debugf("Received broadcast from %+v %+v", nodeAddr.String(), string(buf[:n]))
				p.InputChannel <- nodeAddr.String()
			}
		}
	}()
	return nil
}

func (p *BroadcastListener) Close() error {
	log.Debugf("[BroadcastListener] Closing")
	if p.CancelFunction == nil {
		return nil
	}
	(*p.CancelFunction)()
	return nil
}
