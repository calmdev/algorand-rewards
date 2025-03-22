package ui

import (
	"github.com/calmdev/algorand-rewards/internal/algo"
	"github.com/calmdev/algorand-rewards/internal/app"
)

const (
	// Main window dimensions
	MainWindowWidth  float32 = 800
	MainWindowHeight float32 = 400
)

// RenderView renders the given view.
func RenderView(v View) {
	v.Render(app.CurrentApp())
}

// View interface represents a view.
type View interface {
	Render(a *app.App)
}

// RewardsView struct represents the rewards view.
type RewardsView struct{}

// Render renders the rewards view.
func (v *RewardsView) Render(a *app.App) {
	Layout.loading()

	go func() {
		account := algo.FetchAccount(a.Address())
		rewards := algo.FetchRewards(a.Address())

		Layout.updateTopBar(Header(account))
		Layout.updateMainContent(RewardsList(account, rewards))
		Layout.currentView = v
	}()

	Layout.markActiveButton(0)
}

// SettingsView struct represents the settings view.
type SettingsView struct{}

// Render renders the settings view.
func (v *SettingsView) Render(a *app.App) {
	Layout.loading()

	go func() {
		account := algo.FetchAccount(a.Address())

		Layout.updateTopBar(Header(account))
	}()

	Layout.updateMainContent(SettingsForm(a))
	Layout.markActiveButton(2)
	Layout.currentView = v
}

// TransactionsView struct represents the transactions view.
type TransactionsView struct{}

// Render renders the transactions view.
func (v *TransactionsView) Render(a *app.App) {
	Layout.loading()

	go func() {
		account := algo.FetchAccount(a.Address())
		transactions := algo.FetchTransactions(a.Address())

		Layout.updateTopBar(Header(account))
		Layout.updateMainContent(TransactionsList(account, transactions))
		Layout.currentView = v
	}()

	Layout.markActiveButton(1)
}
