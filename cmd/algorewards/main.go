package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/calmdev/algorand-rewards/internal/algo"
	"github.com/calmdev/algorand-rewards/internal/app"
	"github.com/calmdev/algorand-rewards/internal/ui"
)

func main() {
	// App
	a := app.NewApp()
	a.SetIcon(ui.AlgoBlackIconResource)
	a.Settings().SetTheme(&ui.AppTheme{})

	// Version Check
	a.VersionCheck(func() {
		_ = a.ClearCacheFile(
			algo.RewardsCacheFile,
			algo.TransactionCacheFile,
		)
	})

	// Window
	w := a.NewWindow(app.AppName)
	if a.IsWindows() {
		w.Resize(fyne.NewSize(ui.MainWindowWidth, ui.MainWindowHeight+30))
	} else {
		w.Resize(fyne.NewSize(ui.MainWindowWidth, ui.MainWindowHeight))
	}
	w.SetFixedSize(true)
	w.CenterOnScreen()
	w.SetMaster()
	w.SetCloseIntercept(func() {
		w.Hide()
	})
	w.SetContent(ui.RenderLayout())

	// Initial View
	if a.Address() == "" {
		ui.RenderView(&ui.SettingsView{})
	} else {
		ui.RenderView(&ui.RewardsView{})
	}

	// System Tray
	if d, ok := a.App.(desktop.App); ok {
		d.SetSystemTrayMenu(ui.SystemTray(a, w))
	}

	// Menu
	w.SetMainMenu(ui.MainMenu(a, w))

	// Watch for OS theme variant changes.
	var themeVariant = a.Settings().ThemeVariant()
	go ui.WatchForThemeVariantChanges(a, w, &themeVariant)

	// Run
	w.ShowAndRun()
}
