package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// Color constants for graphs
var (
	GreenLight  = color.RGBA{R: 26, G: 155, B: 12, A: 255}  // Light theme green
	YellowLight = color.RGBA{R: 190, G: 161, B: 14, A: 255} // Light theme yellow
	RedLight    = color.RGBA{R: 186, G: 14, B: 23, A: 255}  // Light theme red

	GreenDark  = color.RGBA{R: 21, G: 222, B: 0, A: 255}  // Dark theme green
	YellowDark = color.RGBA{R: 255, G: 214, B: 0, A: 255} // Dark theme yellow
	RedDark    = color.RGBA{R: 252, G: 0, B: 13, A: 255}  // Dark theme red
)

// Color constants for the Memory tab's gauge, breakdown bars, and history
// chart. These identify a data series (memory, cached, buffers, swap) rather
// than a utilization status, so unlike the green/yellow/red set above they
// don't change meaning based on percentage.
var (
	BlueLight   = color.RGBA{R: 43, G: 92, B: 191, A: 255}   // Used/Memory series (light)
	PurpleLight = color.RGBA{R: 110, G: 76, B: 175, A: 255}  // Cached series (light)
	TealLight   = color.RGBA{R: 42, G: 101, B: 115, A: 255}  // Buffers series (light)
	GrayLight   = color.RGBA{R: 150, G: 156, B: 166, A: 255} // Free/track (light)

	BlueDark   = color.RGBA{R: 91, G: 142, B: 244, A: 255} // Used/Memory series (dark)
	PurpleDark = color.RGBA{R: 138, G: 99, B: 209, A: 255} // Cached series (dark)
	TealDark   = color.RGBA{R: 61, G: 122, B: 140, A: 255} // Buffers series (dark)
	GrayDark   = color.RGBA{R: 58, G: 63, B: 71, A: 255}   // Free/track (dark)
)

// GetSeriesColor returns the fixed identity color for a named data series,
// independent of utilization status (see GetGraphLineColor for that).
func GetSeriesColor(series string) color.RGBA {
	isDark := isDarkTheme()

	switch series {
	case "blue":
		if isDark {
			return BlueDark
		}
		return BlueLight
	case "purple":
		if isDark {
			return PurpleDark
		}
		return PurpleLight
	case "teal":
		if isDark {
			return TealDark
		}
		return TealLight
	case "gray":
		if isDark {
			return GrayDark
		}
		return GrayLight
	case "yellow":
		if isDark {
			return YellowDark
		}
		return YellowLight
	}

	if isDark {
		return GrayDark
	}
	return GrayLight
}

// isDarkTheme checks if the current theme variant is dark. It reads the
// resolved variant from Settings (which accounts for "auto"/OS-detected
// preference) rather than comparing theme instances - CustomTheme is always
// the app's active theme.Theme(), so a direct comparison against
// theme.LightTheme() never matches and previously made this always report
// dark, regardless of the user's actual selection.
func isDarkTheme() bool {
	return fyne.CurrentApp().Settings().ThemeVariant() == theme.VariantDark
}

// UtilizationStatus classifies a utilization percentage into a status band.
// This is the single source of truth for the green/yellow/red thresholds
// used by tile borders, status dots, bars, and graph line colors.
func UtilizationStatus(percent float64) string {
	switch {
	case percent >= 85:
		return "red"
	case percent >= 60:
		return "yellow"
	default:
		return "green"
	}
}

// severityRank orders status bands from least to most severe for comparison.
var severityRank = map[string]int{"green": 0, "yellow": 1, "red": 2}

// moreSevereStatus returns whichever of two status bands is more severe.
func moreSevereStatus(a, b string) string {
	if severityRank[b] > severityRank[a] {
		return b
	}
	return a
}

// GetGraphLineColor returns the appropriate color based on utilization status and theme
func GetGraphLineColor(status string) color.RGBA {
	isDark := isDarkTheme()

	switch status {
	case "green":
		if isDark {
			return GreenDark
		}
		return GreenLight
	case "yellow":
		if isDark {
			return YellowDark
		}
		return YellowLight
	case "red":
		if isDark {
			return RedDark
		}
		return RedLight
	}

	// Default to green
	if isDark {
		return GreenDark
	}
	return GreenLight
}

// ApplyTheme applies the specified theme variant to the application
func ApplyTheme(app fyne.App, themeVariant fyne.ThemeVariant) {
	app.Settings().SetTheme(&CustomTheme{variant: themeVariant})
}
