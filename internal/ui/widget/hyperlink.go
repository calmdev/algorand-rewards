package widget

import (
	"image/color"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type Hyperlink struct {
	widget.Hyperlink
	TextSize  float32
	Color     color.Color
	TextStyle fyne.TextStyle
}

func NewHyperlink(text string, url *url.URL) *Hyperlink {
	hl := &Hyperlink{}
	hl.ExtendBaseWidget(hl)
	hl.SetText(text)
	hl.SetURL(url)
	return hl
}

func (h *Hyperlink) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

func (h *Hyperlink) SetTextSize(size float32) {
	h.TextSize = size
	h.Refresh()
}

func (h *Hyperlink) SetColor(color color.Color) {
	h.Color = color
	h.Refresh()
}

func (h *Hyperlink) SetTextStyle(style fyne.TextStyle) {
	h.TextStyle = style
	h.Refresh()
}

func (h *Hyperlink) CreateRenderer() fyne.WidgetRenderer {
	h.ExtendBaseWidget(h)
	text := canvas.NewText(h.Text, h.Color)
	text.TextSize = h.TextSize
	text.TextStyle = h.TextStyle

	line := canvas.NewLine(h.Color)
	line.StrokeWidth = .5

	objects := []fyne.CanvasObject{text, line}
	return &hyperlinkRenderer{objects: objects, text: text, line: line, hyperlink: h}
}

type hyperlinkRenderer struct {
	objects   []fyne.CanvasObject
	text      *canvas.Text
	line      *canvas.Line
	hyperlink *Hyperlink
}

func (r *hyperlinkRenderer) Layout(size fyne.Size) {
	r.text.Resize(size)
	r.line.Position1 = fyne.NewPos(0, r.text.Size().Height)
	r.line.Position2 = fyne.NewPos(r.text.Size().Width, r.text.Size().Height)
}

func (r *hyperlinkRenderer) MinSize() fyne.Size {
	return r.text.MinSize()
}

func (r *hyperlinkRenderer) Refresh() {
	r.text.Text = r.hyperlink.Text
	r.text.Color = r.hyperlink.Color
	r.text.TextSize = r.hyperlink.TextSize
	r.text.TextStyle = r.hyperlink.TextStyle
	r.line.StrokeColor = r.hyperlink.Color
	canvas.Refresh(r.hyperlink)
}

func (r *hyperlinkRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (r *hyperlinkRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *hyperlinkRenderer) Destroy() {}
