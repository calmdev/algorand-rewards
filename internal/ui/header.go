package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/calmdev/algorand-rewards/internal/algo"
	iw "github.com/calmdev/algorand-rewards/internal/ui/widget"
)

// Header returns the header of the application.
func Header() fyne.CanvasObject {
	header := container.NewHBox(
		LogoWhiteWordMark(70),
		layout.NewSpacer(),
		AccountStatusIcon(),
		canvas.NewText(algo.ShortAddress(), White),
	)

	return container.New(layout.NewCustomPaddedLayout(0, 0, 10, 10), header)
}

// AccountStatusIcon returns the account status icon.
func AccountStatusIcon() fyne.Widget {
	activity := iw.NewActivity()
	activity.Start()

	return activity
}
