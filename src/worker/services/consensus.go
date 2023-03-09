package services

import (
	"github.com/sonntuet1997/avalanche-simplified/worker/entities"
	"github.com/sonntuet1997/avalanche-simplified/worker/properties"
)

type ConsensusService struct {
	ConsensusProperties          *properties.ConsensusProperties
	ConfirmedTransactions        map[string]*entities.Transaction
	RootConfirmedTransactionNode *entities.TransactionNode
}

func NewConsensusService(
	ConsensusProperties *properties.ConsensusProperties,
) *ConsensusService {
	return &ConsensusService{
		ConsensusProperties: ConsensusProperties,
	}
}

func (c *ConsensusService) CountVote(transaction *entities.TransactionNode) {

}

func (c *ConsensusService) PrintConfirmedTree(transaction *entities.TransactionNode) {

}
