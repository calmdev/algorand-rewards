package ui

import (
	_ "embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

//go:embed assets/algo-white-icon.png
var algoWhiteIcon []byte
var algoWhiteIconResource = &fyne.StaticResource{
	StaticName:    "algo-white-icon.png",
	StaticContent: algoWhiteIcon,
}

//go:embed assets/algo-white-wordmark.png
var algoWhiteWordmark []byte
var algoWhiteWordmarkResource = &fyne.StaticResource{
	StaticName:    "algo-white-wordmark.png",
	StaticContent: algoWhiteWordmark,
}

//go:embed assets/algo-black-icon.png
var algoBlackIcon []byte
var AlgoBlackIconResource = &fyne.StaticResource{
	StaticName:    "algo-black-icon.png",
	StaticContent: algoBlackIcon,
}

//go:embed assets/algo-black-wordmark.png
var algoBlackWordmark []byte
var algoBlackWordmarkResource = &fyne.StaticResource{
	StaticName:    "algo-black-wordmark.png",
	StaticContent: algoBlackWordmark,
}

// LogoIcon returns the white Algorand icon.
func LogoIcon(size float32) fyne.CanvasObject {
	if fyne.CurrentApp().Settings().ThemeVariant() == theme.VariantLight {
		logo := canvas.NewImageFromResource(AlgoBlackIconResource)
		logo.FillMode = canvas.ImageFillContain
		logo.SetMinSize(fyne.NewSize(size, size))

		return logo
	}

	logo := canvas.NewImageFromResource(algoWhiteIconResource)
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(size, size))

	return logo
}

// LogoWordmark returns the white Algorand wordmark.
func LogoWordmark(size float32) fyne.CanvasObject {
	// if theme is light return the black logo
	if fyne.CurrentApp().Settings().ThemeVariant() == theme.VariantLight {
		logo := canvas.NewImageFromResource(algoBlackWordmarkResource)
		logo.FillMode = canvas.ImageFillContain
		logo.SetMinSize(fyne.NewSize(size, size/2))

		return logo
	}

	logo := canvas.NewImageFromResource(algoWhiteWordmarkResource)
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(size, size/2))

	return logo
}
