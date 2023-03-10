package properties

import (
	"gitlab.com/golibs-starter/golib/config"
)

type RandomProperties struct {
	Range int
}

func NewRandomProperties(loader config.Loader) (*RandomProperties, error) {
	props := RandomProperties{}
	err := loader.Bind(&props)
	if err != nil {
		return nil, err
	}
	return &props, err
}

func (t *RandomProperties) Prefix() string {
	return "app.random"
}
