package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var (
	// Colors
	Black     = color.Black
	Grey      = color.RGBA{128, 128, 128, 255}
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
		return Black
	}

	if name == theme.ColorNameForeground {
		if variant == theme.VariantLight {
			return Black
		}
		return White
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
