package properties

import (
	"fmt"
	"gitlab.com/golibs-starter/golib/config"
	"time"
)

type ConsensusProperties struct {
	K       int
	Alpha   int
	Beta    int64
	Timeout time.Duration `mapstructure:"_"`

	TimeoutString string `mapstructure:"timeout"`
}

func NewConsensusProperties(loader config.Loader) (*ConsensusProperties, error) {
	props := ConsensusProperties{}
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

func (t *ConsensusProperties) Prefix() string {
	return "app.consensus"
}
