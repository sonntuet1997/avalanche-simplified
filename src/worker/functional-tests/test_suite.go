package tests

import (
	"context"
	"github.com/sonntuet1997/avalanche-simplyfied/worker/bootstrap"
	"gitlab.com/golibs-starter/golib"
	golibtest "gitlab.com/golibs-starter/golib-test"
	"go.uber.org/fx"
)

func init() {
	err := fx.New(
		bootstrap.All(),
		golib.ProvidePropsOption(golib.WithActiveProfiles([]string{"functional_testing"})),
		golib.ProvidePropsOption(golib.WithPaths([]string{"../config/"})),
		golibtest.EnableWebTestUtil(),
	).Start(context.Background())
	if err != nil {
		panic(err)
	}
}
