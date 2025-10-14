package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

// DrawGraph draws a utilization graph on the provided canvas image
func DrawGraph(img *canvas.Image, data []float64) {
	const width, height = 120, 50 // Graph dimensions

	// Create a new image for the graph
	rect := image.Rect(0, 0, width, height)
	dst := image.NewRGBA(rect)

	// Set background color
	backgroundColor := theme.BackgroundColor()
	draw.Draw(dst, dst.Bounds(), &image.Uniform{backgroundColor}, image.ZP, draw.Src)

	// Draw the box around the graph
	borderColor := color.RGBA{128, 128, 128, 255}
	drawLine(dst, 0, 0, width-1, 0, borderColor)               // Top border
	drawLine(dst, 0, height-1, width-1, height-1, borderColor) // Bottom border
	drawLine(dst, 0, 0, 0, height-1, borderColor)              // Left border
	drawLine(dst, width-1, 0, width-1, height-1, borderColor)  // Right border

	// Check if there's data to draw
	if len(data) < 2 {
		img.Image = dst
		img.Refresh()
		return
	}

	// Calculate the x-axis step
	step := width / (len(data) - 1)

	// Draw the graph lines
	for i := 0; i < len(data)-1; i++ {
		x1 := i * step
		y1 := height - int(data[i]/100*float64(height))
		x2 := (i + 1) * step
		y2 := height - int(data[i+1]/100*float64(height))

		// Determine line color based on utilization
		lineColor := GetGraphLineColor("green") // Green for utilization under 75%
		if data[i] >= 75 || data[i+1] >= 75 {
			lineColor = GetGraphLineColor("red") // Red for utilization 75% or above
		}

		drawLine(dst, x1, y1, x2, y2, lineColor)
	}

	img.Image = dst
	img.Refresh()
}

// drawLine draws a line using Bresenham's line algorithm
// https://en.wikipedia.org/wiki/Bresenham%27s_line_algorithm
func drawLine(img *image.RGBA, x1, y1, x2, y2 int, col color.RGBA) {
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
	return fmt.Sprintf("Core #%d", coreNum)
}

// formatUtilLabel formats the utilization label text
func formatUtilLabel(util float64) string {
	return fmt.Sprintf("Util: %.2f%%", util)
}

// formatClockLabel formats the clock speed label text
func formatClockLabel(freq float64) string {
	return fmt.Sprintf("Clock: %.2f MHz", freq)
}
