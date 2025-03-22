package ui

import (
	"fyne.io/fyne/v2"
	"github.com/calmdev/algorand-rewards/internal/app"
)

// SystemTray returns the system tray menu.
func SystemTray(a *app.App, w fyne.Window) *fyne.Menu {
	return fyne.NewMenu(app.AppName,
		fyne.NewMenuItem("Show", func() {
			w.Show()
		}),
		fyne.NewMenuItem("Hide", func() {
			w.Hide()
		}),
		fyne.NewMenuItem("Refresh", func() {
			RenderView(&RewardsView{})
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Quit", func() {
			a.Quit()
		}),
	)
}
