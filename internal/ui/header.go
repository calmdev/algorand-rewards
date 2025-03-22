package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"github.com/calmdev/algorand-rewards/internal/algo"
	"github.com/calmdev/algorand-rewards/internal/format"
)

// Header returns the header of the application.
func Header(account *algo.Account) fyne.CanvasObject {
	header := []fyne.CanvasObject{
		AlgoWordmark(70),
		layout.NewSpacer(),
	}

	if account != nil {
		balanceTotal := canvas.NewText(format.Float(account.AlgoBalance()), theme.Color(theme.ColorNameForeground))
		balanceTotal.TextStyle.Bold = true

		header = append(
			header,
			AlgoIcon(10),
			balanceTotal,
			AccountStatusIcon(account),
			canvas.NewText(format.AddressShort(account.Address), theme.Color(theme.ColorNameForeground)),
		)
	}

	return container.New(layout.NewCustomPaddedLayout(0, 0, 5, 5), container.NewHBox(
		header...,
	))
}
