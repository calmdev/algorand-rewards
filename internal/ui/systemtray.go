package ui

import (
	"fyne.io/fyne/v2"
)

// SystemTray returns the system tray menu.
func SystemTray(w fyne.Window) *fyne.Menu {
	return fyne.NewMenu("Algorand Rewards",
		fyne.NewMenuItem("Show", func() {
			w.Show()
		}),
		fyne.NewMenuItem("Hide", func() {
			w.Hide()
		}),
		fyne.NewMenuItem("Refresh", func() {
			w.SetContent(RenderMainView())
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Quit", func() {
			fyne.CurrentApp().Quit()
		}),
	)
}
