package models

import (
	"github.com/sonntuet1997/avalanche-simplified/worker/entities"
	"gitlab.com/golibs-starter/golib/web/response"
)

type PreferenceResponse struct {
	Meta response.Meta        `json:"meta"`
	Data entities.Transaction `json:"data"`
}
