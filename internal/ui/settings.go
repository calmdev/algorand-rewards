package ui

import (
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/calmdev/algorand-rewards/internal/algo"
	"github.com/calmdev/algorand-rewards/internal/app"
)

// SettingsForm returns the settings form.
func SettingsForm(a *app.App) fyne.CanvasObject {
	// createLabel creates a new label with the given text.
	createLabel := func(text string) *widget.Label {
		label := widget.NewLabel(text)
		label.Alignment = fyne.TextAlignLeading
		label.TextStyle = fyne.TextStyle{Bold: false}
		return label
	}

	// createEntry creates a new entry with the given placeholder and initial text.
	createEntry := func(placeholder, initialText string) *widget.Entry {
		entry := widget.NewEntry()
		entry.SetPlaceHolder(placeholder)
		if initialText != "" {
			entry.SetText(initialText)
		}
		return entry
	}

	// Algorand Wallet Address setting
	algorandWalletAddress := createEntry("Enter your Algorand wallet address", a.Address())

	// GUID setting
	guid := createEntry("Enter your GUID", a.GUID())

	// Progress indicator
	progressLabel := canvas.NewText("", theme.Color(theme.ColorNameForeground))
	progressLabel.TextSize = 12

	// Save button
	var saveButton *widget.Button
	saveButton = widget.NewButton("Save", func() {
		// Disable the save button to prevent multiple clicks
		saveButton.Disable()

		// Save the preferences
		a.SetAddress(algorandWalletAddress.Text)
		a.SetGUID(guid.Text)

		// Clear the cache
		_ = a.ClearCacheFile(
			algo.TransactionCacheFile,
			algo.RewardsCacheFile,
		)

		// Show progress indicator
		progressLabel.Text = "Rebuilding cache..."

		// Fetch rewards and transactions concurrently
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			algo.FetchRewards(a.Address())
		}()

		go func() {
			defer wg.Done()
			algo.FetchTransactions(a.Address())
		}()

		go func() {
			wg.Wait()
			progressLabel.Text = "Cache rebuild complete."
			time.Sleep(500 * time.Millisecond)
			progressLabel.Text = ""
			// Close the settings window
			if a.Address() != "" {
				RenderView(&RewardsView{})
			} else {
				RenderView(&SettingsView{})
			}
		}()
	})

	// Form
	form := container.NewVBox(
		createLabel("Algorand Wallet Address:"),
		algorandWalletAddress,
		createLabel("Telemetry GUID:"),
		guid,
	)

	l := newAppLayout()
	l.mainContent = form
	l.bottomBar = container.NewHBox(
		progressLabel,
		layout.NewSpacer(),
		saveButton,
	)

	return container.New(layout.NewCustomPaddedLayout(0, 5, 5, 5), l.render())
}
