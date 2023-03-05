package services

import (
	"fmt"
	"github.com/sonntuet1997/avalanche-simplyfied/worker/properties"
	"gitlab.com/golibs-starter/golib/log"
	"net"
)

type P2pService struct {
	P2pProperties *properties.P2pProperties
}

func NewP2pService(
	P2pProperties *properties.P2pProperties,
) *P2pService {
	return &P2pService{
		P2pProperties: P2pProperties,
	}
}

const (
	broadcastAddrTmp = "255.255.255.255:%d"
	protocol         = "udp"
	portTmp          = ":%d"
)

func (p *P2pService) Broadcast() error {
	addr, err := net.ResolveUDPAddr(protocol, fmt.Sprintf(broadcastAddrTmp, p.P2pProperties.Port))
	if err != nil {
		return fmt.Errorf("failed to ResolveUDPAddr with error: %w", err)
	}
	conn, err := net.DialUDP(protocol, nil, addr)
	if err != nil {
		return fmt.Errorf("failed to DialUDP with error: %w", err)
	}
	defer conn.Close()

	message := "Hello from node"
	_, err = conn.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("failed to write message with error: %w", err)
	}
	return nil
}

func (p *P2pService) ListenForBroadcasts() error {
	addr, err := net.ResolveUDPAddr(protocol, fmt.Sprintf(portTmp, p.P2pProperties.Port))
	if err != nil {
		return fmt.Errorf("failed to ResolveUDPAddr with error: %w", err)
	}
	conn, err := net.ListenUDP(protocol, addr)
	if err != nil {
		return fmt.Errorf("failed to ListenUDP with error: %w", err)
	}
	defer conn.Close()

	// read incoming broadcast messages and print the source address
	buf := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Errorf("Error reading UDP message: %+v", err)
			continue
		}
		log.Debugf("Received broadcast from", addr.String(), string(buf[:n]))
	}
}
