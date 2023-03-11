package jobs

import (
	"github.com/robfig/cron/v3"
	"github.com/sonntuet1997/avalanche-simplified/worker/services"
	"gitlab.com/golibs-starter/golib/log"
	"go.uber.org/fx"
)

type ScanNodesJob struct {
	P2pService *services.P2pService
}

type ScanNodeJobParams struct {
	fx.In
	P2pService *services.P2pService
}

func NewScanNodesJob(params ScanNodeJobParams) cron.Job {
	return &ScanNodesJob{
		P2pService: params.P2pService,
	}
}

func (r *ScanNodesJob) Run() {
	log.Debugf("[ScanNodesJob] job start")
	err := r.P2pService.ScanNodes()
	if err != nil {
		log.Errorf("[ScanNodesJob] failed to ScanNodes with error: %+v", err)
	}
	log.Debugf("[ScanNodesJob] job stop")
}
