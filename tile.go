package main

import (
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

// CoreState holds one CPU core's live data. It's the single source of truth
// that both the Tiles view (CoreTile) and the List view (CoreListRow)
// render from, so switching between views never shows stale data — whichever
// widgets exist for a core are refreshed together on every update.
type CoreState struct {
	Index       int
	Percent     float64
	ClockMHz    float64
	UtilHistory []float64

	tile *CoreTile
	row  *CoreListRow
}

// NewCoreState creates the state holder for one core.
func NewCoreState(index int) *CoreState {
	return &CoreState{Index: index}
}

// Update records a new sample and refreshes any attached widgets.
func (c *CoreState) Update(percent, clockMHz float64, historySize int) {
	c.Percent = percent
	c.ClockMHz = clockMHz

	c.UtilHistory = append(c.UtilHistory, percent)
	if len(c.UtilHistory) > historySize {
		c.UtilHistory = c.UtilHistory[1:]
	}

	if c.tile != nil {
		c.tile.refresh()
	}
	if c.row != nil {
		c.row.refresh()
	}
}

// RefreshTheme re-applies theme-derived colors on this core's tile and row.
func (c *CoreState) RefreshTheme() {
	if c.tile != nil {
		c.tile.RefreshTheme()
	}
	if c.row != nil {
		c.row.RefreshTheme()
	}
}

// CoreTile is a single core's card in the CPU tab's Tiles view: core name +
// status dot, a large percentage, clock speed, and a sparkline that renders
// at the tile's actual displayed size.
type CoreTile struct {
	state *CoreState

	container   *fyne.Container
	bg          *canvas.Rectangle
	dot         *canvas.Circle
	coreLabel   *canvas.Text
	percentText *canvas.Text
	clockLabel  *canvas.Text
	sparkline   *canvas.Raster
}

// NewCoreTile creates a tile bound to the given core state.
func NewCoreTile(state *CoreState) *CoreTile {
	bg := canvas.NewRectangle(theme.InputBackgroundColor())
	bg.StrokeColor = theme.ShadowColor()
	bg.StrokeWidth = 1
	bg.CornerRadius = 8

	coreLabel := canvas.NewText(formatCoreLabel(state.Index), theme.ForegroundColor())
	coreLabel.TextStyle = fyne.TextStyle{Bold: true}
	coreLabel.TextSize = 12

	dot := canvas.NewCircle(GetGraphLineColor("green"))
	dotBox := container.NewGridWrap(fyne.NewSize(7, 7), dot)

	header := container.NewBorder(nil, nil, coreLabel, dotBox)

	percentText := canvas.NewText("--%", theme.ForegroundColor())
	percentText.TextStyle = fyne.TextStyle{Bold: true}
	percentText.TextSize = 22

	clockLabel := canvas.NewText("-- MHz", theme.PlaceHolderColor())
	clockLabel.TextSize = 11

	t := &CoreTile{state: state, bg: bg, dot: dot, coreLabel: coreLabel, percentText: percentText, clockLabel: clockLabel}

	t.sparkline = canvas.NewRaster(func(w, h int) image.Image {
		return DrawSparkline(w, h, t.state.UtilHistory)
	})

	info := container.NewVBox(header, percentText, clockLabel)
	body := container.NewBorder(info, nil, nil, nil, t.sparkline)
	padded := container.NewPadded(body)

	t.container = container.NewStack(bg, padded)
	state.tile = t
	return t
}

// GetContainer returns the tile's canvas object.
func (t *CoreTile) GetContainer() fyne.CanvasObject {
	return t.container
}

func (t *CoreTile) refresh() {
	status := UtilizationStatus(t.state.Percent)
	statusColor := GetGraphLineColor(status)

	t.percentText.Text = formatUtilLabel(t.state.Percent)
	t.percentText.Refresh()

	t.clockLabel.Text = formatClockLabel(t.state.ClockMHz)
	t.clockLabel.Refresh()

	t.dot.FillColor = statusColor
	t.dot.Refresh()

	t.sparkline.Refresh()

	if status == "red" {
		t.bg.StrokeColor = statusColor
		t.bg.StrokeWidth = 2
	} else {
		t.bg.StrokeColor = theme.ShadowColor()
		t.bg.StrokeWidth = 1
	}
	t.bg.Refresh()
}

// RefreshTheme re-applies theme-derived colors. canvas.Text and
// canvas.Rectangle primitives only read theme colors once, at construction,
// so switching the app theme live (without restarting) would otherwise leave
// these stuck on the old theme's colors even though built-in widgets like
// the Tiles/List buttons update automatically.
func (t *CoreTile) RefreshTheme() {
	t.bg.FillColor = theme.InputBackgroundColor()
	t.coreLabel.Color = theme.ForegroundColor()
	t.percentText.Color = theme.ForegroundColor()
	t.clockLabel.Color = theme.PlaceHolderColor()
	t.coreLabel.Refresh()
	t.percentText.Refresh()
	t.clockLabel.Refresh()
	t.refresh()
}

// CoreListRow is a single core's row in the CPU tab's List view: core name,
// an inline utilization bar with its percentage overlaid, clock speed, and a
// compact sparkline.
type CoreListRow struct {
	state *CoreState

	container  *fyne.Container
	bar        *Bar
	coreLabel  *canvas.Text
	clockLabel *canvas.Text
	sparkline  *canvas.Raster
}

// NewCoreListRow creates a list row bound to the given core state.
func NewCoreListRow(state *CoreState) *CoreListRow {
	coreLabel := canvas.NewText(formatCoreLabel(state.Index), theme.ForegroundColor())
	coreLabel.TextStyle = fyne.TextStyle{Bold: true}
	coreLabel.TextSize = 12
	coreLabelBox := container.NewGridWrap(fyne.NewSize(64, 20), coreLabel)

	bar := NewBar(GetGraphLineColor("green"), theme.InputBackgroundColor(), true)

	clockLabel := canvas.NewText("--", theme.PlaceHolderColor())
	clockLabel.TextSize = 11
	clockLabel.Alignment = fyne.TextAlignTrailing
	clockLabelBox := container.NewGridWrap(fyne.NewSize(44, 20), clockLabel)

	r := &CoreListRow{state: state, bar: bar, coreLabel: coreLabel, clockLabel: clockLabel}

	r.sparkline = canvas.NewRaster(func(w, h int) image.Image {
		return DrawSparkline(w, h, r.state.UtilHistory)
	})
	sparklineBox := container.NewGridWrap(fyne.NewSize(64, 20), r.sparkline)

	trailing := container.NewHBox(clockLabelBox, sparklineBox)

	// NewBorder stretches its center object (the bar) to fill the remaining
	// row width between the fixed-width label and trailer, which is what
	// makes the bar flexible rather than a fixed cell size like the others.
	r.container = container.NewBorder(nil, nil, coreLabelBox, trailing, bar.Container())
	state.row = r
	return r
}

// GetContainer returns the row's canvas object.
func (r *CoreListRow) GetContainer() fyne.CanvasObject {
	return r.container
}

func (r *CoreListRow) refresh() {
	status := UtilizationStatus(r.state.Percent)
	statusColor := GetGraphLineColor(status)

	labelColor := theme.ForegroundColor()
	if status != "green" {
		labelColor = theme.BackgroundColor()
	}

	r.bar.SetFillColor(statusColor)
	r.bar.SetPercent(r.state.Percent, formatUtilLabel(r.state.Percent), labelColor)

	r.clockLabel.Text = formatClockNumber(r.state.ClockMHz)
	r.clockLabel.Refresh()

	r.sparkline.Refresh()
}

// RefreshTheme re-applies theme-derived colors — see CoreTile.RefreshTheme
// for why this is needed.
func (r *CoreListRow) RefreshTheme() {
	r.coreLabel.Color = theme.ForegroundColor()
	r.coreLabel.Refresh()
	r.clockLabel.Color = theme.PlaceHolderColor()
	r.clockLabel.Refresh()
	r.bar.SetTrackColor(theme.InputBackgroundColor())
	r.refresh()
}
