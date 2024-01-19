package main

import (
	"fmt"
	"image/color"
	"image/draw"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/shirou/gopsutil/cpu"

	"bufio"
	"os"
	"strconv"
	"strings"

	"image"
)

type CoreTile struct {
	CoreLabel   *widget.Label
	UtilLabel   *widget.Label
	ClockLabel  *widget.Label
	container   *fyne.Container
	UtilHistory []float64 // Slice to store utilization history
	GraphImg    *canvas.Image
}

func main() {
	myApp := app.New()

	icon, err := fyne.LoadResourceFromPath("icon.png")
	if err != nil {
		//log.Fatal("Could not load icon:", err)
	}

	myWindow := myApp.NewWindow("Cycles | 0.3.4")
	myWindow.SetIcon(icon)

	// Determine the number of CPU cores
	numCores, _ := cpu.Counts(true) // True because we want logical cores, this should be a toggle in settings
	tiles := make([]*CoreTile, numCores)

	// Create a grid container
	grid := container.NewGridWithColumns(4) // Adjust number of columns as needed

	for i := 0; i < numCores; i++ {
		tiles[i] = NewCoreTile()
		grid.Add(tiles[i].GetContainer())
	}

	myWindow.SetContent(grid)

	// Update CPU info periodically
	go func() {
		for {
			updateCPUInfo(tiles)
			time.Sleep(2 * time.Second)
		}
	}()

	myWindow.ShowAndRun()
}

func updateCPUInfo(tiles []*CoreTile) {
	percent, err := cpu.Percent(0, true)
	if err != nil {
		return
	}

	// Read current frequency from /proc/cpuinfo
	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var freqs []float64
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu MHz") {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				freq, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
				if err == nil {
					freqs = append(freqs, freq)
				}
			}
		}
	}

	for i, tile := range tiles {
		// Assuming freqs and percent are fetched as before
		tile.CoreLabel.SetText(fmt.Sprintf("Core #%d", i))
		tile.UtilLabel.SetText(fmt.Sprintf("Util: %.2f%%", percent[i]))
		tile.ClockLabel.SetText(fmt.Sprintf("Clock: %.2f MHz", freqs[i]))

		// Update utilization history
		tile.UtilHistory = append(tile.UtilHistory, percent[i])
		if len(tile.UtilHistory) > 30 { // Keep only the last minute of measurements
			tile.UtilHistory = tile.UtilHistory[1:]
		}

		// Draw graph
		drawGraph(tile.GraphImg, tile.UtilHistory)
	}
}

func drawGraph(img *canvas.Image, data []float64) {
	const width, height = 120, 50 // Graph dimensions

	// Create a new image for the graph
	rect := image.Rect(0, 0, width, height)
	dst := image.NewRGBA(rect)

	// Set background color
	backgroundColor := theme.BackgroundColor()
	draw.Draw(dst, dst.Bounds(), &image.Uniform{backgroundColor}, image.ZP, draw.Src)

	// Draw the box around the graph
	borderColor := color.RGBA{128, 128, 128, 255}              // Grey color
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

		drawLine(dst, x1, y1, x2, y2, lineColor) // Perform type assertion to convert lineColor to color.RGBA

	}

	img.Image = dst
	img.Refresh()
}

// Bresenham's line algorithm
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

// simple abs so that I don't need a whole math import for one function
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func NewCoreTile() *CoreTile {
	coreLabel := widget.NewLabel("Core #")
	utilLabel := widget.NewLabel("Util %")
	clockLabel := widget.NewLabel("Clock MHz")

	// Create a background rectangle with rounded corners
	bg := canvas.NewRectangle(theme.BackgroundColor())
	bg.SetMinSize(fyne.NewSize(100, 100)) // Set the size as needed
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

func (t *CoreTile) GetContainer() fyne.CanvasObject {
	return t.container
}

// Theme table, mostly for graphs, need to move out of this into a real theme
var (
	GreenLight  = color.RGBA{R: 26, G: 155, B: 12, A: 255}  // Light theme green
	YellowLight = color.RGBA{R: 190, G: 161, B: 14, A: 255} // Light theme yellow
	RedLight    = color.RGBA{R: 186, G: 14, B: 23, A: 255}  // Light theme red

	GreenDark  = color.RGBA{R: 21, G: 222, B: 0, A: 255}  // Dark theme green
	YellowDark = color.RGBA{R: 255, G: 214, B: 0, A: 255} // Dark theme yellow
	RedDark    = color.RGBA{R: 252, G: 0, B: 13, A: 255}  // Dark theme red
)

// This doesn't actually work so far as I can tell because light theme is gone
func GetGraphLineColor(status string) color.RGBA {
	currentTheme := fyne.CurrentApp().Settings().Theme()
	isDark := true

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

	return GreenLight
}
