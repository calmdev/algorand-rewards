package ui

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"github.com/calmdev/algorand-rewards/internal/app"
)

var (
	// Colors
	Primary     = color.RGBA{128, 128, 128, 255}
	Transparent = color.Transparent
	White       = color.White
	Black       = color.Black
	Grey        = color.RGBA{128, 128, 128, 255}
	LightGrey   = color.RGBA{211, 211, 211, 255}
	DarkGrey    = color.RGBA{23, 23, 24, 255}
	DarkGreen   = color.RGBA{0, 128, 0, 255}
	DarkRed     = color.RGBA{234, 47, 73, 255}
)

// AppTheme is a custom theme for the application.
type AppTheme struct{}

// Assert CustomTheme implements fyne.Theme interface.
var _ fyne.Theme = (*AppTheme)(nil)

// Color returns a color for the theme.
func (m AppTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNamePrimary {
		if variant == theme.VariantLight {
			return LightGrey
		}
		return Primary
	}

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
		return color.RGBA{R: 30, G: 30, B: 30, A: 255}
	}

	if name == theme.ColorNameOverlayBackground {
		if variant == theme.VariantLight {
			return White
		}
		return DarkGrey
	}

	if name == theme.ColorNameShadow {
		if variant == theme.VariantLight {
			return White
		}
		return DarkGrey
	}

	if name == theme.ColorNameHyperlink {
		if variant == theme.VariantLight {
			return Black
		}
		return White
	}

	return theme.DefaultTheme().Color(name, variant)
}

// Icon returns an icon for the theme.
func (m AppTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Font returns a font for the theme.
func (m AppTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

// Size returns a size for the theme.
func (m AppTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

// WatchForThemeVariantChanges watches for OS theme variant changes.
func WatchForThemeVariantChanges(a *app.App, w fyne.Window, themeVariant *fyne.ThemeVariant) {
	for {
		if a.Settings().ThemeVariant() != *themeVariant {
			*themeVariant = a.Settings().ThemeVariant()
			currentView := Layout.currentView
			w.SetContent(RenderLayout())
			RenderView(currentView)
		}

		time.Sleep(500 * time.Millisecond)
	}
}
