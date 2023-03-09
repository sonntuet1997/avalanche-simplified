package tests

import (
	"context"
	"github.com/ecodia/golang-awaitility/awaitility"
	"github.com/sonntuet1997/avalanche-simplyfied/worker/tests/mocks"
	"github.com/stretchr/testify/assert"
	golibcrontestsuite "gitlab.com/golibs-starter/golib-cron/testsuite"
	"gitlab.com/golibs-starter/golib/log"
	"testing"
	"time"
)

func TestSelfIntroductionJob(t *testing.T) {
	err := P2pService.Close()
	assert.Nil(t, err)
	defer P2pService.ListenForBroadcasts(context.Background())
	broadcastListener := mocks.NewBroadcastListener(5555)
	err = awaitility.Await(time.Second, 5*time.Second, func() bool {
		err = broadcastListener.ListenForBroadcasts(context.Background())
		if err != nil {
			log.Debugf("failed with error: %w", err)
		}
		return err == nil
	})
	assert.Nil(t, err)
	defer broadcastListener.Close()
	t.Run("given normal condition when run self introduction job should success", func(t *testing.T) {
		golibcrontestsuite.RunJob("SelfIntroductionJob")
		assert.Nil(t, err)
		select {
		case output := <-broadcastListener.InputChannel:
			log.Debugf("%+v", output)
			break
		case <-time.After(5 * time.Second):
			assert.Fail(t, "Timeout")
			break
		}
	})
}
