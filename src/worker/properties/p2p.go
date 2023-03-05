package properties

import "gitlab.com/golibs-starter/golib/config"

type P2pProperties struct {
	Port int
}

func NewP2pProperties(loader config.Loader) (*P2pProperties, error) {
	props := P2pProperties{}
	err := loader.Bind(&props)
	if err != nil {
		return nil, err
	}
	if props.Port < 1000 || props.Port > 65535 {
		panic("invalid port number")
	}
	return &props, err
}

func (t *P2pProperties) Prefix() string {
	return "app.p2p"
}
