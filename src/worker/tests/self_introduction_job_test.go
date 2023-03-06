package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelfIntroductionJob(t *testing.T) {
	t.Run("given normal condition when run self introduction job should success", func(t *testing.T) {
		err := P2pService.Close()
		assert.Nil(t, err)
	})
}
