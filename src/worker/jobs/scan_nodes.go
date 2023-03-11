package jobs

import (
	"github.com/robfig/cron/v3"
	"github.com/sonntuet1997/avalanche-simplified/worker/services"
	"gitlab.com/golibs-starter/golib/log"
	"go.uber.org/fx"
)

type ScanNodeJob struct {
	P2pService *services.P2pService
}

type ScanNodeJobParams struct {
	fx.In
	P2pService *services.P2pService
}

func NewScanNodeJob(params ScanNodeJobParams) cron.Job {
	return &ScanNodeJob{
		P2pService: params.P2pService,
	}
}

func (r *ScanNodeJob) Run() {
	log.Debugf("[ScanNodeJob] job start")
	err := r.P2pService.ScanNodes()
	if err != nil {
		log.Errorf("[ScanNodeJob] failed to ScanNodes with error: %+v", err)
	}
	log.Debugf("[ScanNodeJob] job stop")
}
