package services

import (
	"context"
	"fmt"
	"github.com/deliveryhero/pipeline/v2"
	"github.com/sonntuet1997/avalanche-simplified/worker/entities"
	"github.com/sonntuet1997/avalanche-simplified/worker/properties"
	"gitlab.com/golibs-starter/golib/log"
)

type ConsensusService struct {
	ConsensusProperties          *properties.ConsensusProperties
	ConfirmedTransactions        map[string]*entities.Transaction
	RootConfirmedTransactionNode *entities.TransactionNode
	P2pService                   *P2pService
}

func NewConsensusService(
	ConsensusProperties *properties.ConsensusProperties,
	P2pService *P2pService,
) *ConsensusService {
	consensusService := &ConsensusService{
		ConsensusProperties:   ConsensusProperties,
		P2pService:            P2pService,
		ConfirmedTransactions: make(map[string]*entities.Transaction, 0),
	}
	genesisTransaction := &entities.TransactionNode{
		Transaction: &entities.Transaction{
			ID:    "0-0",
			Major: 0,
			Minor: 0,
		},
	}
	consensusService.RootConfirmedTransactionNode = genesisTransaction
	consensusService.ConfirmedTransactions[genesisTransaction.Transaction.ID] = genesisTransaction.Transaction
	return consensusService
}

func (c *ConsensusService) AskNeighborsForPreferences(preferredTransactionNode *entities.TransactionNode) int {
	ctx := context.Background()
	nodes, err := c.P2pService.GetRandomNodes(c.ConsensusProperties.K)
	if err != nil {
		return 0
	}
	processor := pipeline.NewProcessor(
		func(ctx context.Context, input *entities.Node) (*entities.TransactionNode, error) {
			preference, err := c.askNeighborForPreference(preferredTransactionNode)
			if err != nil {
				return nil, fmt.Errorf("failed to askNeighborForPreference with error: %w", err)
			}
			return preference, nil
		}, func(i *entities.Node, err error) {
			log.Errorf("[ConsensusService] failed to process with error: %w", err)
		})
	processNeighborsPreferencesChan := pipeline.ProcessConcurrently(
		ctx,
		c.ConsensusProperties.K,
		processor,
		pipeline.Emit(nodes...),
	)
	for neighborPreference := range processNeighborsPreferencesChan {
		log.Debugf("%+v", neighborPreference)
		_, err = c.updateTransaction(neighborPreference)
		if err != nil {
			log.Errorf("[ConsensusService] failed to updateTransaction with error: %w", err)
			return 0
		}
	}
	return 0
}

func (c *ConsensusService) askNeighborForPreference(preferredTransactionNode *entities.TransactionNode) (*entities.TransactionNode, error) {
	return nil, nil
}

func (c *ConsensusService) updateTransaction(transactionNode *entities.TransactionNode) (*entities.TransactionNode, error) {
	return nil, nil
}

func (c *ConsensusService) PrintConfirmedTree(transaction *entities.TransactionNode) {

}

func (c *ConsensusService) MyPreference(askedTransactionNode *entities.TransactionNode) (*entities.TransactionNode, error) {
	return nil, nil
}
