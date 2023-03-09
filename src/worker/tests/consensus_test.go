package tests

import (
	"context"
	"github.com/deliveryhero/pipeline/v2"
	"gitlab.com/golibs-starter/golib/log"
	"testing"
)

func TestConsensus(t *testing.T) {
	t.Run("given normal condition when run consensus should success", func(t *testing.T) {
		ctx := context.Background()
		processor := pipeline.NewProcessor(
			func(ctx context.Context, input int) (interface{}, error) {
				log.Debugf("[processor] %+v", input)
				return input, nil
			}, func(i int, err error) {
				log.Errorf("[ConsensusService] failed to process with error: %w", err)
			})
		processOutputChan := pipeline.ProcessConcurrently(
			ctx,
			20,
			processor,
			pipeline.Emit([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}...),
		)
		_ = processOutputChan
		for output := range processOutputChan {
			log.Debugf("[processOutputChan] %+v", output)
		}
		log.Debugf("[processOutputChan] done")
	})
}
