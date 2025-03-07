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

	// Algorand Wallet Address Label
	algorandWalletAddressLabel := widget.NewLabel("Algorand Wallet Address:")
	algorandWalletAddressLabel.Alignment = fyne.TextAlignLeading
	algorandWalletAddressLabel.TextStyle = fyne.TextStyle{Bold: false}

	// Algorand Wallet Address setting
	algorandWalletAddress := widget.NewEntry()
	algorandWalletAddress.SetPlaceHolder("Enter your Algorand wallet address")
	if a.Preferences().String("Address") != "" {
		algorandWalletAddress.SetText(a.Preferences().String("Address"))
	}

	// GUID Label
	guidLabel := widget.NewLabel("Telemetry GUID:")
	guidLabel.Alignment = fyne.TextAlignLeading
	guidLabel.TextStyle = fyne.TextStyle{Bold: false}

	// GUID setting
	guid := widget.NewEntry()
	guid.SetPlaceHolder("Enter your GUID")
	if a.Preferences().String("GUID") != "" {
		guid.SetText(a.Preferences().String("GUID"))
	}

	// Save button
	saveButton := widget.NewButton("Save", func() {
		// Save the preferences
		algo.Address = algorandWalletAddress.Text
		a.Preferences().SetString("Address", algo.Address)
		a.Preferences().SetString("GUID", guid.Text)
		// Clear the cache
		algo.ClearCache()
		// Close the settings window
		a.Driver().AllWindows()[0].SetContent(RenderMainView())
	})

	// Form
	form := container.NewVBox(
		algorandWalletAddressLabel,
		algorandWalletAddress,
		guidLabel,
		guid,
		saveButton,
	)

	return container.New(layout.NewCustomPaddedLayout(100, 100, 100, 100), form)
}
