package ui

import (
	"fmt"
	"image/color"
	"net/url"

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

// TransactionsPanel returns a panel of reward stats.
func TransactionsPanel(tx *algo.TransactionDetail) fyne.CanvasObject {
	// createStat creates a new state for the transaction panel.
	createStat := func(labelText, valueText string, valueURL *url.URL, icon fyne.CanvasObject) *fyne.Container {
		label := canvas.NewText(labelText, theme.Color(theme.ColorNameForeground))
		label.TextStyle.Bold = true
		label.TextSize = 12

		var value fyne.CanvasObject
		if valueURL != nil {
			value = iw.NewHyperlink(valueText, valueURL)
			value.(*iw.Hyperlink).TextSize = 12
			value.(*iw.Hyperlink).TextStyle.Bold = true
			value.(*iw.Hyperlink).Color = theme.Color(theme.ColorNameHyperlink)
		} else {
			value = canvas.NewText(valueText, theme.Color(theme.ColorNameForeground))
			value.(*canvas.Text).TextSize = 12
			value.(*canvas.Text).TextStyle.Bold = true
		}

		components := []fyne.CanvasObject{label, layout.NewSpacer()}
		if icon != nil {
			components = append(components, icon)
		}
		components = append(components, value)

		return container.NewHBox(components...)
	}

	txId := createStat("ID: ", format.AddressShort(tx.ID), &url.URL{
		Scheme: "https",
		Host:   "allo.info",
		Path:   fmt.Sprintf("/tx/%s", tx.ID),
	}, nil)
	txType := createStat("Type: ", tx.TypeString(), nil, nil)
	txFee := createStat("Fee: ", format.FloatShort(tx.AlgoFee()), nil, AlgoIcon(10))
	txBlock := createStat("Block: ", fmt.Sprintf("%d", tx.ConfirmedRound), &url.URL{
		Scheme: "https",
		Host:   "allo.info",
		Path:   fmt.Sprintf("/block/%d", tx.ConfirmedRound),
	}, nil)

	stats := container.NewHBox(txId, layout.NewSpacer(), txType, layout.NewSpacer(), txFee, layout.NewSpacer(), txBlock)

	return container.New(layout.NewCustomPaddedLayout(5, 5, 5, 5), stats)
}

const batchSize = 25

// TransactionsList returns a list of transactions.
func TransactionsList(account *algo.Account, t *algo.TransactionList) fyne.CanvasObject {
	var l *appLayout
	var data = t.Transactions
	var offset int
	var content *fyne.Container

	var selected *iw.TappableRectangle

	// createHeader creates a new header for the transaction list.
	createHeader := func(date string) *fyne.Container {
		bg := iw.NewTappableRectangle(theme.Color(theme.ColorNameHeaderBackground), func() {
			fmt.Println("Tapped on date:", date)
		}, &selected)
		bg.SetStroke(theme.Color(theme.ColorNameSeparator))
		bg.SetStrokeWidth(1)
		bg.SetCornerRadius(5)

		title := canvas.NewText(date, theme.Color(theme.ColorNameForeground))
		title.TextStyle.Bold = true

		transactions := canvas.NewText(fmt.Sprintf("%s Transactions", format.Int(int64(len(t.TransactionsByDate[date])))), theme.Color(theme.ColorNameForeground))
		transactions.TextSize = 12
		transactions.Color = Grey

		return container.NewStack(
			bg,
			container.New(layout.NewCustomPaddedLayout(5, 5, 5, 5), container.NewHBox(
				title,
				layout.NewSpacer(),
				transactions,
			)),
		)
	}

	// createTransactionItem creates a new transaction item.
	createTransactionItem := func(row algo.TransactionDetail, l *appLayout, selected **iw.TappableRectangle) *fyne.Container {
		rec := iw.NewTappableRectangle(color.Transparent, func() {
			fmt.Println("Tapped on transaction:", row.ID)
			l.bottomBar = TransactionsPanel(&row)
			l.container.Objects[2].(*fyne.Container).RemoveAll()
			for _, obj := range l.bottomBar.(*fyne.Container).Objects {
				l.container.Objects[2].(*fyne.Container).Add(obj)
			}
			l.container.Objects[2].(*fyne.Container).Refresh()
		}, selected)
		rec.SetHoverColor(theme.Color(theme.ColorNameHover))
		rec.SetSelectedColor(theme.Color(theme.ColorNameSelection))

		if *selected == nil {
			rec.Select()
		}

		var amount int64
		if row.Payment != nil {
			amount = row.Payment.Amount
		}

		var amountLabel *iw.ColorLabel
		if amount > 0 {
			if row.Sender == account.Address {
				amountLabel = iw.NewColorLabel("-"+format.Float(float64(amount)/1e6), DarkRed)
			} else {
				amountLabel = iw.NewColorLabel("+"+format.Float(float64(amount)/1e6), DarkGreen)
			}
		} else {
			amountLabel = iw.NewColorLabel(format.Float(float64(amount)/1e6), theme.Color(theme.ColorNameForeground))
		}

		item := container.New(layout.NewCustomPaddedLayout(0, 0, 5, 5), container.NewStack(
			rec,
			container.New(layout.NewCustomPaddedLayout(0, 0, 5, 5),
				container.NewHBox(
					iw.NewColorLabel(row.Time().Format("03:04:05 PM"), Grey),
					iw.NewColorLabel(format.AddressShort(row.Sender), theme.Color(theme.ColorNameForeground)),
					layout.NewSpacer(),
					amountLabel,
					TransactionTypeIcon(row.Type, 12),
				),
			),
		))
		return item
	}

	// loadMore loads more transactions.
	loadMore := func() {
		end := min(offset+batchSize, len(data))
		var lastDate string
		for _, row := range data[offset:end] {
			if date := row.Time().Format("2006-01-02"); date != lastDate {
				content.Add(createHeader(date))
				lastDate = date
			} else {
				content.Add(container.New(layout.NewCustomPaddedLayout(0, 0, 5, 5), widget.NewSeparator()))
			}
			item := createTransactionItem(row, l, &selected)
			content.Add(item)
		}
		offset = end
	}

	content = container.NewVBox()
	scroll := container.NewVScroll(container.New(layout.NewCustomPaddedLayout(0, 0, 5, 5), content))
	scroll.OnScrolled = func(p fyne.Position) {
		// Load more transactions when near the bottom.
		if p.Y > scroll.Content.Size().Height-scroll.Size().Height-100 {
			loadMore()
		}
	}

	l = newAppLayout()
	l.mainContent = scroll
	l.bottomBar = TransactionsPanel(&data[0])
	l.container = l.render()

	loadMore() // Initial load

	return l.container
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// TransactionsExportDialog opens a dialog to export rewards.
func TransactionsExportDialog(a *app.App, w fyne.Window) {
	d := dialog.NewFileSave(
		func(writer fyne.URIWriteCloser, err error) {
			if err != nil {
				return
			}
			if writer == nil {
				return
			}
			defer writer.Close()
			algo.ExportTransactions(a.Address(), writer)
		},
		w,
	)
	d.SetFileName("transactions.csv")
	d.SetFilter(storage.NewExtensionFileFilter([]string{".csv"}))
	d.SetView(dialog.ListView)
	d.Resize(fyne.NewSize(MainWindowWidth-20, MainWindowHeight-20))
	d.Show()
}
