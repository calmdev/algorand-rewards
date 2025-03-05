package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"github.com/calmdev/algorand-rewards/internal/algo"
	iw "github.com/calmdev/algorand-rewards/internal/ui/widget"
)

// Header returns the header of the application.
func Header(account *algo.Account) fyne.CanvasObject {
	balanceTotal := canvas.NewText(algo.FormatFloat(account.FractionalBalance()), theme.Color(theme.ColorNameForeground))
	balanceTotal.TextStyle.Bold = true

	header := container.NewHBox(
		LogoWordmark(70),
		layout.NewSpacer(),
		LogoIcon(10),
		balanceTotal,
		AccountStatusIcon(),
		canvas.NewText(algo.ShortAddress(), theme.Color(theme.ColorNameForeground)),
	)

	return container.New(layout.NewCustomPaddedLayout(0, 0, 10, 10), header)
}

// AccountStatusIcon returns the account status icon.
func AccountStatusIcon() fyne.Widget {
	activity := iw.NewActivity()
	activity.Start()

	return activity
}
