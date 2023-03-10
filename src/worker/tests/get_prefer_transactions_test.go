package tests

import (
	"fmt"
	"github.com/sonntuet1997/avalanche-simplified/worker/entities"
	golibtest "gitlab.com/golibs-starter/golib-test"
	"net/http"
	"testing"
)

func TestGetPreferTransaction_ShouldReturnSuccess(t *testing.T) {
	urlTemplate := `/v1/node/prefer-transactions/%d`
	transaction0 := &entities.Transaction{ID: "0-0", Major: 0, Minor: 0}
	transaction1 := &entities.Transaction{ID: "1-1", Major: 1, Minor: 1}
	transaction2 := &entities.Transaction{ID: "2-2", Major: 2, Minor: 2}
	ConsensusService.ConfirmedTransactions = make([]*entities.Transaction, 0)
	ConsensusService.ConfirmedTransactions = append(ConsensusService.ConfirmedTransactions, transaction0, transaction1, transaction2)
	ConsensusService.CurrentBlockNumber.Store(2)
	t.Run("given normal condition when query confirmed transaction should return valid data", func(t *testing.T) {
		url := fmt.Sprintf(urlTemplate, 1)
		golibtest.NewRestAssured(t).
			When().
			Get(url).Then().Status(http.StatusOK).
			Body("meta.code", 200).
			Body("data.id", "1-1").
			Body("data.major", 1).
			Body("data.minor", 1)
	})
	t.Run("given normal condition when query next transaction should return valid data", func(t *testing.T) {
		url := fmt.Sprintf(urlTemplate, 3)
		ConsensusService.CurrentPreferenceTransaction = &entities.Transaction{ID: "3-3", Major: 3, Minor: 3}
		golibtest.NewRestAssured(t).
			When().
			Get(url).Then().Status(http.StatusOK).
			Body("meta.code", 200).
			Body("data.id", "3-3").
			Body("data.major", 3).
			Body("data.minor", 3)
	})
	t.Run("given normal condition when query invalid transaction should return valid data", func(t *testing.T) {
		url := fmt.Sprintf(urlTemplate, 4)
		ConsensusService.CurrentPreferenceTransaction = &entities.Transaction{ID: "3-3", Major: 3, Minor: 3}
		golibtest.NewRestAssured(t).
			When().
			Get(url).Then().Status(http.StatusOK).
			Body("meta.code", 200).
			Body("data.id", "-1").
			Body("data.major", -1).
			Body("data.minor", 0)
	})
}
