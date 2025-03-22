package widget

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type TappableRectangle struct {
	widget.BaseWidget
	rect              *canvas.Rectangle
	tapped            func()
	hoverColor        color.Color
	normalColor       color.Color
	selectedColor     color.Color
	selectedRectangle **TappableRectangle
}

func NewTappableRectangle(c color.Color, tapped func(), selectedRectangle **TappableRectangle) *TappableRectangle {
	var r *TappableRectangle

	r = &TappableRectangle{
		rect: canvas.NewRectangle(c),
		tapped: func() {
			if *selectedRectangle != nil {
				(*selectedRectangle).rect.FillColor = (*selectedRectangle).normalColor
				(*selectedRectangle).rect.Refresh()
			}
			*selectedRectangle = r
			tapped()
		},
		normalColor:       c,
		selectedRectangle: selectedRectangle,
	}
	r.ExtendBaseWidget(r)
	return r
}

func (r *TappableRectangle) Cursor() desktop.Cursor {
	if r.hoverColor != nil {
		return desktop.PointerCursor
	}
	return desktop.DefaultCursor
}

// set the stroke color and width
func (r *TappableRectangle) SetStroke(color color.Color) {
	r.rect.StrokeColor = color
}

func (r *TappableRectangle) SetStrokeWidth(width float32) {
	r.rect.StrokeWidth = width
}

func (r *TappableRectangle) SetCornerRadius(radius float32) {
	r.rect.CornerRadius = radius
}

func (r *TappableRectangle) SetFillColor(c color.Color) {
	r.rect.FillColor = c
	r.normalColor = c
}

func (r *TappableRectangle) SetHoverColor(c color.Color) {
	r.hoverColor = c
}

func (r *TappableRectangle) SetSelectedColor(c color.Color) {
	r.selectedColor = c
}

func (r *TappableRectangle) CreateRenderer() fyne.WidgetRenderer {
	return &rectangleRenderer{rect: r.rect}
}

func (r *TappableRectangle) Tapped(_ *fyne.PointEvent) {
	if r.tapped != nil {
		r.tapped()
	}
	if r.selectedColor != nil {
		r.rect.FillColor = r.selectedColor
		r.rect.Refresh()
	}
}

func (r *TappableRectangle) TappedSecondary(_ *fyne.PointEvent) {}

func (r *TappableRectangle) MouseIn(_ *desktop.MouseEvent) {
	if r.hoverColor == nil || r == *r.selectedRectangle {
		return
	}
	r.rect.FillColor = r.hoverColor
	r.rect.Refresh()
}

func (r *TappableRectangle) MouseOut() {
	if r.hoverColor == nil || r == *r.selectedRectangle {
		return
	}
	r.rect.FillColor = r.normalColor
	r.rect.Refresh()
}

func (r *TappableRectangle) MouseMoved(_ *desktop.MouseEvent) {}

func (r *TappableRectangle) Select() {
	*r.selectedRectangle = r
	r.rect.FillColor = r.selectedColor
	r.rect.Refresh()
}

type rectangleRenderer struct {
	rect *canvas.Rectangle
}

func (r *rectangleRenderer) Layout(size fyne.Size) {
	r.rect.Resize(size)
}

func (r *rectangleRenderer) MinSize() fyne.Size {
	return r.rect.MinSize()
}

func (r *rectangleRenderer) Refresh() {
	r.rect.Refresh()
}

func (r *rectangleRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (r *rectangleRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.rect}
}

func (r *rectangleRenderer) Destroy() {}
