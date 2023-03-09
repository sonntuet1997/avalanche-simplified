package properties

import "gitlab.com/golibs-starter/golib/config"

type P2pProperties struct {
	BroadcastPort     int
	MinConnectedNodes int
	ListenToBroadcast bool `default:"true"`
}

func NewP2pProperties(loader config.Loader) (*P2pProperties, error) {
	props := P2pProperties{}
	err := loader.Bind(&props)
	if err != nil {
		return nil, err
	}
	if props.BroadcastPort < 1000 || props.BroadcastPort > 65535 {
		panic("invalid port number")
	}
	return &props, err
}

func (t *P2pProperties) Prefix() string {
	return "app.p2p"
}
