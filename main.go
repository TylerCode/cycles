package main

import (
    "fmt"
    "time"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
    "github.com/shirou/gopsutil/cpu"

    "bufio"
    "os"
    "strings"
    "strconv"
)

func main() {
    myApp := app.New()
    myWindow := myApp.NewWindow("Cycles")

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
    CoreLabel    *widget.Label
    UtilLabel    *widget.Label
    ClockLabel   *widget.Label
}

func NewCoreTile() *CoreTile {
    return &CoreTile{
        CoreLabel:  widget.NewLabel("Core #"),
        UtilLabel:  widget.NewLabel("Util %"),
        ClockLabel: widget.NewLabel("Clock MHz"),
    }
}

func (t *CoreTile) GetContainer() fyne.CanvasObject {
    return container.NewVBox(t.CoreLabel, t.UtilLabel, t.ClockLabel)
}




