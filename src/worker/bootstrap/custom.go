package bootstrap

import (
	"github.com/sonntuet1997/avalanche-simplyfied/worker/jobs"
	golibcron "gitlab.com/golibs-starter/golib-cron"
	"go.uber.org/fx"
)

func Custom() fx.Option {
	return fx.Options(
		golibcron.ProvideJob(jobs.NewSelfIntroductionJob),
		P2pOpt(),
	)
}
