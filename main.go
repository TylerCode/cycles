package main

import (
    "fmt"
    "time"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
    "fyne.io/fyne/v2/canvas"
    "fyne.io/fyne/v2/theme"
    "github.com/shirou/gopsutil/cpu"

    "bufio"
    "os"
    "strings"
    "strconv"
)

func main() {
    myApp := app.New()

    icon, err := fyne.LoadResourceFromPath("icon.png")
    if err != nil {
        //log.Fatal("Could not load icon:", err)
    }

    myWindow := myApp.NewWindow("Cycles")
    myWindow.SetIcon(icon)

    // Determine the number of CPU cores
    numCores, _ := cpu.Counts(true) // True because we want logical and physical
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
    }
}


type CoreTile struct {
    CoreLabel  *widget.Label
    UtilLabel  *widget.Label
    ClockLabel *widget.Label
    container  *fyne.Container
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

    // Use a container to overlay the labels on the background
    container := container.NewMax(bg, container.NewVBox(coreLabel, utilLabel, clockLabel))

    return &CoreTile{
        CoreLabel:  coreLabel,
        UtilLabel:  utilLabel,
        ClockLabel: clockLabel,
        container:  container,
    }
}

func (t *CoreTile) GetContainer() fyne.CanvasObject {
    return t.container
}
