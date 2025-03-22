package algo

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"github.com/calmdev/algorand-rewards/internal/app"
	"github.com/calmdev/algorand-rewards/internal/format"
	"github.com/calmdev/algorand-rewards/internal/nodely"
)

// TransactionCacheFile is the rewards cache file.
var TransactionCacheFile = "transactions.json"

type PaymentTransaction struct {
	Amount      int64  `json:"amount"`
	CloseAmount int64  `json:"close-amount"`
	Receiver    string `json:"receiver"`
}

type TransactionDetail struct {
	ID             string              `json:"id"`
	Timestamp      int64               `json:"round-time"`
	Type           string              `json:"tx-type"`
	Sender         string              `json:"sender"`
	Payment        *PaymentTransaction `json:"payment-transaction"`
	ConfirmedRound int64               `json:"confirmed-round"`
	Fee            int64               `json:"fee"`
}

// Time returns the timestamp as a time.Time.
func (t *TransactionDetail) Time() time.Time {
	return time.Unix(t.Timestamp, 0)
}

// AlgoFee returns the fee in Algos.
func (t *TransactionDetail) AlgoFee() float64 {
	return float64(t.Fee) / 1e6
}

// TypeString returns the type of the transaction as a string.
func (t *TransactionDetail) TypeString() string {
	switch t.Type {
	case "pay":
		return "Payment"
	case "hb":
		return "Heartbeat"
	case "keyreg":
		return "Key Registration"
	case "acfg":
		return "Asset Configuration"
	case "axfer":
		return "Asset Transfer"
	case "afrz":
		return "Asset Freeze"
	case "appl":
		return "Application Call"
	case "stpf":
		return "State Proof"
	default:
		return t.Type
	}
}

type TransactionList struct {
	Transactions       []TransactionDetail            `json:"transactions"`
	TransactionsByDate map[string][]TransactionDetail `json:"transactions-by-date"`
	NextToken          string                         `json:"next-token"`
}

// Data returns the data for the transaction list.
func (tl *TransactionList) Data() [][]string {
	var data = [][]string{
		{"Time", "Type", "Sender", "Receiver", "Algo", "Fee", "ConfirmedRound", "ID"},
	}

	// Append transactions to data
	for _, t := range tl.Transactions {
		var receiver string
		var amount int64
		if t.Payment != nil {
			receiver = t.Payment.Receiver
			amount = t.Payment.Amount
		}

		data = append(data, []string{
			t.Time().Format(time.RFC3339),
			t.TypeString(),
			format.AddressShort(t.Sender),
			format.AddressShort(receiver),
			format.Float(float64(amount) / 1e6),
			format.Float(t.AlgoFee()),
			fmt.Sprintf("%d", t.ConfirmedRound),
			t.ID,
		})
	}

	return data
}

// fetchTransactionsRecursive returns a list of block headers recursively.
//
// Docs: https://nodely.io/swagger/index.html?url=/swagger/api/4160/indexer.oas3.yml#/lookup/lookupAccountTransactions
func fetchTransactionsRecursive(client *nodely.ClientIndexer, address, nextToken string, transactions []TransactionDetail, afterTime time.Time) []TransactionDetail {
	url := "/v2/accounts/" + address + "/transactions"
	queryParams := make([]string, 0)
	if nextToken != "" {
		queryParams = append(queryParams, "next="+nextToken)
	}
	if !afterTime.IsZero() {
		queryParams = append(queryParams, "after-time="+afterTime.Format(time.RFC3339))
		fmt.Println("After time:", afterTime.Format(time.RFC3339))
	}
	if len(queryParams) > 0 {
		url += "?" + strings.Join(queryParams, "&")
	}

	var txs TransactionList
	err := client.Get(url, &txs)
	if err != nil {
		return transactions
	}

	transactions = append(transactions, txs.Transactions...)

	fmt.Println("Current token:", nextToken)

	if txs.NextToken != "" {
		fmt.Println("Next token:", txs.NextToken)
		return fetchTransactionsRecursive(client, address, txs.NextToken, transactions, afterTime)
	}

	fmt.Printf("Fetched %d transactions\n", len(transactions))

	return transactions
}

// FetchTransactions returns a list of transactions for the current address.
func FetchTransactions(address string) *TransactionList {
	client := nodely.NewClientIndexer()

	cacheFile, err := app.CurrentApp().CacheFile(TransactionCacheFile)
	if err != nil {
		return nil
	}

	exists, err := storage.Exists(cacheFile)
	if err != nil {
		return nil
	}

	var txs []TransactionDetail
	if exists {
		// Read the cache file
		file, err := os.OpenFile(cacheFile.Path(), os.O_RDONLY, 0644)
		if err != nil {
			return nil
		}
		defer file.Close()
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&txs)
		if err != nil {
			return nil
		}
		fmt.Printf("Read %d transactions from cache\n", len(txs))
	}

	// Determine the latest timestamp from the cached blocks
	var latestTime time.Time
	if len(txs) > 0 {
		latestTime = txs[0].Time()
	}

	// Fetch new transactions recursively starting from the latest timestamp
	newTxs := fetchTransactionsRecursive(client, address, "", []TransactionDetail{}, latestTime)
	txs = append(newTxs, txs...)

	// Write the updated blocks back to the cache file
	file, err := os.OpenFile(cacheFile.Path(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	err = encoder.Encode(txs)
	if err != nil {
		return nil
	}
	fmt.Printf("Wrote %d transactions to cache\n", len(txs))

	// Transaction by date
	transactionsByDate := make(map[string][]TransactionDetail)
	for _, tx := range txs {
		date := tx.Time().Format("2006-01-02")
		transactionsByDate[date] = append(transactionsByDate[date], tx)
	}

	return &TransactionList{Transactions: txs, TransactionsByDate: transactionsByDate}
}

// ExportTransactions exports the transactions to a CSV file.
func ExportTransactions(address string, writeCloser fyne.URIWriteCloser) {
	// Create a new CSV writer
	writer := csv.NewWriter(writeCloser)
	defer writer.Flush()

	data := FetchTransactions(address).Data()

	// Write the CSV header
	err := writer.Write(data[0])
	if err != nil {
		return
	}

	// Write the CSV rows
	for _, payout := range data[1:] {
		err = writer.Write(payout)
		if err != nil {
			return
		}
	}
}
