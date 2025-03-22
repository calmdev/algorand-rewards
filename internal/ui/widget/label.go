package widget

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type ColorLabel struct {
	widget.Label
	color     color.Color
	minWidth  float32
	textStyle fyne.TextStyle
	alignment fyne.TextAlign
}

func NewColorLabel(text string, color color.Color) *ColorLabel {
	label := &ColorLabel{}
	label.ExtendBaseWidget(label)
	label.SetText(text)
	label.SetColor(color)
	return label
}

func (l *ColorLabel) SetColor(c color.Color) {
	l.color = c
	l.Refresh()
}

func (l *ColorLabel) SetMinWidth(width float32) {
	l.minWidth = width
	l.Refresh()
}

func (l *ColorLabel) SetTextStyle(style fyne.TextStyle) {
	l.textStyle = style
	l.Refresh()
}

func (l *ColorLabel) SetTextAlign(align fyne.TextAlign) {
	l.alignment = align
	l.Refresh()
}

func (l *ColorLabel) CreateRenderer() fyne.WidgetRenderer {
	text := canvas.NewText(l.Text, l.color)
	text.TextStyle = l.textStyle
	text.Alignment = l.alignment
	return &colorLabelRenderer{text: text, label: l}
}

type colorLabelRenderer struct {
	text  *canvas.Text
	label *ColorLabel
}

func (r *colorLabelRenderer) Layout(size fyne.Size) {
	r.text.Resize(size)
}

func (r *colorLabelRenderer) MinSize() fyne.Size {
	minSize := r.text.MinSize()
	if r.label.minWidth > minSize.Width {
		minSize.Width = r.label.minWidth
	}
	return fyne.NewSize(minSize.Width+10, minSize.Height+10)
}

func (r *colorLabelRenderer) Refresh() {
	r.text.Text = r.label.Text
	r.text.Color = r.label.color
	r.text.TextStyle = r.label.textStyle
	r.text.Alignment = r.label.alignment
	r.text.Refresh()
}

func (r *colorLabelRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (r *colorLabelRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.text}
}

func (r *colorLabelRenderer) Destroy() {}
