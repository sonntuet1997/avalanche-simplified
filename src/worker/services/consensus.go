package services

import (
	"context"
	"fmt"
	"github.com/deliveryhero/pipeline/v2"
	"github.com/sonntuet1997/avalanche-simplified/worker/entities"
	"github.com/sonntuet1997/avalanche-simplified/worker/properties"
	http_client "github.com/sonntuet1997/avalanche-simplified/worker/repositories/http-client"
	"gitlab.com/golibs-starter/golib/log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type ConsensusService struct {
	P2pService          *P2pService
	ConsensusProperties *properties.ConsensusProperties
	RandomProperties    *properties.RandomProperties
	NodeRepository      *http_client.NodeRepository

	ConfirmedTransactions        []*entities.Transaction
	CurrentBlockNumber           atomic.Int64
	CurrentPreferenceTransaction *entities.Transaction
	ConsecutiveSuccesses         atomic.Int64

	Mutex sync.RWMutex
}

func NewConsensusService(
	ConsensusProperties *properties.ConsensusProperties,
	RandomProperties *properties.RandomProperties,
	P2pService *P2pService,
	NodeRepository *http_client.NodeRepository,
) *ConsensusService {
	rand.Seed(time.Now().UnixNano())
	consensusService := &ConsensusService{
		ConsensusProperties:   ConsensusProperties,
		RandomProperties:      RandomProperties,
		P2pService:            P2pService,
		NodeRepository:        NodeRepository,
		ConfirmedTransactions: make([]*entities.Transaction, 0),
	}
	consensusService.ConfirmedTransactions = append(consensusService.ConfirmedTransactions, &entities.Transaction{
		ID:    "0-0",
		Major: 0,
		Minor: 0,
	})
	consensusService.CurrentPreferenceTransaction = consensusService.MakeRandomTransaction(1)
	return consensusService
}

func (c *ConsensusService) MyPreference(blockNumber int) (*entities.Transaction, error) {
	currentBlockNumber := c.CurrentBlockNumber.Load()
	if currentBlockNumber >= int64(blockNumber) {
		return c.ConfirmedTransactions[blockNumber], nil
	}
	if int64(blockNumber) == currentBlockNumber+1 {
		c.Mutex.RLock()
		defer c.Mutex.RUnlock()
		return c.CurrentPreferenceTransaction, nil
	}
	invalidTransaction := &entities.Transaction{
		ID:    "-1",
		Major: -1,
		Minor: 0,
	}
	return invalidTransaction, nil
}

func (c *ConsensusService) FetchNewTransaction() error {
	nextBlockNumber := int(c.CurrentBlockNumber.Load()) + 1
	preferences, err := c.AskNeighborsForPreferences(nextBlockNumber)
	if err != nil {
		return fmt.Errorf("failed to AskNeighborsForPreferences with err: %w", err)
	}
	err = c.UpdatePreferredTransaction(nextBlockNumber, preferences)
	if err != nil {
		return fmt.Errorf("failed to UpdatePreferredTransaction with err: %w", err)
	}
	return nil
}

func (c *ConsensusService) AskNeighborsForPreferences(blockNumber int) ([]*entities.Transaction, error) {
	ctx := context.Background()
	nodes, err := c.P2pService.GetRandomNodes(c.ConsensusProperties.K)
	if err != nil {
		return nil, fmt.Errorf("failed to GetRandomNodes with err: %w", err)
	}
	processor := pipeline.NewProcessor(
		func(ctx context.Context, neighborNode *entities.Node) (*entities.Transaction, error) {
			preference, err := c.askNeighborForPreference(neighborNode, blockNumber)
			if err != nil {
				return nil, fmt.Errorf("failed to askNeighborForPreference with error: %w", err)
			}
			return preference, nil
		}, func(i *entities.Node, err error) {
			log.Errorf("[ConsensusService] failed to process with error: %w", err)
		})
	neighborsPreferencesChan := pipeline.ProcessConcurrently(
		ctx,
		c.ConsensusProperties.K,
		processor,
		pipeline.Emit(nodes...),
	)
	collectedNeighborsPreferencesChan := pipeline.Collect(ctx, c.ConsensusProperties.K, c.ConsensusProperties.Timeout, neighborsPreferencesChan)
	collectedNeighborsPreferences := <-collectedNeighborsPreferencesChan
	return collectedNeighborsPreferences, err
}

func (c *ConsensusService) askNeighborForPreference(neighborNode *entities.Node, blockNumber int) (*entities.Transaction, error) {
	ctx := context.Background()
	preference, err := c.NodeRepository.AskForPreference(ctx, neighborNode.Address, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to AskForPreference with err: %w", err)
	}
	return preference, nil
}

func (c *ConsensusService) UpdatePreferredTransaction(blockNumber int, collectedNeighborsPreferences []*entities.Transaction) error {
	transactions := make(map[string]*entities.Transaction, 0)
	counts := make(map[string]int, 0)
	for _, neighborPreference := range collectedNeighborsPreferences {
		if neighborPreference.Major != blockNumber {
			log.Warnf("[UpdatePreferredTransaction] block number not matched %+v %+v", neighborPreference, blockNumber)
		}
		transactions[neighborPreference.ID] = neighborPreference
		counts[neighborPreference.ID]++
	}
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	hasMajor := false
	for k, v := range transactions {
		if counts[k] >= c.ConsensusProperties.Alpha {
			if v.ID == c.CurrentPreferenceTransaction.ID {
				c.ConsecutiveSuccesses.Add(1)
			} else {
				c.ConsecutiveSuccesses.Store(1)
			}
			c.CurrentPreferenceTransaction = v
			hasMajor = true
			break
		}
	}
	if !hasMajor {
		c.ConsecutiveSuccesses.Store(0)
		return nil
	}
	if c.ConsecutiveSuccesses.Load() >= c.ConsensusProperties.Beta {
		c.ConsecutiveSuccesses.Store(0)
		c.ConfirmedTransactions = append(c.ConfirmedTransactions, c.CurrentPreferenceTransaction)
		c.CurrentPreferenceTransaction = c.MakeRandomTransaction(blockNumber + 1)
		c.CurrentBlockNumber.Add(1)
	}
	return nil
}

func (c *ConsensusService) MakeRandomTransaction(blockNumber int) *entities.Transaction {
	randomValue := rand.Intn(c.RandomProperties.Range) + 1
	return &entities.Transaction{
		ID:    fmt.Sprintf("%d-%d", blockNumber, randomValue),
		Major: blockNumber,
		Minor: randomValue,
	}
}

func (c *ConsensusService) GetConfirmedTransactions() []*entities.Transaction {
	return c.ConfirmedTransactions
}
