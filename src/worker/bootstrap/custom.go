package bootstrap

import (
	"github.com/sonntuet1997/avalanche-simplyfied/worker/properties"
	"github.com/sonntuet1997/avalanche-simplyfied/worker/services"
	"gitlab.com/golibs-starter/golib"
	"go.uber.org/fx"
)

func Custom() fx.Option {
	return fx.Options(
		golib.ProvideProps(properties.NewP2pProperties),
		fx.Provide(services.NewP2pService),
	)
}
