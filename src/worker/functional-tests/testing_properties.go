package tests

import (
	"fmt"
	"gitlab.com/golibs-starter/golib/config"
	"time"
)

type TestingProperties struct {
	ToBlockNumber int
	NodesNumber   int
	URL           string
	Timeout       time.Duration `mapstructure:"_"`
	TimeoutString string        `mapstructure:"timeout"`
}

func NewTestingProperties(loader config.Loader) (*TestingProperties, error) {
	props := TestingProperties{}
	err := loader.Bind(&props)
	if err != nil {
		return nil, err
	}
	timeout, err := time.ParseDuration(props.TimeoutString)
	if err != nil {
		return nil, fmt.Errorf("failed to ParseDuration TimeoutString with err: %w", err)
	}
	props.Timeout = timeout
	return &props, err
}

func (t *TestingProperties) Prefix() string {
	return "app.testing"
}
