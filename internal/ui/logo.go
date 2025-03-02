package ui

import (
	_ "embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
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

// LogoWhiteIcon returns the white Algorand icon.
func LogoWhiteIcon(size float32) fyne.CanvasObject {
	logo := canvas.NewImageFromResource(algoWhiteIconResource)
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(size, size))

	return logo
}

// LogoWhiteWordMark returns the white Algorand wordmark.
func LogoWhiteWordMark(size float32) fyne.CanvasObject {
	logo := canvas.NewImageFromResource(algoWhiteWordmarkResource)
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(size, size/2))

	return logo
}
