package services

import (
	"context"
	"fmt"
	"github.com/sonntuet1997/avalanche-simplified/worker/entities"
	"github.com/sonntuet1997/avalanche-simplified/worker/properties"
	"gitlab.com/golibs-starter/golib/log"
	"math/rand"
	"net"
	"sync"
	"time"
)

type P2pService struct {
	P2pProperties     *properties.P2pProperties
	NeighborNodes     map[string]*entities.Node // address -> node
	CancelFunction    *context.CancelFunc
	CachedRandomNodes []*entities.Node
	Wg                sync.WaitGroup
}

func NewP2pService(
	P2pProperties *properties.P2pProperties,
) *P2pService {
	service := P2pService{
		P2pProperties: P2pProperties,
		NeighborNodes: make(map[string]*entities.Node, 0),
	}
	service.Wg.Add(P2pProperties.MinConnectedNodes)
	return &service
}

const (
	broadcastAddrTmp = "255.255.255.255:%d"
	protocol         = "udp"
	portTmp          = ":%d"
)

func (p *P2pService) SelfIntroduce() error {
	addr, err := net.ResolveUDPAddr(protocol, fmt.Sprintf(broadcastAddrTmp, p.P2pProperties.BroadcastPort))
	if err != nil {
		return fmt.Errorf("failed to ResolveUDPAddr with error: %w", err)
	}
	conn, err := net.DialUDP(protocol, nil, addr)
	if err != nil {
		return fmt.Errorf("failed to DialUDP with error: %w", err)
	}
	defer conn.Close()
	message := "+"
	_, err = conn.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("failed to write message with error: %w", err)
	}
	return nil
}

func (p *P2pService) GetRandomNodes(nodesNumber int) ([]*entities.Node, error) {
	if nodesNumber > p.P2pProperties.MinConnectedNodes {
		panic("wrong config!")
	}
	p.Wg.Wait()
	if p.CachedRandomNodes != nil && len(p.CachedRandomNodes) == nodesNumber {
		return p.CachedRandomNodes, nil
	}
	neighborNodes := make([]*entities.Node, 0, len(p.NeighborNodes))
	for _, v := range p.NeighborNodes {
		neighborNodes = append(neighborNodes, v)
	}

	rand.Seed(time.Now().UnixNano())
	randomElements := make([]*entities.Node, nodesNumber)

	uniqueIndices := make(map[int]bool)
	for len(uniqueIndices) < nodesNumber {
		uniqueIndices[rand.Intn(len(p.NeighborNodes))] = true
	}
	i := 0
	for index := range uniqueIndices {
		randomElements[i] = neighborNodes[index]
		i++
	}
	// TODO: not safe for concurrency
	p.CachedRandomNodes = randomElements
	return randomElements, nil
}

func (p *P2pService) SendLeavingSignals(numberSignals int) error {
	addr, err := net.ResolveUDPAddr(protocol, fmt.Sprintf(broadcastAddrTmp, p.P2pProperties.BroadcastPort))
	if err != nil {
		return fmt.Errorf("failed to ResolveUDPAddr with error: %w", err)
	}
	conn, err := net.DialUDP(protocol, nil, addr)
	if err != nil {
		return fmt.Errorf("failed to DialUDP with error: %w", err)
	}
	defer conn.Close()
	message := "-"
	for i := 0; i < numberSignals; i++ {
		_, err = conn.Write([]byte(message))
		if err != nil {
			return fmt.Errorf("failed to write message with error: %w", err)
		}
		time.Sleep(time.Second)
	}
	return nil
}

func (p *P2pService) ListenForBroadcasts(ctx context.Context) {
	addr, err := net.ResolveUDPAddr(protocol, fmt.Sprintf(portTmp, p.P2pProperties.BroadcastPort))
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
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Debugf("[P2pService] Received Done Signal")
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
				if _, ok := p.NeighborNodes[nodeAddr.String()]; !ok {
					p.NeighborNodes[nodeAddr.String()] = &entities.Node{
						Address: nodeAddr.String(),
					}
					p.Wg.Done()
				}
			}
		}
	}()
	<-ctx.Done()
	log.Debugf("[P2pService] Leaving")
}

func (p *P2pService) Close() error {
	log.Debugf("[P2pService] Closing")
	if p.CancelFunction == nil {
		return nil
	}
	(*p.CancelFunction)()
	return nil
}
