package main

import (
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// CoreTile represents a single CPU core's display tile
type CoreTile struct {
	CoreLabel   *widget.Label
	UtilLabel   *widget.Label
	ClockLabel  *widget.Label
	container   *fyne.Container
	UtilHistory []float64 // Slice to store utilization history
	GraphImg    *canvas.Image
}

// NewCoreTile creates a new core tile with default styling
func NewCoreTile() *CoreTile {
	coreLabel := widget.NewLabel("Core #")
	utilLabel := widget.NewLabel("Util %")
	clockLabel := widget.NewLabel("Clock MHz")

	// Create a background rectangle with rounded corners
	bg := canvas.NewRectangle(theme.BackgroundColor())
	bg.SetMinSize(fyne.NewSize(100, 100))
	bg.FillColor = theme.BackgroundColor()
	bg.StrokeColor = theme.ShadowColor()
	bg.StrokeWidth = 1
	bg.CornerRadius = 10

	graphImg := canvas.NewImageFromImage(image.NewRGBA(image.Rect(0, 0, 100, 50)))
	graphImg.FillMode = canvas.ImageFillOriginal

	container := container.NewMax(bg, container.NewVBox(coreLabel, utilLabel, clockLabel, graphImg))

	return &CoreTile{
		CoreLabel:  coreLabel,
		UtilLabel:  utilLabel,
		ClockLabel: clockLabel,
		container:  container,
		GraphImg:   graphImg,
	}
}

// GetContainer returns the container of the core tile
func (t *CoreTile) GetContainer() fyne.CanvasObject {
	return t.container
}
