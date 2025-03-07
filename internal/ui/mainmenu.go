package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"github.com/calmdev/algorand-rewards/internal/algo"
)

// MainMenu returns the main menu.
func MainMenu(w fyne.Window) *fyne.MainMenu {
	return fyne.NewMainMenu(
		fyne.NewMenu("Account",
			fyne.NewMenuItem("Settings", func() {
				w.SetContent(RenderSettingsView())
				w.Show()
			}),
			fyne.NewMenuItem("Export Rewards", func() {
				d := dialog.NewFileSave(
					func(writer fyne.URIWriteCloser, err error) {
						if err != nil {
							return
						}
						if writer == nil {
							return
						}
						defer writer.Close()
						algo.ExportRewards(writer)
					},
					w,
				)
				d.SetFileName("rewards.csv")
				d.SetFilter(storage.NewExtensionFileFilter([]string{".csv"}))
				d.SetView(dialog.ListView)
				d.Resize(fyne.NewSize(MainWindowWidth-100, MainWindowHeight-100))
				d.Show()
			}),
		),
	)
}
