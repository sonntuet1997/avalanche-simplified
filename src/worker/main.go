package main

import (
	"github.com/sonntuet1997/avalanche-simplyfied/worker/bootstrap"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		bootstrap.All(),
	).Run()
}
