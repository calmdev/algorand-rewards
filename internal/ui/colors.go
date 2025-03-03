package ui

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var (
	// Colors
	Black     = color.Black
	Grey      = color.RGBA{128, 128, 128, 255}
	LightGrey = color.RGBA{211, 211, 211, 255}
	DarkGrey  = color.RGBA{23, 23, 24, 100}
	DarkGreen = color.RGBA{0, 128, 0, 255}
	White     = color.White
)

// CustomTheme is a custom theme for the application.
type CustomTheme struct{}

// Assert CustomTheme implements fyne.Theme interface.
var _ fyne.Theme = (*CustomTheme)(nil)

// Color returns a color for the theme.
func (m CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameBackground {
		if variant == theme.VariantLight {
			return White
		}
		return DarkGrey
	}

	if name == theme.ColorNameForeground {
		if variant == theme.VariantLight {
			return Black
		}
		return White
	}

	if name == theme.ColorNameSeparator {
		if variant == theme.VariantLight {
			return LightGrey
		}
		return Black
	}

	return theme.DefaultTheme().Color(name, variant)
}

// Icon returns an icon for the theme.
func (m CustomTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Font returns a font for the theme.
func (m CustomTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

// Size returns a size for the theme.
func (m CustomTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

// WatchForThemeVariantChanges watches for OS theme variant changes.
func WatchForThemeVariantChanges(a fyne.App, w fyne.Window, themeVariant *fyne.ThemeVariant) {
	for {
		if a.Settings().ThemeVariant() != *themeVariant {
			*themeVariant = a.Settings().ThemeVariant()
			w.SetContent(RenderMainView())
		}

		time.Sleep(500 * time.Millisecond)
	}
}
