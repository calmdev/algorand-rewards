package algo

import (
	"encoding/json"
	"io"
	"net/http"
)

// Account represents stats of an account.
type Account struct {
	Amount            int64  `json:"amount"`
	Status            string `json:"status"`
	IncentiveEligible bool   `json:"incentive-eligible"`
	LastHeartbeat     int64  `json:"last-heartbeat"`
	LastProposed      int64  `json:"last-proposed"`
}

// FractionalBalance returns the fractional balance.
func (a *Account) FractionalBalance() float64 {
	return float64(a.Amount) / 1e6
}

// FetchAccount fetches account stats from the nodely api.
//
// Docs: https://nodely.io/swagger/index.html?url=/swagger/api/4160/algod.oas3.yml#/public/AccountInformation
func FetchAccount() *Account {
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, "https://mainnet-api.4160.nodely.dev/v2/accounts/"+Address, nil)
	req.Header.Add("accept", "application/json")
	res, _ := client.Do(req)
	body, _ := io.ReadAll(res.Body)

	var account Account
	_ = json.Unmarshal([]byte(body), &account)

	return &account
}
