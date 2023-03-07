package bootstrap

import (
	"context"
	"github.com/sonntuet1997/avalanche-simplyfied/worker/properties"
	"github.com/sonntuet1997/avalanche-simplyfied/worker/services"
	"gitlab.com/golibs-starter/golib"
	"go.uber.org/fx"
)

func P2pOpt() fx.Option {
	return fx.Options(
		golib.ProvideProps(properties.NewP2pProperties),
		fx.Provide(services.NewP2pService),
		fx.Invoke(OnStartP2pService),
		fx.Invoke(OnStopP2pService),
	)
}

func OnStartP2pService(lc fx.Lifecycle, p2pService *services.P2pService, p2pProperties *properties.P2pProperties) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if p2pProperties.ListenToBroadcast {
				go p2pService.ListenForBroadcasts(context.Background())
			}
			return nil
		},
	})
}

func OnStopP2pService(lc fx.Lifecycle, p2pService *services.P2pService) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return p2pService.Close()
		},
	})
}
