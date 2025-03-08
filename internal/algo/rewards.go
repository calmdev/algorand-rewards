package algo

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

// Rewards represents a list of payouts.
type Rewards struct {
	Payouts     []PayoutDate
	TotalPayout float64
	TotalWins   int64
	MinPayout   float64
	MaxPayout   float64
}

// NewRewards creates a new Rewards instance.
func NewRewards(payouts []PayoutDate) *Rewards {
	rewards := Rewards{
		Payouts: payouts,
	}
	rewards.PayoutsByPreference()
	rewards.TotalPayout = TotalPayout(&rewards)
	rewards.TotalWins = TotalWins(&rewards)
	rewards.MinPayout = MinPayout(&rewards)
	rewards.MaxPayout = MaxPayout(&rewards)
	return &rewards
}

// Data returns the data for the rewards.
func (r *Rewards) Data() [][]string {
	var data = [][]string{{"Date", "Wins", "Fees Collected", "Bonus", "Rewards"}}

	// Append payouts to data
	for _, payout := range r.Payouts {
		data = append(data, []string{
			payout.Date,
			FormatInt(payout.TotalWins),
			FormatFloat(payout.FractionalFeesCollected()),
			FormatFloat(payout.FractionalBonus()),
			FormatFloat(payout.FractionalPayout()),
		})
	}

	return data
}

// ByDay sorts the payouts by day.
func (r *Rewards) ByDay() {
	// Create a new slice to hold the sorted payouts
	sortedPayouts := make([]PayoutDate, len(r.Payouts))
	copy(sortedPayouts, r.Payouts)

	// Sort the payouts by date
	for i := range sortedPayouts {
		for j := i + 1; j < len(sortedPayouts); j++ {
			if sortedPayouts[i].Date < sortedPayouts[j].Date {
				sortedPayouts[i], sortedPayouts[j] = sortedPayouts[j], sortedPayouts[i]
			}
		}
	}

	r.Payouts = sortedPayouts
}

// ByMonth returns the payouts by month.
func (r *Rewards) ByMonth() {
	// Group payouts by month
	monthlyPayouts := make(map[string][]PayoutDate)
	for _, payout := range r.Payouts {
		month := payout.Date[:7] // Assuming the date format is YYYY-MM-DD
		monthlyPayouts[month] = append(monthlyPayouts[month], payout)
	}

	var data []PayoutDate
	// Aggregate data by month
	for month, payouts := range monthlyPayouts {
		var totalWins int64
		var totalFees, totalBonus, totalRewards float64

		for _, payout := range payouts {
			totalWins += payout.TotalWins
			totalFees += payout.FractionalFeesCollected()
			totalBonus += payout.FractionalBonus()
			totalRewards += payout.FractionalPayout()
		}

		data = append(data, PayoutDate{
			Date:          month,
			Payout:        int64(totalRewards * 1e6),
			TotalWins:     totalWins,
			Bonus:         int64(totalBonus * 1e6),
			FeesCollected: int64(totalFees * 1e6 * 2),
		})
	}

	// Sort the payouts by date
	for i := range data {
		for j := i + 1; j < len(data); j++ {
			if data[i].Date < data[j].Date {
				data[i], data[j] = data[j], data[i]
			}
		}
	}

	r.Payouts = data
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
func Payouts() *Rewards {
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

	return NewRewards(payouts)
}

// RewardsCacheFile is the rewards cache file.
var RewardsCacheFile = "rewards.json"

// ClearRewardsCache clears the rewards cache.
func ClearRewardsCache() {
	// Delete the rewards cache file
	rewardsCacheFile, _ := storage.Child(fyne.CurrentApp().Storage().RootURI(), RewardsCacheFile)
	_ = storage.Delete(rewardsCacheFile)
}

// ExportRewards exports the rewards to a CSV file.
func ExportRewards(writeCloser fyne.URIWriteCloser) {
	// Create a new CSV writer
	writer := csv.NewWriter(writeCloser)
	defer writer.Flush()

	data := Payouts().Data()

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

// TotalPayout returns the total payout.
func TotalPayout(r *Rewards) float64 {
	var total float64
	for _, payout := range r.Payouts {
		total += payout.FractionalPayout()
	}
	return total
}

// Total wins returns the total wins.
func TotalWins(r *Rewards) int64 {
	var total int64
	for _, payout := range r.Payouts {
		total += payout.TotalWins
	}
	return total
}

// MinPayout returns the minimum payout.
func MinPayout(r *Rewards) float64 {
	var minPayout float64
	for _, payout := range r.Payouts {
		if minPayout == 0 {
			minPayout = payout.FractionalPayout()
			continue
		}
		if payout.FractionalPayout() < minPayout {
			minPayout = payout.FractionalPayout()
		}
	}
	return minPayout
}

// MaxPayout returns the maximum payout.
func MaxPayout(r *Rewards) float64 {
	var maxPayout float64
	for _, payout := range r.Payouts {
		if payout.FractionalPayout() > maxPayout {
			maxPayout = payout.FractionalPayout()
		}
	}
	return maxPayout
}

// PayoutsByPreference returns the payouts by preference.
func (r *Rewards) PayoutsByPreference() {
	// Check if the view is by day or by month
	switch fyne.CurrentApp().Preferences().String("RewardsView") {
	case "day":
		r.ByDay()
	case "month":
		r.ByMonth()
	default:
		r.ByDay()
	}
}
