package tests

import (
	"context"
	"github.com/sonntuet1997/avalanche-simplified/worker/bootstrap"
	"github.com/sonntuet1997/avalanche-simplified/worker/services"
	"gitlab.com/golibs-starter/golib"
	golibcrontestsuite "gitlab.com/golibs-starter/golib-cron/testsuite"
	golibtest "gitlab.com/golibs-starter/golib-test"
	"go.uber.org/fx"
)

func init() {
	err := fx.New(
		bootstrap.All(),
		golib.ProvidePropsOption(golib.WithActiveProfiles([]string{"testing"})),
		golib.ProvidePropsOption(golib.WithPaths([]string{"../config/"})),
		golibtest.EnableWebTestUtil(),
		golibcrontestsuite.EnableCronTestSuite(),
		fx.Invoke(func(client *services.P2pService) {
			P2pService = client
		}),
	).Start(context.Background())
	if err != nil {
		panic(err)
	}
}

var P2pService *services.P2pService
