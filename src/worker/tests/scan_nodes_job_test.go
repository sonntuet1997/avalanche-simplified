package tests

import (
	"github.com/jarcoal/httpmock"
	"github.com/sonntuet1997/avalanche-simplified/worker/entities"
	"github.com/stretchr/testify/assert"
	golibcrontestsuite "gitlab.com/golibs-starter/golib-cron/testsuite"
	"testing"
)

func TestScanNodeJob(t *testing.T) {
	httpmock.ActivateNonDefault(httpClient)
	httpmock.RegisterResponder(
		"GET",
		"http://node-1:8000/actuator/health",
		httpmock.NewStringResponder(200, "{}"),
	)
	httpmock.RegisterResponder(
		"GET",
		"http://node-2:8000/actuator/health",
		httpmock.NewStringResponder(200, "{}"),
	)
	httpmock.RegisterResponder(
		"GET",
		"http://node-3:8000/actuator/health",
		httpmock.NewStringResponder(200, "{}"),
	)
	httpmock.RegisterResponder(
		"GET",
		"http://node-4:8000/actuator/health",
		httpmock.NewStringResponder(400, "{}"),
	)
	t.Run("given normal condition when run scan job should return update nodes data", func(t *testing.T) {
		golibcrontestsuite.RunJob("ScanNodesJob")
		expectedNodes := map[string]*entities.Node{
			"node-1": {
				Address: "node-1",
			},
			"node-2": {
				Address: "node-2",
			},
			"node-3": {
				Address: "node-3",
			},
		}
		for k, v := range expectedNodes {
			assert.Equal(t, v, p2pService.NeighborNodes[k])
		}
		assert.Nil(t, p2pService.NeighborNodes["node-4"])
	})
}
