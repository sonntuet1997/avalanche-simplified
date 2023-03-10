package tests

import (
	"gitlab.com/golibs-starter/golib/config"
)

type TestingProperties struct {
	ToBlockNumber int
	NodesNumber   int
	URL           string
}

func NewTestingProperties(loader config.Loader) (*TestingProperties, error) {
	props := TestingProperties{}
	err := loader.Bind(&props)
	if err != nil {
		return nil, err
	}
	return &props, err
}

func (t *TestingProperties) Prefix() string {
	return "app.testing"
}
