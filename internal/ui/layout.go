package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/calmdev/algorand-rewards/internal/app"
)

// Layout is the application layout.
var Layout *appLayout

// RenderLayout renders the application layout.
func RenderLayout() fyne.CanvasObject {
	Layout = newAppLayout()
	Layout.topBar = Header(nil) // No account yet
	Layout.leftBar = LeftSidebar()

	Layout.container = Layout.render()

	return Layout.container
}

// appLayout represents the layout of the app.
type appLayout struct {
	topBar      fyne.CanvasObject
	leftBar     fyne.CanvasObject
	rightBar    fyne.CanvasObject
	bottomBar   fyne.CanvasObject
	mainContent fyne.CanvasObject

	container   *fyne.Container
	currentView View
}

// newAppLayout returns a new AppLayout.
func newAppLayout() *appLayout {
	return &appLayout{
		topBar:      container.NewHBox(),
		bottomBar:   container.NewHBox(),
		leftBar:     container.NewVBox(),
		rightBar:    container.NewVBox(),
		mainContent: container.NewVBox(),
	}
}

// render renders the AppLayout.
func (l *appLayout) render() *fyne.Container {
	return container.NewBorder(
		l.topBar,
		l.bottomBar,
		l.leftBar,
		l.rightBar,
		l.mainContent,
	)
}

// loading renders a loading message in the main content.
func (l *appLayout) loading() {
	loading := container.NewCenter(widget.NewIcon(AlgoIconResource()))
	l.updateMainContent(loading)
}

// updateMainContent updates the main content.
func (l *appLayout) updateMainContent(content fyne.CanvasObject) {
	l.mainContent = content
	l.container.Objects[0] = l.mainContent
	l.container.Objects[0].Refresh()
}

// updateTopBar updates the top bar.
func (l *appLayout) updateTopBar(content fyne.CanvasObject) {
	l.topBar = content
	topBarContainer := l.container.Objects[1].(*fyne.Container)
	topBarContainer.Objects[0] = l.topBar
	topBarContainer.Refresh()
}

// markActiveButton marks the button at the given index as active.
func (l *appLayout) markActiveButton(index int) {
	leftBar := l.container.Objects[3].(*fyne.Container).Objects[0].(*fyne.Container).Objects

	// Mark all buttons as low importance and enable them
	for _, obj := range leftBar {
		if btn, ok := obj.(*widget.Button); ok {
			btn.Importance = widget.LowImportance
			btn.Enable()
		}
	}

	// Disable all buttons except settings if there is no address configured.
	if app.CurrentApp().Address() == "" {
		for i, obj := range leftBar {
			if btn, ok := obj.(*widget.Button); ok && i != 3 {
				btn.Disable()
			}
		}
	}

	// Mark the selected button as high importance
	if btn, ok := leftBar[index].(*widget.Button); ok {
		btn.Importance = widget.HighImportance
	}

	l.container.Refresh()
}
