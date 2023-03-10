package tests

import (
	"context"
	"gitlab.com/golibs-starter/golib"
	"go.uber.org/fx"
)

func init() {
	err := fx.New(
		golib.PropertiesOpt(),
		golib.ProvideProps(NewTestingProperties),
		golib.ProvidePropsOption(golib.WithActiveProfiles([]string{"default"})),
		golib.ProvidePropsOption(golib.WithPaths([]string{"."})),
		fx.Invoke(func(client *TestingProperties) {
			testingProperties = client
		}),
	).Start(context.Background())
	if err != nil {
		panic(err)
	}
}

var testingProperties *TestingProperties
