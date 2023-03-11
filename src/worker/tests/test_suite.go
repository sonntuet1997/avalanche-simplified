package tests

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/sonntuet1997/avalanche-simplified/worker/bootstrap"
	"github.com/sonntuet1997/avalanche-simplified/worker/services"
	"gitlab.com/golibs-starter/golib"
	golibcrontestsuite "gitlab.com/golibs-starter/golib-cron/testsuite"
	golibtest "gitlab.com/golibs-starter/golib-test"
	"go.uber.org/fx"
	"net/http"
)

func init() {
	err := fx.New(
		bootstrap.All(),
		golib.ProvidePropsOption(golib.WithActiveProfiles([]string{"testing"})),
		golib.ProvidePropsOption(golib.WithPaths([]string{"../config/"})),
		golibtest.EnableWebTestUtil(),
		golibcrontestsuite.EnableCronTestSuite(),
		fx.Invoke(func(client *services.P2pService) {
			p2pService = client
		}),
		fx.Invoke(func(client *services.ConsensusService) {
			consensusService = client
		}),
		fx.Invoke(func(client *http.Client) {
			httpClient = client
		}),
		fx.Invoke(func(client *resty.Client) {
			restyClient = client
		}),
	).Start(context.Background())
	if err != nil {
		panic(err)
	}
}

var p2pService *services.P2pService
var consensusService *services.ConsensusService
var httpClient *http.Client
var restyClient *resty.Client
