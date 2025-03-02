package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
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

	// Stats
	wins := canvas.NewText("Total Wins: "+algo.FormatInt(totalWins), White)
	wins.TextStyle.Bold = true

	rewards := canvas.NewText("Total Rewards: ", White)
	rewards.TextStyle.Bold = true
	rewardsTotal := canvas.NewText(algo.FormatFloat(totalPayout), White)
	rewardsTotal.TextStyle.Bold = true

	spacer := layout.NewSpacer()

	stats := container.NewHBox(
		spacer,
		wins,
		spacer,
		rewards,
		LogoWhiteIcon(10),
		rewardsTotal,
		spacer,
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
			text := canvas.NewText("", White)

			logo := LogoWhiteIcon(10)
			logo.Hide()

			hbox := container.NewHBox(logo, text)

			return container.New(layout.NewCustomPaddedLayout(0, 0, 10, 10), hbox)
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			c := o.(*fyne.Container)

			hbox := c.Objects[0].(*fyne.Container)

			text := hbox.Objects[1].(*canvas.Text)
			text.Text = data[i.Row][i.Col]
			text.Color = White
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
