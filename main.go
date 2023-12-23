package main

import (
    "fmt"
    "time"

    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
    "github.com/shirou/gopsutil/cpu"
)

func main() {
    myApp := app.New()
    myWindow := myApp.NewWindow("CPU Info")

    // Initialize UI components
    cpuInfo := widget.NewLabel("Getting CPU info...")
    content := container.NewVBox(cpuInfo)
    myWindow.SetContent(content)

    // Update CPU info periodically
    go func() {
        for {
            updateCPUInfo(cpuInfo)
            time.Sleep(2 * time.Second)
        }
    }()

    myWindow.ShowAndRun()
}

func updateCPUInfo(label *widget.Label) {
    percent, err := cpu.Percent(0, true)
    if err != nil {
        label.SetText(fmt.Sprintf("Error getting CPU info: %s", err))
        return
    }

    freq, err := cpu.Info()
    if err != nil {
        label.SetText(fmt.Sprintf("Error getting CPU frequency: %s", err))
        return
    }

    info := "CPU Cores:\n"
    for i, p := range percent {
        info += fmt.Sprintf("Core %d: %.2f%%, %.2f MHz\n", i, p, freq[i].Mhz)
    }

    label.SetText(info)
}

