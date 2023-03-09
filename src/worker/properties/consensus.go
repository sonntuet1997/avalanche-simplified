package properties

import "gitlab.com/golibs-starter/golib/config"

type ConsensusProperties struct {
	N     int
	K     int
	Alpha int
	Beta  int
}

func NewConsensusProperties(loader config.Loader) (*ConsensusProperties, error) {
	props := ConsensusProperties{}
	err := loader.Bind(&props)
	if err != nil {
		return nil, err
	}
	return &props, err
}

func (t *ConsensusProperties) Prefix() string {
	return "app.consensus"
}
