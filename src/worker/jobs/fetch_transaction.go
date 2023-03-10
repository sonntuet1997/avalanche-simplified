package jobs

import (
	"errors"
	"github.com/robfig/cron/v3"
	"github.com/sonntuet1997/avalanche-simplified/worker/constants"
	"github.com/sonntuet1997/avalanche-simplified/worker/services"
	"gitlab.com/golibs-starter/golib/log"
	"go.uber.org/fx"
)

type FetchTransactionJob struct {
	ConsensusService *services.ConsensusService
}

type FetchTransactionJobParams struct {
	fx.In
	ConsensusService *services.ConsensusService
}

func NewFetchTransactionJob(params FetchTransactionJobParams) cron.Job {
	return &FetchTransactionJob{
		ConsensusService: params.ConsensusService,
	}
}

func (r *FetchTransactionJob) Run() {
	log.Debugf("[FetchTransactionJob] job start")
	err := r.ConsensusService.FetchNewTransaction()
	if errors.Is(err, constants.ErrNotEnoughNeighborNodes) {
		log.Warnf("[FetchTransactionJob] %+v", err)
	}
	if err != nil {
		log.Errorf("[FetchTransactionJob] failed to FetchNewTransaction with error: %+v", err)
	}
	log.Debugf("[FetchTransactionJob] job stop")
}
