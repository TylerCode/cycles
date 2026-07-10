package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
)

// DrawSparkline renders a utilization history line into an image sized w×h.
// Unlike the old fixed-bitmap DrawGraph, w and h are the actual pixel size
// the sparkline is being displayed at (driven by canvas.Raster's generate
// callback), so it always renders crisply at its real on-screen size.
func DrawSparkline(w, h int, data []float64) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	if w <= 0 || h <= 0 {
		return dst
	}

	if len(data) < 2 {
		return dst
	}

	xStep := float64(w-1) / float64(len(data)-1)
	yAt := func(v float64) int {
		y := h - 1 - int(v/100*float64(h-1))
		return max(0, min(y, h-1))
	}

	for i := 0; i < len(data)-1; i++ {
		x1 := int(float64(i) * xStep)
		y1 := yAt(data[i])
		x2 := int(float64(i+1) * xStep)
		y2 := yAt(data[i+1])

		status := moreSevereStatus(UtilizationStatus(data[i]), UtilizationStatus(data[i+1]))
		lineColor := GetGraphLineColor(status)

		drawLineWithEffect(dst, x1, y1, x2, y2, lineColor)
	}

	return dst
}

// DrawAreaChart renders a filled area for the primary series (e.g. memory
// usage %) with a plain line for the secondary series (e.g. swap usage %)
// on top, plus light horizontal gridlines. Both series are expected in the
// 0-100 range. w and h are the actual pixel size being rendered at.
func DrawAreaChart(w, h int, primary, secondary []float64) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	if w <= 0 || h <= 0 {
		return dst
	}

	draw.Draw(dst, dst.Bounds(), &image.Uniform{color.Transparent}, image.Point{}, draw.Src)

	gridColor := color.NRGBA{R: 128, G: 128, B: 128, A: 60}
	for _, frac := range []float64{0.25, 0.5, 0.75} {
		y := int(float64(h-1) * frac)
		drawLine(dst, 0, y, w-1, y, gridColor)
	}

	if len(primary) < 2 {
		return dst
	}

	primaryColor := GetSeriesColor("blue")
	fillColor := color.NRGBA{R: primaryColor.R, G: primaryColor.G, B: primaryColor.B, A: 40}

	valueAt := func(series []float64, x int) float64 {
		if len(series) < 2 {
			if len(series) == 1 {
				return series[0]
			}
			return 0
		}
		xStep := float64(w-1) / float64(len(series)-1)
		pos := float64(x) / xStep
		i := int(pos)
		if i >= len(series)-1 {
			return series[len(series)-1]
		}
		frac := pos - float64(i)
		return series[i]*(1-frac) + series[i+1]*frac
	}

	yAt := func(v float64) int {
		y := h - 1 - int(v/100*float64(h-1))
		return max(0, min(y, h-1))
	}

	for x := 0; x < w; x++ {
		top := yAt(valueAt(primary, x))
		for y := top; y < h; y++ {
			dst.Set(x, y, fillColor)
		}
	}
	for i := 0; i < w-1; i++ {
		y1 := yAt(valueAt(primary, i))
		y2 := yAt(valueAt(primary, i+1))
		drawLineWithEffect(dst, i, y1, i+1, y2, primaryColor)
	}

	if len(secondary) >= 2 {
		secondaryColor := GetSeriesColor("yellow")
		for i := 0; i < w-1; i++ {
			y1 := yAt(valueAt(secondary, i))
			y2 := yAt(valueAt(secondary, i+1))
			drawLineWithEffect(dst, i, y1, i+1, y2, secondaryColor)
		}
	}

	return dst
}

// DrawRadialGauge renders a donut-style gauge showing percent (0-100) as a
// clockwise sweep starting at the top, matching the mockup's SVG
// stroke-dasharray ring. w and h are the actual pixel size being rendered
// at; the ring is centered and sized to fit within min(w, h).
func DrawRadialGauge(w, h int, percent float64, ringColor, trackColor color.RGBA) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	if w <= 0 || h <= 0 {
		return dst
	}
	draw.Draw(dst, dst.Bounds(), &image.Uniform{color.Transparent}, image.Point{}, draw.Src)

	cx, cy := float64(w)/2, float64(h)/2
	outer := math.Min(float64(w), float64(h)) / 2
	if outer < 2 {
		return dst
	}
	thickness := outer * 0.21
	inner := outer - thickness

	sweep := percent / 100 * 2 * math.Pi
	if sweep > 2*math.Pi {
		sweep = 2 * math.Pi
	}
	if sweep < 0 {
		sweep = 0
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dx := float64(x) + 0.5 - cx
			dy := float64(y) + 0.5 - cy
			dist := math.Hypot(dx, dy)
			if dist < inner || dist > outer {
				continue
			}

			angle := math.Atan2(dx, -dy)
			if angle < 0 {
				angle += 2 * math.Pi
			}

			if angle <= sweep {
				dst.Set(x, y, ringColor)
			} else {
				dst.Set(x, y, trackColor)
			}
		}
	}

	return dst
}

// glowOffsets are the pixel offsets used to build a rounded halo around a
// line in dark mode — a ring at radius 1 and a fainter, wider ring at
// radius 2, rather than just 4 single-pixel taps, so the glow reads as a
// soft bloom instead of a faint cross.
var glowOffsets = []struct {
	dx, dy int
	alpha  uint8
}{
	{-1, 0, 110}, {1, 0, 110}, {0, -1, 110}, {0, 1, 110},
	{-1, -1, 80}, {1, -1, 80}, {-1, 1, 80}, {1, 1, 80},
	{-2, 0, 45}, {2, 0, 45}, {0, -2, 45}, {0, 2, 45},
}

// drawLineWithEffect draws a graph line with a themed accent underneath it:
// a soft colored glow in dark mode, or a dark drop-shadow in light mode.
// The shadow specifically matters for the yellow series, which otherwise has
// poor contrast against a white background.
func drawLineWithEffect(dst *image.RGBA, x1, y1, x2, y2 int, col color.RGBA) {
	if isDarkTheme() {
		for _, o := range glowOffsets {
			glow := color.NRGBA{R: col.R, G: col.G, B: col.B, A: o.alpha}
			drawLine(dst, x1+o.dx, y1+o.dy, x2+o.dx, y2+o.dy, glow)
		}
	} else {
		shadow := color.NRGBA{A: 160}
		drawLine(dst, x1+1, y1+2, x2+1, y2+2, shadow)
		drawLine(dst, x1+2, y1+1, x2+2, y2+1, shadow)
	}
	drawLine(dst, x1, y1, x2, y2, col)
}

// drawLine draws a line using Bresenham's line algorithm
// https://en.wikipedia.org/wiki/Bresenham%27s_line_algorithm
func drawLine(img *image.RGBA, x1, y1, x2, y2 int, col color.Color) {
	dx := abs(x2 - x1)
	sx := -1
	if x1 < x2 {
		sx = 1
	}

	dy := -abs(y2 - y1)
	sy := -1
	if y1 < y2 {
		sy = 1
	}

	err := dx + dy
	for {
		img.Set(x1, y1, col)

		if x1 == x2 && y1 == y2 {
			break
		}

		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x1 += sx
		}

		if e2 <= dx {
			err += dx
			y1 += sy
		}
	}
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// formatCoreLabel formats the core label text
func formatCoreLabel(coreNum int) string {
	return fmt.Sprintf("Core %d", coreNum)
}

// formatUtilLabel formats the utilization label text
func formatUtilLabel(util float64) string {
	return fmt.Sprintf("%.1f%%", util)
}

// formatClockLabel formats the clock speed label text
func formatClockLabel(freq float64) string {
	return fmt.Sprintf("%.0f MHz", freq)
}

// formatClockNumber formats the clock speed without a unit suffix, for use
// in the CPU list view's "Clock" column where the header already labels the
// unit.
func formatClockNumber(freq float64) string {
	return fmt.Sprintf("%.0f", freq)
}

// formatThreadsValue formats the CPU stats strip's thread count.
func formatThreadsValue(threads int) string {
	return fmt.Sprintf("%d", threads)
}

// formatPeakCoreValue formats the CPU stats strip's peak core indicator.
func formatPeakCoreValue(coreIndex int, util float64) string {
	return fmt.Sprintf("Core %d · %.1f%%", coreIndex, util)
}
