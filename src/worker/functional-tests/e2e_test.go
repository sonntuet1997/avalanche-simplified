package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ecodia/golang-awaitility/awaitility"
	"github.com/go-resty/resty/v2"
	"github.com/sonntuet1997/avalanche-simplified/worker/entities"
	"github.com/sonntuet1997/avalanche-simplified/worker/repositories/http-client/models"
	"github.com/stretchr/testify/assert"
	"gitlab.com/golibs-starter/golib/log"
	"net/http"
	"testing"
	"time"
)

func TestE2e(t *testing.T) {
	client := &http.Client{}
	restyClient := resty.NewWithClient(client)
	ctx := context.Background()
	t.Run("given normal condition when query transactions from nodes should return uniform data", func(t *testing.T) {
		for blockNumber := 0; blockNumber < testingProperties.ToBlockNumber; blockNumber++ {
			err := awaitility.Await(time.Second, 15*time.Second, func() bool {
				var checkResult *entities.Transaction
				for nodeIndex := 0; nodeIndex < testingProperties.NodesNumber; nodeIndex++ {
					var response models.PreferenceResponse
					url := fmt.Sprintf(testingProperties.URL,
						nodeIndex+1, blockNumber)
					res, err := restyClient.R().SetContext(ctx).Get(url)
					assert.Nil(t, err)
					err = json.Unmarshal(res.Body(), &response)
					assert.Nil(t, err)
					if response.Data.Major == -1 {
						return false
					}
					if checkResult != nil {
						if checkResult.ID != response.Data.ID {
							log.Debugf("Diff data at block %+v for node %+v between %+v and %+v", blockNumber, nodeIndex, response.Data, checkResult)
							return false
						}
					}
					checkResult = &response.Data
				}
				log.Infof("Passed blockNumber %v with data %v", blockNumber, checkResult)
				return true
			})
			assert.Nil(t, err)
			if err != nil {
				t.FailNow()
			}
		}
	})
}
