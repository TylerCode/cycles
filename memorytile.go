package main

import (
	"fmt"
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// MemoryTile represents a memory display tile
type MemoryTile struct {
	TitleLabel   *widget.Label
	TotalLabel   *widget.Label
	UsedLabel    *widget.Label
	FreeLabel    *widget.Label
	CachedLabel  *widget.Label
	PercentLabel *widget.Label
	container    *fyne.Container
	UsageHistory []float64 // Slice to store memory usage percentage history
	GraphImg     *canvas.Image
}

// NewMemoryTile creates a new memory tile with default styling
func NewMemoryTile(title string) *MemoryTile {
	titleLabel := widget.NewLabel(title)
	totalLabel := widget.NewLabel("Total: -- GB")
	usedLabel := widget.NewLabel("Used: -- GB")
	freeLabel := widget.NewLabel("Free: -- GB")
	cachedLabel := widget.NewLabel("Cached: -- GB")
	percentLabel := widget.NewLabel("Usage: --%")

	// Make title label bold by using a larger/emphasized style
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Create a background rectangle with rounded corners
	bg := canvas.NewRectangle(theme.BackgroundColor())
	bg.SetMinSize(fyne.NewSize(150, 150))
	bg.FillColor = theme.BackgroundColor()
	bg.StrokeColor = theme.ShadowColor()
	bg.StrokeWidth = 1
	bg.CornerRadius = 10

	graphImg := canvas.NewImageFromImage(image.NewRGBA(image.Rect(0, 0, 140, 60)))
	graphImg.FillMode = canvas.ImageFillOriginal

	container := container.NewMax(
		bg,
		container.NewVBox(
			titleLabel,
			percentLabel,
			totalLabel,
			usedLabel,
			freeLabel,
			cachedLabel,
			graphImg,
		),
	)

	return &MemoryTile{
		TitleLabel:   titleLabel,
		TotalLabel:   totalLabel,
		UsedLabel:    usedLabel,
		FreeLabel:    freeLabel,
		CachedLabel:  cachedLabel,
		PercentLabel: percentLabel,
		container:    container,
		GraphImg:     graphImg,
	}
}

// GetContainer returns the container of the memory tile
func (t *MemoryTile) GetContainer() fyne.CanvasObject {
	return t.container
}

// formatMemorySize formats bytes to human-readable format (GB, MB)
func formatMemorySize(kbytes uint64) string {
	gb := float64(kbytes) / 1024.0 / 1024.0
	if gb >= 1.0 {
		return fmt.Sprintf("%.2f GB", gb)
	}
	mb := float64(kbytes) / 1024.0
	return fmt.Sprintf("%.2f MB", mb)
}

// formatMemoryPercent formats memory usage percentage
func formatMemoryPercent(percent float64) string {
	return fmt.Sprintf("%.1f%%", percent)
}
