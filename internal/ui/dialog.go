package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func showTelemetryDialog(w fyne.Window) {
	DialogOpen = true
	dialog.ShowInformation("Telemetry", "Please set your GUID in the settings to enable telemetry.", w)
	DialogOpen = false
}
