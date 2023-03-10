package jobs

import (
	"github.com/robfig/cron/v3"
	"github.com/sonntuet1997/avalanche-simplified/worker/services"
	"gitlab.com/golibs-starter/golib/log"
	"go.uber.org/fx"
	"math/rand"
	"time"
)

type SelfIntroductionJob struct {
	P2pService *services.P2pService
}

type SelfIntroductionJobParams struct {
	fx.In
	P2pService *services.P2pService
}

func NewSelfIntroductionJob(params SelfIntroductionJobParams) cron.Job {
	return &SelfIntroductionJob{
		P2pService: params.P2pService,
	}
}

func (r *SelfIntroductionJob) Run() {
	log.Debugf("[SelfIntroductionJob] job start")
	rand.Seed(time.Now().UnixNano())
	randomValue := rand.Intn(10) + 1
	time.Sleep(time.Second * (time.Duration)(randomValue))
	err := r.P2pService.SelfIntroduce()
	if err != nil {
		log.Errorf("[SelfIntroductionJob] failed to SelfIntroduce with error: %+v", err)
	}
	log.Debugf("[SelfIntroductionJob] job stop")
}
