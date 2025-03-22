package algo

import (
	"github.com/calmdev/algorand-rewards/internal/nodely"
)

// Account represents stats of an account.
type Account struct {
	Address           string `json:"-"`
	Amount            int64  `json:"amount"`
	Status            string `json:"status"`
	IncentiveEligible bool   `json:"incentive-eligible"`
	LastHeartbeat     int64  `json:"last-heartbeat"`
	LastProposed      int64  `json:"last-proposed"`
}

// AlgoBalance returns the balance in Algos.
func (a *Account) AlgoBalance() float64 {
	return float64(a.Amount) / 1e6
}

// FetchAccount fetches account stats from the nodely api.
//
// Docs: https://nodely.io/swagger/index.html?url=/swagger/api/4160/algod.oas3.yml#/public/AccountInformation
func FetchAccount(address string) *Account {
	client := nodely.NewClient()

	account := Account{Address: address}
	err := client.Get("/v2/accounts/"+account.Address, &account)
	if err != nil {
		return nil
	}

	return &account
}
