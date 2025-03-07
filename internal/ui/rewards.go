package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/calmdev/algorand-rewards/internal/algo"
)

// RewardsPanel returns a panel of reward stats.
func RewardsPanel(payouts []algo.PayoutDate) fyne.CanvasObject {
	var totalWins int64
	for _, payout := range payouts {
		totalWins += payout.TotalWins
	}

	var totalPayout float64
	for _, payout := range payouts {
		totalPayout += payout.FractionalPayout()
	}

	var minPayout float64
	for _, payout := range payouts {
		if minPayout == 0 {
			minPayout = payout.FractionalPayout()
			continue
		}
		if payout.FractionalPayout() < minPayout {
			minPayout = payout.FractionalPayout()
		}
	}

	var maxPayout float64
	for _, payout := range payouts {
		if payout.FractionalPayout() > maxPayout {
			maxPayout = payout.FractionalPayout()
		}
	}

	// Stats
	wins := canvas.NewText("Wins: "+algo.FormatInt(totalWins), theme.Color(theme.ColorNameForeground))
	wins.TextStyle.Bold = true
	wins.TextSize = 12

	minRewards := canvas.NewText("Min: ", theme.Color(theme.ColorNameForeground))
	minRewards.TextStyle.Bold = true
	minRewards.TextSize = 12
	minRewardsTotal := canvas.NewText(algo.FormatFloatShort(minPayout), theme.Color(theme.ColorNameForeground))
	minRewardsTotal.TextStyle.Bold = true
	minRewardsTotal.TextSize = 12

	maxRewards := canvas.NewText("Max: ", theme.Color(theme.ColorNameForeground))
	maxRewards.TextStyle.Bold = true
	maxRewards.TextSize = 12
	maxRewardsTotal := canvas.NewText(algo.FormatFloatShort(maxPayout), theme.Color(theme.ColorNameForeground))
	maxRewardsTotal.TextStyle.Bold = true
	maxRewardsTotal.TextSize = 12

	rewards := canvas.NewText("Rewards: ", theme.Color(theme.ColorNameForeground))
	rewards.TextStyle.Bold = true
	rewards.TextSize = 12
	rewardsTotal := canvas.NewText(algo.FormatFloatShort(totalPayout), theme.Color(theme.ColorNameForeground))
	rewardsTotal.TextStyle.Bold = true
	rewardsTotal.TextSize = 12

	spacer := layout.NewSpacer()

	stats := container.NewHBox(
		wins,
		spacer,
		minRewards,
		LogoIcon(10),
		minRewardsTotal,
		spacer,
		maxRewards,
		LogoIcon(10),
		maxRewardsTotal,
		spacer,
		rewards,
		LogoIcon(10),
		rewardsTotal,
	)

	return container.New(layout.NewCustomPaddedLayout(12, 12, 10, 10), stats)
}

// RewardsTable returns a table of rewards.
func RewardsTable(payouts []algo.PayoutDate) fyne.CanvasObject {
	var data = [][]string{{"Date", "Wins", "Fees Collected", "Bonus", "Rewards"}}

	// Append payouts to data
	for _, payout := range payouts {
		data = append(data, []string{
			payout.Date,
			algo.FormatInt(payout.TotalWins),
			algo.FormatFloat(payout.FractionalFeesCollected()),
			algo.FormatFloat(payout.FractionalBonus()),
			algo.FormatFloat(payout.FractionalPayout()),
		})
	}

	// Table
	table := widget.NewTable(
		func() (int, int) {
			return len(data), len(data[0])
		},
		func() fyne.CanvasObject {
			text := canvas.NewText("", theme.Color(theme.ColorNameForeground))

			logo := LogoIcon(10)
			logo.Hide()

			hbox := container.NewHBox(logo, text)

			return container.New(layout.NewCustomPaddedLayout(0, 0, 10, 10), hbox)
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			c := o.(*fyne.Container)

			hbox := c.Objects[0].(*fyne.Container)

			text := hbox.Objects[1].(*canvas.Text)
			text.Text = data[i.Row][i.Col]
			text.Color = theme.Color(theme.ColorNameForeground)
			text.TextStyle = fyne.TextStyle{Bold: false, Italic: false}

			logo := hbox.Objects[0].(*canvas.Image)

			// Row 0 is the header
			if i.Row == 0 {
				text.TextStyle = fyne.TextStyle{Bold: true}
			} else {
				// date column is gray
				if i.Col == 0 {
					text.Color = Grey
				}

				// rewards column
				if i.Col == 4 {
					logo.Show()
				} else {
					logo.Hide()
				}
			}
		})

	table.StickyRowCount = 1 // Sticky header
	table.HideSeparators = false

	// Set column widths
	for i := range data[0] {
		table.SetColumnWidth(i, 150)
	}

	// Set row heights
	for i := range data {
		table.SetRowHeight(i, 30)
	}

	return container.New(layout.NewCustomPaddedLayout(0, 0, 10, 10), table)
}
