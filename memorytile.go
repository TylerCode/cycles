package main

import (
	"fmt"
	"image"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

// gaugeSize is the fixed size of the radial usage gauge.
var gaugeSize = fyne.NewSize(160, 160)

// breakdownRow is one labeled, colored bar in the Memory tab's breakdown
// panel (Used/Cached/Buffers/Free/Swap).
type breakdownRow struct {
	container   fyne.CanvasObject
	bar         *Bar
	swatchRect  *canvas.Rectangle
	series      string
	labelText   *canvas.Text
	valueText   *canvas.Text
	percentText *canvas.Text
}

func newBreakdownRow(label, series string) *breakdownRow {
	swatchRect := canvas.NewRectangle(GetSeriesColor(series))
	swatchRect.CornerRadius = 2

	labelText := canvas.NewText(label, theme.ForegroundColor())
	labelText.TextSize = 13
	labelBox := container.NewGridWrap(fyne.NewSize(64, 18), labelText)

	leading := container.NewHBox(swatch(swatchRect), labelBox)

	bar := NewBar(theme.PrimaryColor(), theme.InputBackgroundColor(), false)

	valueText := canvas.NewText("--", theme.ForegroundColor())
	valueText.TextSize = 13
	valueText.TextStyle = fyne.TextStyle{Bold: true}
	valueText.Alignment = fyne.TextAlignTrailing
	valueBox := container.NewGridWrap(fyne.NewSize(70, 18), valueText)

	percentText := canvas.NewText("--", theme.PlaceHolderColor())
	percentText.TextSize = 12
	percentText.Alignment = fyne.TextAlignTrailing
	percentBox := container.NewGridWrap(fyne.NewSize(46, 18), percentText)

	trailing := container.NewHBox(valueBox, percentBox)

	row := container.NewBorder(nil, nil, leading, trailing, bar.Container())

	return &breakdownRow{
		container: row, bar: bar, swatchRect: swatchRect, series: series,
		labelText: labelText, valueText: valueText, percentText: percentText,
	}
}

func swatch(c fyne.CanvasObject) fyne.CanvasObject {
	return container.NewGridWrap(fyne.NewSize(10, 10), c)
}

func (r *breakdownRow) update(percent float64, value, percentText string) {
	r.bar.SetFillColor(GetSeriesColor(r.series))
	r.bar.SetPercent(percent, "", nil)
	r.valueText.Text = value
	r.valueText.Refresh()
	r.percentText.Text = percentText
	r.percentText.Refresh()
}

// refreshTheme re-applies theme-derived colors — see CoreTile.RefreshTheme
// for why this is needed.
func (r *breakdownRow) refreshTheme() {
	r.labelText.Color = theme.ForegroundColor()
	r.labelText.Refresh()
	r.valueText.Color = theme.ForegroundColor()
	r.valueText.Refresh()
	r.percentText.Color = theme.PlaceHolderColor()
	r.percentText.Refresh()
	r.bar.SetTrackColor(theme.InputBackgroundColor())
	r.swatchRect.FillColor = GetSeriesColor(r.series)
	r.swatchRect.Refresh()
}

// MemoryDashboard is the Memory tab's whole layout: a radial usage gauge,
// used/cached/buffers/free/swap breakdown rows, and a full-width memory+swap
// history chart. It replaces the old copy-pasted single grid tile — memory
// is one aggregate metric with named sub-components, not a repeating
// collection, so it gets a dashboard rather than a tile grid.
type MemoryDashboard struct {
	container fyne.CanvasObject

	UsageHistory []float64
	SwapHistory  []float64

	gaugePercentText *canvas.Text
	gaugeSubText     *canvas.Text
	gaugePercent     float64

	used, cached, buffers, free, swap *breakdownRow

	divider      *canvas.Rectangle
	chartTitle   *canvas.Text
	memoryLegend *legendItem
	swapLegend   *legendItem

	gaugeRaster *canvas.Raster
	chartRaster *canvas.Raster
}

// legendItem is a small color swatch + label used in the history chart's
// legend (e.g. "■ Memory").
type legendItem struct {
	swatch    *canvas.Rectangle
	labelText *canvas.Text
}

func newLegendItem(label string, swatchColor color.Color) (*legendItem, fyne.CanvasObject) {
	rect := canvas.NewRectangle(swatchColor)
	rect.CornerRadius = 2
	labelText := canvas.NewText(label, theme.PlaceHolderColor())
	labelText.TextSize = 11
	item := &legendItem{swatch: rect, labelText: labelText}
	return item, container.NewHBox(swatch(rect), labelText)
}

func (l *legendItem) refreshTheme(swatchColor color.Color) {
	l.swatch.FillColor = swatchColor
	l.swatch.Refresh()
	l.labelText.Color = theme.PlaceHolderColor()
	l.labelText.Refresh()
}

// NewMemoryDashboard creates the Memory tab's dashboard layout.
func NewMemoryDashboard() *MemoryDashboard {
	d := &MemoryDashboard{}

	d.gaugeRaster = canvas.NewRaster(func(w, h int) image.Image {
		return DrawRadialGauge(w, h, d.gaugePercent, GetSeriesColor("blue"), GetSeriesColor("gray"))
	})
	d.gaugeRaster.SetMinSize(gaugeSize)

	d.gaugePercentText = canvas.NewText("--%", theme.ForegroundColor())
	d.gaugePercentText.TextStyle = fyne.TextStyle{Bold: true}
	d.gaugePercentText.TextSize = 26
	d.gaugePercentText.Alignment = fyne.TextAlignCenter

	d.gaugeSubText = canvas.NewText("of -- GB", theme.PlaceHolderColor())
	d.gaugeSubText.TextSize = 12
	d.gaugeSubText.Alignment = fyne.TextAlignCenter

	gaugeLabels := container.NewVBox(d.gaugePercentText, d.gaugeSubText)
	gaugeStack := container.NewStack(container.NewGridWrap(gaugeSize, d.gaugeRaster), container.NewCenter(gaugeLabels))
	gaugeBox := container.NewCenter(gaugeStack)

	d.used = newBreakdownRow("Used", "blue")
	d.cached = newBreakdownRow("Cached", "purple")
	d.buffers = newBreakdownRow("Buffers", "teal")
	d.free = newBreakdownRow("Free", "gray")
	d.swap = newBreakdownRow("Swap", "yellow")

	d.divider = canvas.NewRectangle(theme.ShadowColor())
	d.divider.SetMinSize(fyne.NewSize(0, 1))

	breakdown := container.NewVBox(
		d.used.container,
		d.cached.container,
		d.buffers.container,
		d.free.container,
		d.divider,
		d.swap.container,
	)

	top := container.NewBorder(nil, nil, gaugeBox, nil, container.NewPadded(breakdown))

	var memoryLegendBox, swapLegendBox fyne.CanvasObject
	d.memoryLegend, memoryLegendBox = newLegendItem("Memory", GetSeriesColor("blue"))
	d.swapLegend, swapLegendBox = newLegendItem("Swap", GetSeriesColor("yellow"))

	d.chartTitle = canvas.NewText("Usage history", theme.ForegroundColor())
	d.chartTitle.TextStyle = fyne.TextStyle{Bold: true}
	d.chartTitle.TextSize = 13

	chartHeader := container.NewBorder(nil, nil, d.chartTitle, container.NewHBox(memoryLegendBox, swapLegendBox))

	d.chartRaster = canvas.NewRaster(func(w, h int) image.Image {
		return DrawAreaChart(w, h, d.UsageHistory, d.SwapHistory)
	})

	chartPanel := container.NewBorder(chartHeader, nil, nil, nil, d.chartRaster)

	d.container = container.NewBorder(container.NewPadded(top), nil, nil, nil, container.NewPadded(chartPanel))

	return d
}

// GetContainer returns the dashboard's canvas object.
func (d *MemoryDashboard) GetContainer() fyne.CanvasObject {
	return d.container
}

// Update refreshes the dashboard with the latest memory and swap readings.
func (d *MemoryDashboard) Update(mem MemoryInfo, swap SwapInfo, usagePercent, swapPercent float64) {
	d.gaugePercent = usagePercent
	d.gaugePercentText.Text = formatMemoryPercent(usagePercent)
	d.gaugePercentText.Refresh()
	d.gaugeSubText.Text = "of " + formatMemorySize(mem.Total)
	d.gaugeSubText.Refresh()
	d.gaugeRaster.Refresh()

	pct := func(part, total uint64) float64 {
		if total == 0 {
			return 0
		}
		return float64(part) / float64(total) * 100
	}

	d.used.update(usagePercent, formatMemorySize(mem.Used), formatMemoryPercent(usagePercent))
	d.cached.update(pct(mem.Cached, mem.Total), formatMemorySize(mem.Cached), formatMemoryPercent(pct(mem.Cached, mem.Total)))
	d.buffers.update(pct(mem.Buffers, mem.Total), formatMemorySize(mem.Buffers), formatMemoryPercent(pct(mem.Buffers, mem.Total)))
	d.free.update(pct(mem.Free, mem.Total), formatMemorySize(mem.Free), formatMemoryPercent(pct(mem.Free, mem.Total)))
	d.swap.update(swapPercent, formatMemorySize(swap.Used), "of "+formatMemorySize(swap.Total))

	d.chartRaster.Refresh()
}

// RefreshTheme re-applies theme-derived colors across the whole dashboard.
// Needed because canvas.Text/Rectangle primitives only read theme colors
// once, at construction — see CoreTile.RefreshTheme.
func (d *MemoryDashboard) RefreshTheme() {
	d.gaugePercentText.Color = theme.ForegroundColor()
	d.gaugePercentText.Refresh()
	d.gaugeSubText.Color = theme.PlaceHolderColor()
	d.gaugeSubText.Refresh()
	d.gaugeRaster.Refresh()

	for _, row := range []*breakdownRow{d.used, d.cached, d.buffers, d.free, d.swap} {
		row.refreshTheme()
	}

	d.divider.FillColor = theme.ShadowColor()
	d.divider.Refresh()

	d.chartTitle.Color = theme.ForegroundColor()
	d.chartTitle.Refresh()

	d.memoryLegend.refreshTheme(GetSeriesColor("blue"))
	d.swapLegend.refreshTheme(GetSeriesColor("yellow"))

	d.chartRaster.Refresh()
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
