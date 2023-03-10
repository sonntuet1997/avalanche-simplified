package tests

import (
	"github.com/sonntuet1997/avalanche-simplified/worker/entities"
	golibtest "gitlab.com/golibs-starter/golib-test"
	"net/http"
	"testing"
)

func TestGetConfirmedTransaction_ShouldReturnSuccess(t *testing.T) {
	url := `/v1/node/transactions`
	transaction0 := &entities.Transaction{ID: "0-0", Major: 0, Minor: 0}
	transaction1 := &entities.Transaction{ID: "1-1", Major: 1, Minor: 1}
	transaction2 := &entities.Transaction{ID: "2-2", Major: 2, Minor: 2}
	consensusService.ConfirmedTransactions = make([]*entities.Transaction, 0)
	consensusService.ConfirmedTransactions = append(consensusService.ConfirmedTransactions, transaction0, transaction1, transaction2)
	t.Run("given normal condition when query list confirmed transactions should return valid data", func(t *testing.T) {
		golibtest.NewRestAssured(t).
			When().
			Get(url).Then().Status(http.StatusOK).
			Body("meta.code", 200).
			Body("data.#", 3).
			Body("data.0.id", "0-0").
			Body("data.0.major", 0).
			Body("data.0.minor", 0).
			Body("data.1.id", "1-1").
			Body("data.1.major", 1).
			Body("data.1.minor", 1).
			Body("data.2.id", "2-2").
			Body("data.2.major", 2).
			Body("data.2.minor", 2)
	})
}
