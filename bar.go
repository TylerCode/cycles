package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

// barLayout resizes a track (full width) and a fill (percent-of-width)
// rectangle, plus an optional overlaid label, to fill the container. It
// implements fyne.Layout so the bar reflows automatically whenever its
// parent container is resized, not just when SetPercent is called.
type barLayout struct {
	percent float64
}

func (b *barLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(objects) < 2 {
		return
	}
	track, fill := objects[0], objects[1]

	track.Move(fyne.NewPos(0, 0))
	track.Resize(size)

	fillWidth := size.Width * float32(b.percent/100)
	if fillWidth < 0 {
		fillWidth = 0
	}
	if fillWidth > size.Width {
		fillWidth = size.Width
	}
	fill.Move(fyne.NewPos(0, 0))
	fill.Resize(fyne.NewSize(fillWidth, size.Height))

	if len(objects) > 2 {
		label := objects[2]
		label.Move(fyne.NewPos(4, 0))
		label.Resize(fyne.NewSize(size.Width-4, size.Height))
	}
}

func (b *barLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(40, 8)
}

// Bar is a colored horizontal progress bar that reflows with its container.
// It's used for the CPU list view's inline utilization bar and the Memory
// tab's used/cached/buffers/free/swap breakdown rows, each of which needs
// its own fixed fill color rather than the single theme primary color that
// widget.ProgressBar is hardcoded to use.
type Bar struct {
	container *fyne.Container
	layout    *barLayout
	track     *canvas.Rectangle
	fill      *canvas.Rectangle
	label     *canvas.Text
}

// NewBar creates a bar. If showLabel is true, a text label is overlaid on
// top of the bar (used by the CPU list view's inline "48.5%" style label).
func NewBar(fillColor, trackColor color.Color, showLabel bool) *Bar {
	track := canvas.NewRectangle(trackColor)
	track.CornerRadius = 3
	fill := canvas.NewRectangle(fillColor)
	fill.CornerRadius = 3

	layout := &barLayout{}
	objects := []fyne.CanvasObject{track, fill}

	var label *canvas.Text
	if showLabel {
		label = canvas.NewText("", color.White)
		label.TextSize = 10
		label.TextStyle = fyne.TextStyle{Bold: true}
		objects = append(objects, label)
	}

	c := container.New(layout, objects...)

	return &Bar{container: c, layout: layout, track: track, fill: fill, label: label}
}

// Container returns the bar's canvas object for embedding in a layout.
func (b *Bar) Container() fyne.CanvasObject {
	return b.container
}

// SetFillColor updates the bar's fill color (used when the color depends on
// a status that can change, e.g. CPU utilization bands).
func (b *Bar) SetFillColor(c color.Color) {
	b.fill.FillColor = c
}

// SetTrackColor updates the bar's track (background) color and repaints it.
// Used to re-apply the theme's color when the app theme changes live, since
// the track color is otherwise only ever read once at construction time.
func (b *Bar) SetTrackColor(c color.Color) {
	b.track.FillColor = c
	b.track.Refresh()
}

// SetPercent updates the fill width and, if this bar has a label, its text
// and color. It re-runs the layout immediately so the change is visible
// without waiting for an unrelated resize.
func (b *Bar) SetPercent(percent float64, labelText string, labelColor color.Color) {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	b.layout.percent = percent

	if b.label != nil {
		b.label.Text = labelText
		if labelColor != nil {
			b.label.Color = labelColor
		}
	}

	b.layout.Layout(b.container.Objects, b.container.Size())
	b.container.Refresh()
}
