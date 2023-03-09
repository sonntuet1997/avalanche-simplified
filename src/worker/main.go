package main

import (
	"github.com/sonntuet1997/avalanche-simplified/worker/bootstrap"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		bootstrap.All(),
	).Run()
}
