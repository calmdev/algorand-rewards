package ui

import (
	"fmt"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"github.com/calmdev/algorand-rewards/internal/algo"
)

// MainMenu returns the main menu.
func MainMenu(w fyne.Window) *fyne.MainMenu {
	return fyne.NewMainMenu(
		fyne.NewMenu("Account",
			&fyne.MenuItem{
				Label: "Refresh",
				Action: func() {
					w.SetContent(RenderMainView())
					w.Show()
				},
			},
			&fyne.MenuItem{
				Label: "Settings",
				Action: func() {
					w.SetContent(RenderSettingsView())
					w.Show()
				},
			},
			&fyne.MenuItem{
				Label: "Export Rewards",
				Action: func() {
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
				},
			},
			fyne.NewMenuItemSeparator(),
			&fyne.MenuItem{
				Label: "Telemetry",
				Action: func() {
					guid := fyne.CurrentApp().Preferences().String("GUID")
					if guid == "" {
						dialog.ShowInformation("Telemetry", "Please set your GUID in the settings to enable telemetry.", w)
						return
					}
					url := url.URL{
						Scheme:   "https",
						Host:     "g.nodely.io",
						Path:     "/d/telemetry/node-telemetry",
						RawQuery: fmt.Sprintf("var-GUID=%s&orgId=1&from=now-24h&to=now", guid),
					}
					if err := fyne.CurrentApp().OpenURL(&url); err != nil {
						dialog.ShowError(err, w)
					}
				},
			},
		),
	)
}
