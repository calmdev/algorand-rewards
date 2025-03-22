package ui

import (
	_ "embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"github.com/calmdev/algorand-rewards/internal/algo"
	"github.com/calmdev/algorand-rewards/internal/app"
	iw "github.com/calmdev/algorand-rewards/internal/ui/widget"
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

//go:embed assets/heart.svg
var heartIcon []byte
var heartIconResource = &fyne.StaticResource{
	StaticName:    "heart.svg",
	StaticContent: heartIcon,
}

func HeartIcon(size float32) fyne.CanvasObject {
	icon := theme.NewThemedResource(heartIconResource)
	img := canvas.NewImageFromResource(icon)
	img.SetMinSize(fyne.NewSize(size, size))
	img.FillMode = canvas.ImageFillContain

	return img
}

//go:embed assets/key.svg
var keyIcon []byte
var keyIconResource = &fyne.StaticResource{
	StaticName:    "key.svg",
	StaticContent: keyIcon,
}

func KeyIcon(size float32) fyne.CanvasObject {
	icon := theme.NewThemedResource(keyIconResource)
	img := canvas.NewImageFromResource(icon)
	img.SetMinSize(fyne.NewSize(size, size))
	img.FillMode = canvas.ImageFillContain

	return img
}

//go:embed assets/screwdriver-wrench.svg
var screwDriverWrenchIcon []byte
var screwDriverWrenchResource = &fyne.StaticResource{
	StaticName:    "screwdriver-wrench.svg",
	StaticContent: screwDriverWrenchIcon,
}

func ScrewDriverWrenchIcon(size float32) fyne.CanvasObject {
	icon := theme.NewThemedResource(screwDriverWrenchResource)
	img := canvas.NewImageFromResource(icon)
	img.SetMinSize(fyne.NewSize(size, size))
	img.FillMode = canvas.ImageFillContain

	return img
}

//go:embed assets/coins.svg
var coinsIcon []byte
var coinsIconResource = &fyne.StaticResource{
	StaticName:    "coins.svg",
	StaticContent: coinsIcon,
}

func CoinsIcon(size float32) fyne.CanvasObject {
	icon := theme.NewThemedResource(coinsIconResource)
	img := canvas.NewImageFromResource(icon)
	img.SetMinSize(fyne.NewSize(size, size))
	img.FillMode = canvas.ImageFillContain

	return img
}

//go:embed assets/snowflake.svg
var snowflakeIcon []byte
var snowflakeIconResource = &fyne.StaticResource{
	StaticName:    "snowflake.svg",
	StaticContent: snowflakeIcon,
}

func SnowflakeIcon(size float32) fyne.CanvasObject {
	icon := theme.NewThemedResource(snowflakeIconResource)
	img := canvas.NewImageFromResource(icon)
	img.SetMinSize(fyne.NewSize(size, size))
	img.FillMode = canvas.ImageFillContain

	return img
}

//go:embed assets/laptop-code.svg
var laptopCodeIcon []byte
var laptopCodeIconResource = &fyne.StaticResource{
	StaticName:    "laptop-code.svg",
	StaticContent: laptopCodeIcon,
}

func LaptopCodeIcon(size float32) fyne.CanvasObject {
	icon := theme.NewThemedResource(laptopCodeIconResource)
	img := canvas.NewImageFromResource(icon)
	img.SetMinSize(fyne.NewSize(size, size))
	img.FillMode = canvas.ImageFillContain

	return img
}

//go:embed assets/fingerprint.svg
var fingerprintIcon []byte
var fingerprintIconResource = &fyne.StaticResource{
	StaticName:    "fingerprint.svg",
	StaticContent: fingerprintIcon,
}

func FingerprintIcon(size float32) fyne.CanvasObject {
	icon := theme.NewThemedResource(fingerprintIconResource)
	img := canvas.NewImageFromResource(icon)
	img.SetMinSize(fyne.NewSize(size, size))
	img.FillMode = canvas.ImageFillContain

	return img
}

// AlgoIconResource returns icon resource.
func AlgoIconResource() *fyne.StaticResource {
	if app.CurrentApp().Settings().ThemeVariant() == theme.VariantLight {
		return AlgoBlackIconResource
	}
	return algoWhiteIconResource
}

// AlgoIcon returns the Algorand icon.
func AlgoIcon(size float32) fyne.CanvasObject {
	if app.CurrentApp().Settings().ThemeVariant() == theme.VariantLight {
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

// AlgoWordmark returns the white Algorand wordmark.
func AlgoWordmark(size float32) fyne.CanvasObject {
	// if theme is light return the black logo
	if app.CurrentApp().Settings().ThemeVariant() == theme.VariantLight {
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

// TransactionTypeIcon returns the icon for the transaction type.
func TransactionTypeIcon(t string, size float32) fyne.CanvasObject {
	switch t {
	case "pay":
		return AlgoIcon(size)
	case "hb":
		return HeartIcon(size)
	case "keyreg":
		return KeyIcon(size)
	case "acfg":
		return ScrewDriverWrenchIcon(size)
	case "axfer":
		return CoinsIcon(size)
	case "afrz":
		return SnowflakeIcon(size)
	case "appl":
		return LaptopCodeIcon(size)
	case "stpf":
		return FingerprintIcon(size)
	default:
		return AlgoIcon(size)
	}
}

// AccountStatusIcon returns the account status icon.
func AccountStatusIcon(account *algo.Account) fyne.Widget {
	var activity *iw.Activity

	// Check if account is eligible for rewards.
	if account.IncentiveEligible {
		activity = iw.NewActivity(DarkGreen, 20)
	} else {
		activity = iw.NewActivity(DarkRed, 20)
	}

	activity.Start()

	return activity
}
