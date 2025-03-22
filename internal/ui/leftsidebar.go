package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func LeftSidebar() fyne.CanvasObject {
	var rewards *widget.Button
	var settings *widget.Button
	var transactions *widget.Button

	var iconContainer *fyne.Container

	rewards = &widget.Button{
		Importance: widget.LowImportance,
		Icon:       theme.HomeIcon(),
		OnTapped: func() {
			RenderView(&RewardsView{})
			iconContainer.Refresh()
		},
	}

	settings = &widget.Button{
		Importance: widget.LowImportance,
		Icon:       theme.SettingsIcon(),
		OnTapped: func() {
			RenderView(&SettingsView{})
			iconContainer.Refresh()
		},
	}

	transactions = &widget.Button{
		Importance: widget.LowImportance,
		Icon:       theme.HistoryIcon(),
		OnTapped: func() {
			RenderView(&TransactionsView{})
			iconContainer.Refresh()
		},
	}

	iconContainer = container.NewVBox(
		rewards,
		transactions,
		settings,
	)

	return container.New(layout.NewCustomPaddedLayout(0, 5, 5, 0), iconContainer)
}
