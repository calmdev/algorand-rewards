package ui

import (
	"fmt"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/calmdev/algorand-rewards/internal/app"
)

// MainMenu returns the main menu.
func MainMenu(a *app.App, w fyne.Window) *fyne.MainMenu {
	return fyne.NewMainMenu(
		accountMenu(a, w),
		rewardsMenu(a, w),
		transactionsMenu(a, w),
	)
}

// accountMenu returns the account menu.
func accountMenu(a *app.App, w fyne.Window) *fyne.Menu {
	var refresh *fyne.MenuItem
	var settings *fyne.MenuItem
	var telemetry *fyne.MenuItem

	refresh = &fyne.MenuItem{
		Label: "Refresh",
		Action: func() {
			RenderView(&RewardsView{})
			w.Show()
		},
	}

	settings = &fyne.MenuItem{
		Label: "Settings",
		Action: func() {
			RenderView(&SettingsView{})
			w.Show()
		},
	}

	telemetry = &fyne.MenuItem{
		Label: "Telemetry",
		Action: func() {
			guid := a.GUID()
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
			if err := a.OpenURL(&url); err != nil {
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
func rewardsMenu(a *app.App, w fyne.Window) *fyne.Menu {
	var rewardsByDay *fyne.MenuItem
	var rewardsByDayOfWeek *fyne.MenuItem
	var rewardsByWeek *fyne.MenuItem
	var rewardsByMonth *fyne.MenuItem
	var rewardsByQuarter *fyne.MenuItem
	var rewardsByYear *fyne.MenuItem
	var exportRewards *fyne.MenuItem

	rewardsViewPref := a.RewardsView()
	toggleChecked := func(view string) {
		rewardsByDay.Checked = false
		rewardsByDayOfWeek.Checked = false
		rewardsByWeek.Checked = false
		rewardsByMonth.Checked = false
		rewardsByQuarter.Checked = false
		rewardsByYear.Checked = false
		switch view {
		case "day":
			rewardsByDay.Checked = true
			a.SetRewardsView(view)
		case "dayOfWeek":
			rewardsByDayOfWeek.Checked = true
			a.SetRewardsView(view)
		case "week":
			rewardsByWeek.Checked = true
			a.SetRewardsView(view)
		case "month":
			rewardsByMonth.Checked = true
			a.SetRewardsView(view)
		case "quarter":
			rewardsByQuarter.Checked = true
			a.SetRewardsView(view)
		case "year":
			rewardsByYear.Checked = true
			a.SetRewardsView(view)
		}
	}

	rewardsByDay = &fyne.MenuItem{
		Label:   "By Day",
		Checked: rewardsViewPref == "day" || rewardsViewPref == "",
		Action: func() {
			toggleChecked("day")
			RenderView(&RewardsView{})
			w.Show()
		},
	}

	rewardsByDayOfWeek = &fyne.MenuItem{
		Label:   "By Day of Week",
		Checked: rewardsViewPref == "dayOfWeek",
		Action: func() {
			toggleChecked("dayOfWeek")
			RenderView(&RewardsView{})
			w.Show()
		},
	}

	rewardsByWeek = &fyne.MenuItem{
		Label:   "By Week",
		Checked: rewardsViewPref == "week",
		Action: func() {
			toggleChecked("week")
			RenderView(&RewardsView{})
			w.Show()
		},
	}

	rewardsByMonth = &fyne.MenuItem{
		Label:   "By Month",
		Checked: rewardsViewPref == "month",
		Action: func() {
			toggleChecked("month")
			RenderView(&RewardsView{})
			w.Show()
		},
	}

	rewardsByQuarter = &fyne.MenuItem{
		Label:   "By Quarter",
		Checked: false,
		Action: func() {
			toggleChecked("quarter")
			RenderView(&RewardsView{})
			w.Show()
		},
	}

	rewardsByYear = &fyne.MenuItem{
		Label:   "By Year",
		Checked: false,
		Action: func() {
			toggleChecked("year")
			RenderView(&RewardsView{})
			w.Show()
		},
	}

	exportRewards = &fyne.MenuItem{
		Label: "Export Rewards",
		Action: func() {
			RewardsExportDialog(a, w)
		},
	}

	sep := fyne.NewMenuItemSeparator()

	return fyne.NewMenu(
		"Rewards",
		rewardsByDay,
		rewardsByDayOfWeek,
		rewardsByWeek,
		rewardsByMonth,
		rewardsByQuarter,
		rewardsByYear,
		sep,
		exportRewards,
	)
}

// transactionsMenu returns the transactions menu.
func transactionsMenu(a *app.App, w fyne.Window) *fyne.Menu {
	history := &fyne.MenuItem{
		Label: "History",
		Action: func() {
			RenderView(&TransactionsView{})
			w.Show()
		},
	}
	exportTransactions := &fyne.MenuItem{
		Label: "Export Transactions",
		Action: func() {
			TransactionsExportDialog(a, w)
		},
	}

	sep := fyne.NewMenuItemSeparator()

	return fyne.NewMenu(
		"Transactions",
		history,
		sep,
		exportTransactions,
	)
}
