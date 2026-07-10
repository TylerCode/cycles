package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// coreTileSize is the fixed cell size used by the Tiles view's reflowing
// grid (see container.NewGridWrap in NewCPUTab). Fyne's GridWrap reflows the
// column count to fit available width but doesn't stretch cells, so a fixed
// size here is what actually changes as the window is resized: the number
// of columns, not the tile size.
var coreTileSize = fyne.NewSize(168, 112)

// CPUTab builds and owns the whole CPU tab: an aggregate stats strip, a
// Tiles/List view toggle, and both view containers (kept alive
// simultaneously — see CoreState — so toggling views never shows stale
// data).
type CPUTab struct {
	Content fyne.CanvasObject
	Cores   []*CoreState

	threadsLabel *canvas.Text
	avgLabel     *canvas.Text
	peakLabel    *canvas.Text
	clockLabel   *canvas.Text

	threadsValue *canvas.Text
	avgValue     *canvas.Text
	peakValue    *canvas.Text
	clockValue   *canvas.Text

	tilesBtn *widget.Button
	listBtn  *widget.Button
	tilesC   fyne.CanvasObject
	listC    fyne.CanvasObject

	onViewChange func(view string)
}

// NewCPUTab creates the CPU tab for numCores cores. initialView is "tiles"
// or "list". onViewChange, if non-nil, is called whenever the user toggles
// the view so the caller can persist the choice.
func NewCPUTab(numCores int, initialView string, onViewChange func(view string)) *CPUTab {
	t := &CPUTab{onViewChange: onViewChange}

	t.Cores = make([]*CoreState, numCores)
	tiles := make([]fyne.CanvasObject, numCores)
	rows := make([]fyne.CanvasObject, numCores)
	for i := 0; i < numCores; i++ {
		state := NewCoreState(i)
		t.Cores[i] = state
		tiles[i] = NewCoreTile(state).GetContainer()
		rows[i] = NewCoreListRow(state).GetContainer()
	}

	tilesGrid := container.NewGridWrap(coreTileSize, tiles...)
	t.tilesC = container.NewVScroll(tilesGrid)

	// Split rows into two side-by-side column groups, matching the mockup's
	// list view (this stays two columns regardless of window width, per
	// design decision). HSplit gives them a visible divider — without one,
	// the right column's label sits flush against the left column's
	// trailing clock/sparkline group, reading as one merged row.
	half := (numCores + 1) / 2
	col1 := container.NewVBox(rows[:half]...)
	col2 := container.NewVBox(rows[half:]...)
	listSplit := container.NewHSplit(col1, col2)
	listSplit.Offset = 0.5
	t.listC = container.NewVScroll(listSplit)

	views := container.NewStack(t.tilesC, t.listC)

	statsStrip := t.buildStatsStrip()
	toggle := t.buildViewToggle()
	header := container.NewBorder(nil, nil, statsStrip, toggle)

	t.Content = container.NewBorder(container.NewPadded(header), nil, nil, nil, views)

	t.setView(initialView, false)

	return t
}

func (t *CPUTab) buildStatsStrip() fyne.CanvasObject {
	stat := func(label string) (fyne.CanvasObject, *canvas.Text, *canvas.Text) {
		labelText := canvas.NewText(label, theme.PlaceHolderColor())
		labelText.TextSize = 11
		valueText := canvas.NewText("--", theme.ForegroundColor())
		valueText.TextStyle = fyne.TextStyle{Bold: true}
		valueText.TextSize = 15
		return container.NewVBox(labelText, valueText), labelText, valueText
	}

	var threads, avg, peak, clock fyne.CanvasObject
	threads, t.threadsLabel, t.threadsValue = stat("Threads")
	avg, t.avgLabel, t.avgValue = stat("Avg util")
	peak, t.peakLabel, t.peakValue = stat("Peak core")
	clock, t.clockLabel, t.clockValue = stat("Max clock")

	return container.NewHBox(threads, avg, peak, clock)
}

func (t *CPUTab) buildViewToggle() fyne.CanvasObject {
	t.tilesBtn = widget.NewButton("Tiles", func() { t.setView("tiles", true) })
	t.listBtn = widget.NewButton("List", func() { t.setView("list", true) })
	return container.NewHBox(t.tilesBtn, layout.NewSpacer(), t.listBtn)
}

func (t *CPUTab) setView(view string, notify bool) {
	if view == "list" {
		t.tilesC.Hide()
		t.listC.Show()
		t.tilesBtn.Importance = widget.LowImportance
		t.listBtn.Importance = widget.HighImportance
	} else {
		view = "tiles"
		t.listC.Hide()
		t.tilesC.Show()
		t.tilesBtn.Importance = widget.HighImportance
		t.listBtn.Importance = widget.LowImportance
	}
	t.tilesBtn.Refresh()
	t.listBtn.Refresh()

	if notify && t.onViewChange != nil {
		t.onViewChange(view)
	}
}

// UpdateStats refreshes the aggregate stats strip (Threads/Avg util/Peak
// core/Max clock).
func (t *CPUTab) UpdateStats(stats CPUAggregateStats) {
	t.threadsValue.Text = formatThreadsValue(stats.Threads)
	t.threadsValue.Refresh()

	t.avgValue.Text = formatUtilLabel(stats.AvgUtil)
	t.avgValue.Refresh()

	t.peakValue.Text = formatPeakCoreValue(stats.PeakCoreIndex, stats.PeakCoreUtil)
	t.peakValue.Refresh()

	t.clockValue.Text = formatClockLabel(stats.MaxClock)
	t.clockValue.Refresh()
}

// RefreshTheme re-applies theme-derived colors across the stats strip and
// every core's tile/row. Needed because canvas.Text/Rectangle primitives
// only read theme colors once, at construction — see CoreTile.RefreshTheme.
func (t *CPUTab) RefreshTheme() {
	for _, label := range []*canvas.Text{t.threadsLabel, t.avgLabel, t.peakLabel, t.clockLabel} {
		label.Color = theme.PlaceHolderColor()
		label.Refresh()
	}
	for _, value := range []*canvas.Text{t.threadsValue, t.avgValue, t.peakValue, t.clockValue} {
		value.Color = theme.ForegroundColor()
		value.Refresh()
	}
	for _, core := range t.Cores {
		core.RefreshTheme()
	}
}
