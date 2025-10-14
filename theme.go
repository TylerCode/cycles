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

// GetGraphLineColor returns the appropriate color based on utilization status and theme
func GetGraphLineColor(status string) color.RGBA {
	currentTheme := fyne.CurrentApp().Settings().Theme()
	isDark := true

	// Check if the current theme is light
	if currentTheme == theme.LightTheme() {
		isDark = false
	}

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
