package algo

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"github.com/calmdev/algorand-rewards/internal/app"
	"github.com/calmdev/algorand-rewards/internal/format"
	"github.com/calmdev/algorand-rewards/internal/nodely"
)

// RewardsCacheFile is the rewards cache file.
var RewardsCacheFile = "rewards.json"

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
	rewards.SortByView(app.CurrentApp().RewardsView())
	rewards.TotalPayout = TotalPayout(rewards.Payouts)
	rewards.TotalWins = TotalWins(rewards.Payouts)
	rewards.MinPayout = MinPayout(rewards.Payouts)
	rewards.MaxPayout = MaxPayout(rewards.Payouts)
	return &rewards
}

// Data returns the data for the rewards.
func (r *Rewards) Data() [][]string {
	var data = [][]string{{"Date", "Wins", "Fees Collected", "Bonus", "Rewards"}}

	// Append payouts to data
	for _, payout := range r.Payouts {
		data = append(data, []string{
			payout.Date,
			format.Int(payout.TotalWins),
			format.Float(payout.AlgoFeesCollected()),
			format.Float(payout.AlgoBonus()),
			format.Float(payout.AlgoPayout()),
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

// ByDayofWeek sorts the payouts by day of the week.
func (r *Rewards) ByDayofWeek() {
	// Group payouts by day of the week
	weeklyPayouts := make(map[string][]PayoutDate)
	for _, payout := range r.Payouts {
		// Get the day of the week from string YYYY-MM-DD
		date := payout.Date
		dateTime, err := time.Parse("2006-01-02", date)
		if err != nil {
			continue
		}
		dayOfWeek := dateTime.Weekday().String()
		weeklyPayouts[dayOfWeek] = append(weeklyPayouts[dayOfWeek], payout)
	}
	var data []PayoutDate
	// Aggregate data by day of the week
	for day, payouts := range weeklyPayouts {
		var totalWins int64
		var totalFees, totalBonus, totalRewards float64

		for _, payout := range payouts {
			totalWins += payout.TotalWins
			totalFees += payout.AlgoFeesCollected()
			totalBonus += payout.AlgoBonus()
			totalRewards += payout.AlgoPayout()
		}

		data = append(data, PayoutDate{
			Date:          day,
			Payout:        int64(totalRewards * 1e6),
			TotalWins:     totalWins,
			Bonus:         int64(totalBonus * 1e6),
			FeesCollected: int64(totalFees * 1e6 * 2),
		})
	}
	// Sort the payouts by day of the week start at Monday
	daysOfWeek := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	var sortedData []PayoutDate
	for _, day := range daysOfWeek {
		for _, payout := range data {
			if payout.Date == day {
				sortedData = append(sortedData, payout)
				break
			}
		}
	}

	r.Payouts = sortedData
}

// ByWeek sorts the payouts by week.
func (r *Rewards) ByWeek() {
	// Group payouts by week
	weeklyPayouts := make(map[string][]PayoutDate)
	for _, payout := range r.Payouts {
		date := payout.Date
		dateTime, err := time.Parse("2006-01-02", date)
		if err != nil {
			continue
		}
		year, week := dateTime.ISOWeek()
		weekKey := fmt.Sprintf("%d-W%02d", year, week)
		weeklyPayouts[weekKey] = append(weeklyPayouts[weekKey], payout)
	}

	var data []PayoutDate
	// Aggregate data by week
	for week, payouts := range weeklyPayouts {
		var totalWins int64
		var totalFees, totalBonus, totalRewards float64

		for _, payout := range payouts {
			totalWins += payout.TotalWins
			totalFees += payout.AlgoFeesCollected()
			totalBonus += payout.AlgoBonus()
			totalRewards += payout.AlgoPayout()
		}

		data = append(data, PayoutDate{
			Date:          week,
			Payout:        int64(totalRewards * 1e6),
			TotalWins:     totalWins,
			Bonus:         int64(totalBonus * 1e6),
			FeesCollected: int64(totalFees * 1e6 * 2),
		})
	}

	// Sort the payouts by week using sort.Slice
	sort.Slice(data, func(i, j int) bool {
		return data[i].Date > data[j].Date
	})

	r.Payouts = data
}

// ByMonth sorts the payouts by month.
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
			totalFees += payout.AlgoFeesCollected()
			totalBonus += payout.AlgoBonus()
			totalRewards += payout.AlgoPayout()
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

// ByQuarter sorts the payouts by quarter.
func (r *Rewards) ByQuarter() {
	// Group payouts by quarter
	quarterlyPayouts := make(map[string][]PayoutDate)
	for _, payout := range r.Payouts {
		year := payout.Date[:4] // Assuming the date format is YYYY-MM-DD
		month := payout.Date[5:7]
		var quarter string
		switch month {
		case "01", "02", "03":
			quarter = year + "-Q1"
		case "04", "05", "06":
			quarter = year + "-Q2"
		case "07", "08", "09":
			quarter = year + "-Q3"
		case "10", "11", "12":
			quarter = year + "-Q4"
		}
		quarterlyPayouts[quarter] = append(quarterlyPayouts[quarter], payout)
	}

	var data []PayoutDate
	// Aggregate data by quarter
	for quarter, payouts := range quarterlyPayouts {
		var totalWins int64
		var totalFees, totalBonus, totalRewards float64

		for _, payout := range payouts {
			totalWins += payout.TotalWins
			totalFees += payout.AlgoFeesCollected()
			totalBonus += payout.AlgoBonus()
			totalRewards += payout.AlgoPayout()
		}

		data = append(data, PayoutDate{
			Date:          quarter,
			Payout:        int64(totalRewards * 1e6),
			TotalWins:     totalWins,
			Bonus:         int64(totalBonus * 1e6),
			FeesCollected: int64(totalFees * 1e6 * 2),
		})
	}

	r.Payouts = data
}

// ByYear sorts the payouts by year.
func (r *Rewards) ByYear() {
	// Group payouts by year
	yearlyPayouts := make(map[string][]PayoutDate)
	for _, payout := range r.Payouts {
		year := payout.Date[:4] // Assuming the date format is YYYY-MM-DD
		yearlyPayouts[year] = append(yearlyPayouts[year], payout)
	}

	var data []PayoutDate
	// Aggregate data by year
	for year, payouts := range yearlyPayouts {
		var totalWins int64
		var totalFees, totalBonus, totalRewards float64

		for _, payout := range payouts {
			totalWins += payout.TotalWins
			totalFees += payout.AlgoFeesCollected()
			totalBonus += payout.AlgoBonus()
			totalRewards += payout.AlgoPayout()
		}

		data = append(data, PayoutDate{
			Date:          year,
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

// TotalPayout returns the total payout.
func TotalPayout(payouts []PayoutDate) float64 {
	var total float64
	for _, payout := range payouts {
		total += payout.AlgoPayout()
	}
	return total
}

// Total wins returns the total wins.
func TotalWins(payouts []PayoutDate) int64 {
	var total int64
	for _, payout := range payouts {
		total += payout.TotalWins
	}
	return total
}

// MinPayout returns the minimum payout.
func MinPayout(payouts []PayoutDate) float64 {
	var minPayout float64
	for _, payout := range payouts {
		if minPayout == 0 {
			minPayout = payout.AlgoPayout()
			continue
		}
		if payout.AlgoPayout() < minPayout {
			minPayout = payout.AlgoPayout()
		}
	}
	return minPayout
}

// MaxPayout returns the maximum payout.
func MaxPayout(payouts []PayoutDate) float64 {
	var maxPayout float64
	for _, payout := range payouts {
		if payout.AlgoPayout() > maxPayout {
			maxPayout = payout.AlgoPayout()
		}
	}
	return maxPayout
}

// SortByView sorts the rewards by the given view.
func (r *Rewards) SortByView(view string) {
	switch view {
	case "day":
		r.ByDay()
	case "dayOfWeek":
		r.ByDayofWeek()
	case "week":
		r.ByWeek()
	case "month":
		r.ByMonth()
	case "quarter":
		r.ByQuarter()
	case "year":
		r.ByYear()
	default:
		r.ByDay()
	}
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

// AlgoPayout returns the payout in Algos.
func (pd *PayoutDate) AlgoPayout() float64 {
	return float64(pd.Payout) / 1e6
}

// AlgoBonus returns the bonus in Algos.
func (pd *PayoutDate) AlgoBonus() float64 {
	return float64(pd.Bonus) / 1e6
}

// AlgoFeesCollected returns the fees collected in Algos.
func (pd *PayoutDate) AlgoFeesCollected() float64 {
	return float64(pd.FeesCollected) / 1e6 / 2
}

// fetchBlockHeadersRecursive returns a list of block headers recursively.
//
// Docs: https://nodely.io/swagger/index.html?url=/swagger/api/4160/indexer.oas3.yml#/search/searchForBlockHeaders
func fetchBlockHeadersRecursive(client *nodely.ClientIndexer, address, nextToken string, blocks []BlockHeader, afterTime time.Time) []BlockHeader {
	url := "/v2/block-headers?proposers=" + address
	if nextToken != "" {
		url += "&next=" + nextToken
	}
	if !afterTime.IsZero() {
		url += "&after-time=" + afterTime.Format(time.RFC3339)
		fmt.Println("After time:", afterTime.Format(time.RFC3339))
	}

	var blockHeaders Blocks
	err := client.Get(url, &blockHeaders)
	if err != nil {
		return blocks
	}

	blocks = append(blocks, blockHeaders.Blocks...)

	fmt.Println("Current token:", nextToken)

	if blockHeaders.NextToken != "" {
		fmt.Println("Next token:", blockHeaders.NextToken)
		return fetchBlockHeadersRecursive(client, address, blockHeaders.NextToken, blocks, afterTime)
	}

	fmt.Printf("Fetched %d blocks\n", len(blocks))

	return blocks
}

// FetchRewards returns a list of payouts for the current address.
func FetchRewards(address string) *Rewards {
	client := nodely.NewClientIndexer()

	cacheFile, err := app.CurrentApp().CacheFile(RewardsCacheFile)
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
	newBlocks := fetchBlockHeadersRecursive(client, address, "", []BlockHeader{}, latestTime)
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

// ExportRewards exports the rewards to a CSV file.
func ExportRewards(address string, writeCloser fyne.URIWriteCloser) {
	// Create a new CSV writer
	writer := csv.NewWriter(writeCloser)
	defer writer.Flush()

	data := FetchRewards(address).Data()

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
