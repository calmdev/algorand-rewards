package algo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	// Address is the Algorand address to fetch rewards for.
	Address = ""
	// printer is a message printer for formatting numbers.
	printer = message.NewPrinter(message.MatchLanguage(language.Make(os.Getenv("LANG")).String()))
)

// ShortAddress returns a shortened version of the address.
func ShortAddress() string {
	if len(Address) < 8 {
		return Address
	}
	return Address[:4] + "..." + Address[len(Address)-4:]
}

// Blocks represents a list of block headers.
type Blocks struct {
	Blocks    []BlockHeader `json:"blocks"`
	NextToken string        `json:"next-token"`
}

// BlockHeader represents a block header.
type BlockHeader struct {
	Timestamp      int64 `json:"timestamp"`
	ProposerPayout int64 `json:"proposer-payout"`
	Bonus          int64 `json:"bonus"`
	FeesCollected  int64 `json:"fees-collected"`
}

// PayoutAlgos returns the total payout in microalgos.
func (b *BlockHeader) PayoutAlgos() int64 {
	return b.ProposerPayout
}

// Time returns the timestamp as a time.Time.
func (b *BlockHeader) Time() time.Time {
	return time.Unix(b.Timestamp, 0)
}

// PayoutDate represents a payout for a specific date.
type PayoutDate struct {
	Date          string `json:"date"`
	Payout        int64  `json:"payout"`
	Bonus         int64  `json:"bonus"`
	FeesCollected int64  `json:"fees-collected"`
	TotalWins     int64  `json:"totalWins"`
	BestDay       bool   `json:"bestDay"`
}

// FractionalPayout returns the payout in fractional algos.
func (pd *PayoutDate) FractionalPayout() float64 {
	return float64(pd.Payout) / 1e6
}

// FractionalBonus returns the bonus in fractional algos.
func (pd *PayoutDate) FractionalBonus() float64 {
	return float64(pd.Bonus) / 1e6
}

// FractionalFeesCollected returns the fees collected in fractional algos.
func (pd *PayoutDate) FractionalFeesCollected() float64 {
	return float64(pd.FeesCollected) / 1e6 / 2
}

// FormatFloat formats a float64 as a string.
func FormatFloat(f float64) string {
	return printer.Sprintf("%.6f", f)
}

// FormatInt formats an int64 as a string.
func FormatInt(i int64) string {
	return printer.Sprintf("%d", i)
}

// fetchBlockHeadersRecursive returns a list of block headers recursively.
func fetchBlockHeadersRecursive(nextToken string, blocks []BlockHeader, afterTime time.Time) []BlockHeader {
	client := &http.Client{}
	url := "https://mainnet-idx.4160.nodely.dev/v2/block-headers?proposers=" + Address
	if nextToken != "" {
		url += "&next=" + nextToken
	}
	if !afterTime.IsZero() {
		url += "&after-time=" + afterTime.Format(time.RFC3339)
		fmt.Println("After time:", afterTime.Format(time.RFC3339))
	}
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("accept", "application/json")
	res, _ := client.Do(req)
	body, _ := io.ReadAll(res.Body)

	// unmarshal json
	var blockHeaders Blocks
	_ = json.Unmarshal([]byte(body), &blockHeaders)

	blocks = append(blocks, blockHeaders.Blocks...)

	fmt.Println("Current token:", nextToken)

	if blockHeaders.NextToken != "" {
		fmt.Println("Next token:", blockHeaders.NextToken)
		return fetchBlockHeadersRecursive(blockHeaders.NextToken, blocks, afterTime)
	}

	fmt.Printf("Fetched %d blocks\n", len(blocks))

	return blocks
}

// Payouts returns a list of payouts for the current address.
func Payouts() []PayoutDate {
	cacheFile, err := storage.Child(fyne.CurrentApp().Storage().RootURI(), RewardsCacheFile)
	if err != nil {
		return nil
	}

	exists, err := storage.Exists(cacheFile)
	if err != nil {
		return nil
	}

	var blocks []BlockHeader
	if exists {
		// Read the cache file
		file, err := os.OpenFile(cacheFile.Path(), os.O_RDONLY, 0644)
		if err != nil {
			return nil
		}
		defer file.Close()
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&blocks)
		if err != nil {
			return nil
		}
		fmt.Printf("Read %d blocks from cache\n", len(blocks))
	}

	// Determine the latest timestamp from the cached blocks
	var latestTime time.Time
	if len(blocks) > 0 {
		latestTime = blocks[len(blocks)-1].Time()
	}

	// Fetch new blocks recursively starting from the latest timestamp
	newBlocks := fetchBlockHeadersRecursive("", []BlockHeader{}, latestTime)
	blocks = append(blocks, newBlocks...)

	// Write the updated blocks back to the cache file
	file, err := os.OpenFile(cacheFile.Path(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	err = encoder.Encode(blocks)
	if err != nil {
		return nil
	}
	fmt.Printf("Wrote %d blocks to cache\n", len(blocks))

	// Create a map of payouts by date
	payoutsByDate := make(map[string]PayoutDate)
	for _, block := range blocks {
		date := block.Time().Format("2006-01-02")
		if _, ok := payoutsByDate[date]; !ok {
			payoutsByDate[date] = PayoutDate{
				Date:          date,
				Payout:        block.PayoutAlgos(),
				TotalWins:     1,
				Bonus:         block.Bonus,
				FeesCollected: block.FeesCollected,
			}
		} else {
			payout := payoutsByDate[date].Payout + block.PayoutAlgos()
			bonus := payoutsByDate[date].Bonus + block.Bonus
			feesCollected := payoutsByDate[date].FeesCollected + block.FeesCollected
			totalWins := payoutsByDate[date].TotalWins + 1
			payoutsByDate[date] = PayoutDate{
				Date:          date,
				Payout:        payout,
				TotalWins:     totalWins,
				Bonus:         bonus,
				FeesCollected: feesCollected,
			}
		}
	}

	// Fill any missing dates
	var startDate time.Time
	var endDate time.Time
	for _, block := range blocks {
		if startDate.IsZero() || block.Time().Before(startDate) {
			startDate = block.Time()
		}
		if endDate.IsZero() || block.Time().After(endDate) {
			endDate = block.Time()
		}
	}
	for d := startDate; d.Before(endDate); d = d.AddDate(0, 0, 1) {
		date := d.Format("2006-01-02")
		if _, ok := payoutsByDate[date]; !ok {
			payoutsByDate[date] = PayoutDate{Date: date, Payout: 0, TotalWins: 0}
		}
	}

	// Append today's date if not in the map
	today := time.Now().Format("2006-01-02")
	if _, ok := payoutsByDate[today]; !ok {
		payoutsByDate[today] = PayoutDate{Date: today, Payout: 0, TotalWins: 0}
	}

	// Create a slice of PayoutDate
	var payouts []PayoutDate
	for _, payout := range payoutsByDate {
		payouts = append(payouts, payout)
	}
	// Find the max payout
	var maxPayout int64
	for _, payout := range payouts {
		if payout.Payout > maxPayout {
			maxPayout = payout.Payout
		}
	}
	// Set the best day flag
	for i := range payouts {
		if payouts[i].Payout == maxPayout {
			payouts[i].BestDay = true
		}
	}
	// Create a new slice sorted by payout date in descending order
	for i := range payouts {
		for j := i + 1; j < len(payouts); j++ {
			if payouts[i].Date < payouts[j].Date {
				payouts[i], payouts[j] = payouts[j], payouts[i]
			}
		}
	}

	return payouts
}

// RewardsCacheFile is the name of the cache file for rewards.
var RewardsCacheFile = "rewards.json"

// ClearRewardsCache clears the rewards cache.
func ClearRewardsCache() {
	// Delete the cache file
	cacheFile, _ := storage.Child(fyne.CurrentApp().Storage().RootURI(), RewardsCacheFile)
	_ = storage.Delete(cacheFile)
}
