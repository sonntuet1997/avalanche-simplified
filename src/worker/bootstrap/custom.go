package bootstrap

import (
	"github.com/sonntuet1997/avalanche-simplified/worker/controllers"
	"github.com/sonntuet1997/avalanche-simplified/worker/jobs"
	"github.com/sonntuet1997/avalanche-simplified/worker/properties"
	"github.com/sonntuet1997/avalanche-simplified/worker/services"
	"gitlab.com/golibs-starter/golib"
	golibcron "gitlab.com/golibs-starter/golib-cron"
	"go.uber.org/fx"
)

func Custom() fx.Option {
	return fx.Options(
		golibcron.ProvideJob(jobs.NewSelfIntroductionJob),
		golib.ProvideProps(properties.NewConsensusProperties),
		fx.Provide(services.NewConsensusService),
		fx.Provide(controllers.NewNodeController),
		P2pOpt(),
	)
}
