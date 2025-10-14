package main

import (
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"github.com/shirou/gopsutil/cpu"
)

func main() {
	// Load configuration
	config := DefaultConfig()
	config.ParseFlags()

	myApp := app.New()

	icon, err := fyne.LoadResourceFromPath("icon.png")
	if err != nil {
		log.Printf("Warning: Could not load icon: %v", err)
	}

	windowTitle := fmt.Sprintf("Cycles | %s", config.Version)
	myWindow := myApp.NewWindow(windowTitle)
	myWindow.SetIcon(icon)

	// Set up menu
	aboutItem := fyne.NewMenuItem("About", func() {
		dialog.ShowInformation("About Cycles", GetAppInfo(), myWindow)
	})
	
	helpMenu := fyne.NewMenu("Help", aboutItem)
	mainMenu := fyne.NewMainMenu(helpMenu)
	myWindow.SetMainMenu(mainMenu)

	// Determine the number of CPU cores
	numCores, err := cpu.Counts(config.LogicalCores)
	if err != nil {
		log.Fatalf("Error getting CPU core count: %v", err)
	}

	tiles := make([]*CoreTile, numCores)

	// Create a grid container
	grid := container.NewGridWithColumns(config.GridColumns)

	for i := 0; i < numCores; i++ {
		tiles[i] = NewCoreTile()
		grid.Add(tiles[i].GetContainer())
	}

	myWindow.SetContent(grid)

	// Update CPU info periodically
	go func() {
		for {
			UpdateCPUInfo(tiles)
			time.Sleep(config.UpdateInterval)
		}
	}()

	myWindow.ShowAndRun()
}
