package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/calmdev/algorand-rewards/internal/algo"
)

// SettingsForm returns the settings form.
func SettingsForm() fyne.CanvasObject {
	// App
	a := fyne.CurrentApp()

	// Input Label
	label := widget.NewLabel("Algorand Wallet Address:")
	label.Alignment = fyne.TextAlignLeading
	label.TextStyle = fyne.TextStyle{Bold: false}

	// Algorand Wallet Address settings
	algorandWalletAddress := widget.NewEntry()
	algorandWalletAddress.SetPlaceHolder("Enter your Algorand wallet address")
	if a.Preferences().String("Address") != "" {
		algorandWalletAddress.SetText(a.Preferences().String("Address"))
	}
	algorandWalletAddress.OnChanged = func(s string) {
		// Update the Algorand wallet address
	}

	// Save button
	saveButton := widget.NewButton("Save", func() {
		// Save the Algorand wallet address
		algo.Address = algorandWalletAddress.Text
		a.Preferences().SetString("Address", algo.Address)
		// Clear the cache
		algo.ClearRewardsCache()
		// Close the settings window
		a.Driver().AllWindows()[0].SetContent(RenderMainView())
	})

	// Form
	form := container.NewVBox(
		label,
		algorandWalletAddress,
		saveButton,
	)

	return container.New(layout.NewCustomPaddedLayout(100, 100, 100, 100), form)
}
