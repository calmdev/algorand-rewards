package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/calmdev/algorand-rewards/internal/algo"
	"github.com/calmdev/algorand-rewards/internal/ui"
)

func main() {
	// App
	a := app.NewWithID("com.calmdev.algorand-rewards")
	a.SetIcon(ui.AlgoBlackIconResource)
	a.Settings().SetTheme(&ui.CustomTheme{})

	// Preferences
	if a.Preferences().String("Address") != "" {
		algo.Address = a.Preferences().String("Address")
	}

	// Window
	w := a.NewWindow("Algorand Rewards")
	w.Resize(fyne.NewSize(ui.MainWindowWidth, ui.MainWindowHeight))
	w.SetFixedSize(true)
	w.CenterOnScreen()
	w.SetMaster()
	w.SetCloseIntercept(func() {
		w.Hide()
	})

	// Window Content
	if a.Preferences().String("Address") == "" {
		w.SetContent(ui.RenderSettingsView())
	} else {
		w.SetContent(ui.RenderMainView())
	}

	// System Tray
	if d, ok := a.(desktop.App); ok {
		d.SetSystemTrayMenu(ui.SystemTray(w))
	}

	// Menu
	w.SetMainMenu(ui.MainMenu(w))

	// Watch for OS theme variant changes.
	var themeVariant = a.Settings().ThemeVariant()
	go ui.WatchForThemeVariantChanges(a, w, &themeVariant)

	// Run
	w.ShowAndRun()
}
