package tests

import (
	"fmt"
	"github.com/ecodia/golang-awaitility/awaitility"
	"github.com/jarcoal/httpmock"
	"github.com/sonntuet1997/avalanche-simplified/worker/entities"
	"github.com/sonntuet1997/avalanche-simplified/worker/tests/mock_data"
	"github.com/stretchr/testify/assert"
	golibcrontestsuite "gitlab.com/golibs-starter/golib-cron/testsuite"
	"testing"
	"time"
)

func setup(t *testing.T) {
	httpmock.ActivateNonDefault(httpClient)
	httpmock.RegisterResponder(
		"GET",
		"http://truthful-node:8000/v1/node/prefer-transactions/1",
		httpmock.NewStringResponder(200, fmt.Sprintf(mock_data.PreferTransactionTemplate, "1-1", 1, 1)),
	)
	httpmock.RegisterResponder(
		"GET",
		"http://adversary-node:8000/v1/node/prefer-transactions/1",
		httpmock.NewStringResponder(200, fmt.Sprintf(mock_data.PreferTransactionTemplate, "1-2", 1, 2)),
	)
}

func TestFetchNewTransactionJob(t *testing.T) {
	setup(t)
	t.Run("given normal condition when run fetch new introduction job should return truthful-node data", func(t *testing.T) {
		p2pService.NeighborNodes = map[string]*entities.Node{
			"1": {
				Address: "truthful-node",
			},
			"2": {
				Address: "truthful-node",
			},
			"3": {
				Address: "truthful-node",
			},
			"4": {
				Address: "adversary-node",
			},
		}
		consensusService.CurrentPreferenceTransaction = &entities.Transaction{
			ID:    "1-2",
			Major: 1,
			Minor: 2,
		}
		err := awaitility.Await(time.Second, 5*time.Second, func() bool {
			for {
				golibcrontestsuite.RunJob("FetchTransactionJob")
				if consensusService.CurrentBlockNumber.Load() == int64(1) {
					return true
				}
			}
		})
		assert.Nil(t, err)
		assert.Equal(t, *consensusService.ConfirmedTransactions[1], entities.Transaction{
			ID:    "1-1",
			Major: 1,
			Minor: 1,
		})
	})
}
