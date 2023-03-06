package services

import (
	"context"
	"fmt"
	"github.com/sonntuet1997/avalanche-simplyfied/worker/entities"
	"github.com/sonntuet1997/avalanche-simplyfied/worker/properties"
	"gitlab.com/golibs-starter/golib/log"
	"net"
)

type P2pService struct {
	P2pProperties  *properties.P2pProperties
	NeighborNodes  map[string]*entities.Node // address -> node
	CancelFunction *context.CancelFunc
}

func NewP2pService(
	P2pProperties *properties.P2pProperties,
) *P2pService {
	return &P2pService{
		P2pProperties: P2pProperties,
		NeighborNodes: make(map[string]*entities.Node, 0),
	}
}

const (
	broadcastAddrTmp = "255.255.255.255:%d"
	protocol         = "udp"
	portTmp          = ":%d"
)

func (p *P2pService) SelfIntroduce() error {
	addr, err := net.ResolveUDPAddr(protocol, fmt.Sprintf(broadcastAddrTmp, p.P2pProperties.Port))
	if err != nil {
		return fmt.Errorf("failed to ResolveUDPAddr with error: %w", err)
	}
	conn, err := net.DialUDP(protocol, nil, addr)
	if err != nil {
		return fmt.Errorf("failed to DialUDP with error: %w", err)
	}
	defer conn.Close()
	message := "."
	_, err = conn.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("failed to write message with error: %w", err)
	}
	return nil
}

func (p *P2pService) ListenForBroadcasts(ctx context.Context) {
	addr, err := net.ResolveUDPAddr(protocol, fmt.Sprintf(portTmp, p.P2pProperties.Port))
	if err != nil {
		log.Errorf("failed to ResolveUDPAddr with error: %w", err)
		return
	}
	conn, err := net.ListenUDP(protocol, addr)
	if err != nil {
		log.Errorf("failed to ListenUDP with error: %w", err)
		return
	}
	defer conn.Close()
	ctx, cancelFunc := context.WithCancel(ctx)
	p.CancelFunction = &cancelFunc
	opening := true
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Debugf("[P2pService] Received Done Signal")
				opening = false
				return
			}
		}
	}()
	for opening {
		buf := make([]byte, 1024)
		n, nodeAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Errorf("Error reading UDP message: %+v", err)
			continue
		}
		log.Debugf("Received broadcast from %+v %+v", nodeAddr.String(), string(buf[:n]))
		p.NeighborNodes[nodeAddr.String()] = &entities.Node{
			Address: nodeAddr.String(),
		}
	}
	log.Debugf("[P2pService] Exiting ListenForBroadcasts")
}

func (p *P2pService) Close() error {
	log.Debugf("[P2pService] Closing")
	if p.CancelFunction == nil {
		return nil
	}
	(*p.CancelFunction)()
	return nil
}
