package ui

import (
	"fmt"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

// MainMenu returns the main menu.
func MainMenu(w fyne.Window) *fyne.MainMenu {
	return fyne.NewMainMenu(
		accountMenu(w),
		rewardsMenu(w),
	)
}

// accountMenu returns the account menu.
func accountMenu(w fyne.Window) *fyne.Menu {
	var refresh *fyne.MenuItem
	var settings *fyne.MenuItem
	var telemetry *fyne.MenuItem

	refresh = &fyne.MenuItem{
		Label: "Refresh",
		Action: func() {
			w.SetContent(RenderMainView())
			w.Show()
		},
	}

	settings = &fyne.MenuItem{
		Label: "Settings",
		Action: func() {
			w.SetContent(RenderSettingsView())
			w.Show()
		},
	}

	telemetry = &fyne.MenuItem{
		Label: "Telemetry",
		Action: func() {
			guid := fyne.CurrentApp().Preferences().String("GUID")
			if guid == "" {
				showTelemetryDialog(w)
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
	}

	sep := fyne.NewMenuItemSeparator()

	return fyne.NewMenu("Account",
		refresh,
		settings,
		sep,
		telemetry,
	)
}

// rewardsMenu returns the rewards menu.
func rewardsMenu(w fyne.Window) *fyne.Menu {
	var rewardsByDay *fyne.MenuItem
	var rewardsByMonth *fyne.MenuItem
	var exportRewards *fyne.MenuItem

	rewardsViewPref := fyne.CurrentApp().Preferences().String("RewardsView")

	rewardsByDay = &fyne.MenuItem{
		Label:   "By Day",
		Checked: rewardsViewPref == "day" || rewardsViewPref == "",
		Action: func() {
			fyne.CurrentApp().Preferences().SetString("RewardsView", "day")
			rewardsByDay.Checked = true
			rewardsByMonth.Checked = false
			w.SetContent(RenderMainView())
			w.Show()
		},
	}

	rewardsByMonth = &fyne.MenuItem{
		Label:   "By Month",
		Checked: rewardsViewPref == "month",
		Action: func() {
			fyne.CurrentApp().Preferences().SetString("RewardsView", "month")
			rewardsByDay.Checked = false
			rewardsByMonth.Checked = true
			w.SetContent(RenderMainView())
			w.Show()
		},
	}

	exportRewards = &fyne.MenuItem{
		Label: "Export Rewards",
		Action: func() {
			RewardsExportDialog(w)
		},
	}

	sep := fyne.NewMenuItemSeparator()

	return fyne.NewMenu("Rewards", rewardsByDay, rewardsByMonth, sep, exportRewards)
}
