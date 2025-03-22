package ui

import (
	"fmt"
	"image/color"
	"net/url"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/calmdev/algorand-rewards/internal/algo"
	"github.com/calmdev/algorand-rewards/internal/app"
	"github.com/calmdev/algorand-rewards/internal/format"
	iw "github.com/calmdev/algorand-rewards/internal/ui/widget"
)

// RewardsPanel returns a panel of reward stats.
func RewardsPanel(account *algo.Account, r *algo.Rewards) fyne.CanvasObject {
	// createText creates a new text or hyperlink for the rewards panel.
	createText := func(label, value string, bold bool, valueURL *url.URL, icon fyne.CanvasObject) *fyne.Container {
		text := canvas.NewText(label, theme.Color(theme.ColorNameForeground))
		text.TextStyle.Bold = bold
		text.TextSize = 12

		var valueObj fyne.CanvasObject
		if valueURL != nil {
			valueObj = iw.NewHyperlink(value, valueURL)
			valueObj.(*iw.Hyperlink).TextSize = 12
			valueObj.(*iw.Hyperlink).TextStyle.Bold = true
			valueObj.(*iw.Hyperlink).Color = theme.Color(theme.ColorNameHyperlink)
		} else {
			valueObj = canvas.NewText(value, theme.Color(theme.ColorNameForeground))
			valueObj.(*canvas.Text).TextSize = 12
			valueObj.(*canvas.Text).TextStyle.Bold = true
		}

		components := []fyne.CanvasObject{text, layout.NewSpacer()}
		if icon != nil {
			components = append(components, icon)
		}
		components = append(components, valueObj)

		return container.NewHBox(components...)
	}

	wins := createText("Wins: ", format.Int(r.TotalWins), true, nil, nil)
	minRewards := createText("Min: ", format.FloatShort(r.MinPayout), true, nil, AlgoIcon(10))
	maxRewards := createText("Max: ", format.FloatShort(r.MaxPayout), true, nil, AlgoIcon(10))
	rewards := createText("Rewards: ", format.FloatShort(r.TotalPayout), true, &url.URL{
		Scheme: "https",
		Host:   "algonoderewards.com",
		Path:   fmt.Sprintf("/%s", account.Address),
	}, AlgoIcon(10))

	spacer := layout.NewSpacer()

	stats := container.NewHBox(
		wins,
		spacer,
		minRewards,
		spacer,
		maxRewards,
		spacer,
		rewards,
	)

	return container.New(layout.NewCustomPaddedLayout(5, 5, 5, 5), stats)
}

// RewardsList returns a list of rewards.
func RewardsList(account *algo.Account, r *algo.Rewards) fyne.CanvasObject {
	var l *appLayout
	var data = r.Payouts
	var offset int
	var content *fyne.Container

	var selected *iw.TappableRectangle

	// createHeaderLabel creates a new header label.
	createHeaderLabel := func(text string, c color.Color, width float32) *iw.ColorLabel {
		label := iw.NewColorLabel(text, c)
		label.SetMinWidth(width)
		label.SetTextStyle(fyne.TextStyle{Bold: true})
		return label
	}

	// createCellLabel creates a new cell label.
	createCellLabel := func(text string, c color.Color, width float32) *iw.ColorLabel {
		label := iw.NewColorLabel(text, c)
		label.SetMinWidth(width)
		return label
	}

	// createRewardItem creates a new reward item.
	createRewardItem := func(row algo.PayoutDate, l *appLayout, r *algo.Rewards, selected **iw.TappableRectangle) *fyne.Container {
		rec := iw.NewTappableRectangle(color.Transparent, func() {
			fmt.Println("Tapped on reward:", row.Date)
			l.bottomBar = RewardsPanel(account, r)
			l.container.Objects[2].(*fyne.Container).RemoveAll()
			for _, obj := range l.bottomBar.(*fyne.Container).Objects {
				l.container.Objects[2].(*fyne.Container).Add(obj)
			}
			l.container.Objects[2].(*fyne.Container).Refresh()
		}, selected)
		rec.SetHoverColor(theme.Color(theme.ColorNameHover))
		rec.SetSelectedColor(theme.Color(theme.ColorNameSelection))

		if *selected == nil {
			switch app.CurrentApp().RewardsView() {
			case "dayOfWeek":
				if row.Date == time.Now().Weekday().String() {
					rec.Select()
				}
			default:
				rec.Select()
			}
		}

		item := container.NewStack(
			rec,
			container.New(layout.NewCustomPaddedLayout(0, 0, 5, 5),
				container.NewHBox(
					createCellLabel(row.Date, Grey, 120),
					createCellLabel(fmt.Sprintf("%d", row.TotalWins), theme.Color(theme.ColorNameForeground), 80),
					createCellLabel(format.Float(row.AlgoFeesCollected()), theme.Color(theme.ColorNameForeground), 140),
					createCellLabel(format.Float(row.AlgoBonus()), theme.Color(theme.ColorNameForeground), 140),
					AlgoIcon(10),
					createCellLabel(format.Float(row.AlgoPayout()), theme.Color(theme.ColorNameForeground), 140),
				),
			),
		)
		return item
	}

	// loadMore loads more rewards.
	loadMore := func() {
		end := min(offset+batchSize, len(data))
		for _, row := range data[offset:end] {
			item := createRewardItem(row, l, r, &selected)
			content.Add(item)
			content.Add(widget.NewSeparator())
		}
		offset = end
	}

	// Add sticky header
	header := container.NewHBox(
		createHeaderLabel("Date", theme.Color(theme.ColorNameForeground), 120),
		createHeaderLabel("Wins", theme.Color(theme.ColorNameForeground), 80),
		createHeaderLabel("Fees Collected", theme.Color(theme.ColorNameForeground), 140),
		createHeaderLabel("Bonus", theme.Color(theme.ColorNameForeground), 140),
		createHeaderLabel("Rewards", theme.Color(theme.ColorNameForeground), 140),
	)
	headerContainer := container.New(layout.NewCustomPaddedLayout(0, 0, 10, 0), header)

	content = container.NewVBox()
	scroll := container.NewVScroll(container.New(layout.NewCustomPaddedLayout(0, 0, 5, 5), content))
	scroll.SetMinSize(fyne.NewSize(0, 270))

	scroll.OnScrolled = func(p fyne.Position) {
		// Load more rewards when near the bottom.
		if p.Y > scroll.Content.Size().Height-scroll.Size().Height-100 {
			loadMore()
		}
	}

	l = newAppLayout()
	l.mainContent = container.NewVBox(headerContainer, scroll)
	l.bottomBar = RewardsPanel(account, r)
	l.container = l.render()

	loadMore() // Initial load

	return l.container
}

// RewardsExportDialog opens a dialog to export rewards.
func RewardsExportDialog(a *app.App, w fyne.Window) {
	d := dialog.NewFileSave(
		func(writer fyne.URIWriteCloser, err error) {
			if err != nil {
				return
			}
			if writer == nil {
				return
			}
			defer writer.Close()
			algo.ExportRewards(a.Address(), writer)
		},
		w,
	)
	d.SetFileName("rewards.csv")
	d.SetFilter(storage.NewExtensionFileFilter([]string{".csv"}))
	d.SetView(dialog.ListView)
	d.Resize(fyne.NewSize(MainWindowWidth-20, MainWindowHeight-20))
	d.Show()
}
