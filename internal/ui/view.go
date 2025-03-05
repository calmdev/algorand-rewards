package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/calmdev/algorand-rewards/internal/algo"
)

var (
	// MainWindowWidth is the main window width.
	MainWindowWidth float32 = 800
	// MainWindowHeight is the main window height.
	MainWindowHeight float32 = 400
)

// RenderMainView returns the main view.
func RenderMainView() fyne.CanvasObject {
	payouts := algo.Payouts()
	account := algo.FetchAccount()

	return container.NewBorder(
		Header(account),
		RewardsPanel(payouts),
		nil,
		nil,
		RewardsTable(payouts),
	)
}

// RenderSettingsView returns the settings view.
func RenderSettingsView() fyne.CanvasObject {
	account := algo.FetchAccount()

	return container.NewBorder(
		Header(account),
		nil,
		nil,
		nil,
		SettingsForm(),
	)
}
