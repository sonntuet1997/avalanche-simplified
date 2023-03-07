package tests

import (
	"context"
	"github.com/sonntuet1997/avalanche-simplyfied/worker/bootstrap"
	"github.com/sonntuet1997/avalanche-simplyfied/worker/services"
	"gitlab.com/golibs-starter/golib"
	golibcrontestsuite "gitlab.com/golibs-starter/golib-cron/testsuite"
	golibtest "gitlab.com/golibs-starter/golib-test"
	"go.uber.org/fx"
)

func init() {
	// init node 1
	err := fx.New(
		bootstrap.All(),
		golib.ProvidePropsOption(golib.WithActiveProfiles([]string{"testing_node1"})),
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
	// init node 2
	err = fx.New(
		bootstrap.All(),
		golib.ProvidePropsOption(golib.WithActiveProfiles([]string{"testing_node2"})),
		golib.ProvidePropsOption(golib.WithPaths([]string{"../config/"})),
		golibtest.EnableWebTestUtil(),
		fx.Invoke(func(client *services.P2pService) {
			P2pService = client
		}),
	).Start(context.Background())
	if err != nil {
		panic(err)
	}
}

var P2pService *services.P2pService
