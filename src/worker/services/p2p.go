package services

import (
	"context"
	"fmt"
	"github.com/deliveryhero/pipeline/v2"
	"github.com/sonntuet1997/avalanche-simplified/worker/constants"
	"github.com/sonntuet1997/avalanche-simplified/worker/entities"
	"github.com/sonntuet1997/avalanche-simplified/worker/properties"
	http_client "github.com/sonntuet1997/avalanche-simplified/worker/repositories/http-client"
	"gitlab.com/golibs-starter/golib/config"
	"gitlab.com/golibs-starter/golib/log"
	"math/rand"
	"net"
	"sync"
	"time"
)

type P2pService struct {
	P2pProperties  *properties.P2pProperties
	NeighborNodes  map[string]*entities.Node // address -> node
	NodeRepository *http_client.NodeRepository
	AppProperties  *config.AppProperties
	CancelFunction *context.CancelFunc
	LocalAddresses map[string]interface{}
	RWMutex        sync.RWMutex
}

func NewP2pService(
	P2pProperties *properties.P2pProperties,
	AppProperties *config.AppProperties,
	NodeRepository *http_client.NodeRepository,
) *P2pService {
	service := P2pService{
		P2pProperties:  P2pProperties,
		NodeRepository: NodeRepository,
		AppProperties:  AppProperties,
		NeighborNodes:  make(map[string]*entities.Node, 0),
	}
	service.getLocalAddresses()
	return &service
}

func (p *P2pService) getLocalAddresses() {
	result := make(map[string]interface{})
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
	}
	for _, addr := range addresses {
		ipnet, ok := addr.(*net.IPNet)
		if ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			result[ipnet.IP.String()] = struct{}{}
		}
	}
	p.LocalAddresses = result
}

const (
	broadcastAddrTmp = "255.255.255.255:%d"
	protocol         = "udp"
	portTmp          = ":%d"
)

func (p *P2pService) SelfIntroduce() error {
	if p.P2pProperties.DisableBroadcast {
		log.Debugf("DisableBroadcast SelfIntroduce")
		return nil
	}
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

func (p *P2pService) ScanNodes() error {
	p.RWMutex.Lock()
	totalCurrentNodes := len(p.NeighborNodes)
	p.RWMutex.Unlock()
	if totalCurrentNodes > p.P2pProperties.MinConnectedNodes {
		return nil
	}
	ctx := context.Background()
	processor := pipeline.NewProcessor(
		func(ctx context.Context, nodeNumber int) (string, error) {
			address, err := p.NodeRepository.CheckHealthAndGetAddress(ctx, fmt.Sprintf(p.P2pProperties.NodeHealthURLTemplate, nodeNumber, p.AppProperties.Port))
			if err != nil {
				return "", fmt.Errorf("failed to CheckHealthAndGetAddress with error: %w", err)
			}
			log.Infof("[P2pService] scanned address: %+v", address)
			time.Sleep(50 * time.Millisecond)
			return address, nil
		}, func(nodeNumber int, err error) {
			log.Errorf("[P2pService] failed to process node %+v with error: %w", nodeNumber, err)
		})
	inputChan := make(chan int)
	go func() {
		for i := 1; i <= p.P2pProperties.TotalNodes; i++ {
			inputChan <- i
		}
		close(inputChan)
	}()
	neighborsAddressesChan := pipeline.Process(
		ctx,
		processor,
		inputChan,
	)
	collectedNeighborsAddressesChan := pipeline.Collect(ctx, p.P2pProperties.TotalNodes, 5*time.Minute, neighborsAddressesChan)
	collectedNeighborsAddresses := <-collectedNeighborsAddressesChan
	neighborNodes := make(map[string]*entities.Node, p.P2pProperties.TotalNodes)
	for _, neighborsAddress := range collectedNeighborsAddresses {
		log.Debugf("[P2pService] neighborNodes %+v", neighborsAddress)
		_, ok := neighborNodes[neighborsAddress]
		if !ok {
			neighborNodes[neighborsAddress] = &entities.Node{
				Address: neighborsAddress,
			}
		}
	}
	p.RWMutex.Lock()
	p.NeighborNodes = neighborNodes
	p.RWMutex.Unlock()
	return nil
}

func (p *P2pService) GetRandomNodes(nodesNumber int) ([]*entities.Node, error) {
	if nodesNumber > p.P2pProperties.MinConnectedNodes {
		panic("wrong config!")
	}
	if nodesNumber > len(p.NeighborNodes) {
		return nil, constants.ErrNotEnoughNeighborNodes
	}
	neighborNodes := p.GetNeighborNodes()

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
	return randomElements, nil
}

func (p *P2pService) GetNeighborNodes() []*entities.Node {
	p.RWMutex.RLock()
	defer p.RWMutex.RUnlock()
	neighborNodes := make([]*entities.Node, 0, len(p.NeighborNodes))
	for _, v := range p.NeighborNodes {
		neighborNodes = append(neighborNodes, v)
	}
	return neighborNodes
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
	if p.P2pProperties.DisableBroadcast {
		log.Debugf("DisableBroadcast ListenForBroadcasts")
		return
	}
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
				if _, ok := p.LocalAddresses[nodeAddr.IP.String()]; ok {
					log.Debugf("Filtered your own broadcast message from %+v %+v", nodeAddr.String(), string(buf[:n]))
					continue
				}
				log.Debugf("Received broadcast from %+v %+v", nodeAddr.String(), string(buf[:n]))
				p.RWMutex.RLock()
				_, ok := p.NeighborNodes[nodeAddr.IP.String()]
				p.RWMutex.RUnlock()
				if !ok {
					p.RWMutex.Lock()
					p.NeighborNodes[nodeAddr.IP.String()] = &entities.Node{
						Address: nodeAddr.IP.String(),
					}
					p.RWMutex.Unlock()
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
