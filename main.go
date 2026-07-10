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

	myApp := app.NewWithID("us.tylerc.cycles")

	// Load settings
	settings := NewSettings(myApp)

	// Command-line flags override saved settings
	if config.GridColumns != 8 {
		settings.GridColumns = config.GridColumns
	}
	if config.UpdateInterval != 2*time.Second {
		settings.UpdateInterval = config.UpdateInterval
	}
	if config.HistorySize != 30 {
		settings.HistorySize = config.HistorySize
	}

	// Apply settings to config
	settings.ApplyToConfig(config)

	// Apply theme
	ApplyTheme(myApp, settings.GetThemeVariant())

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

	settingsItem := fyne.NewMenuItem("Preferences...", func() {
		ShowSettingsDialog(settings, myWindow, func() {
			// Settings saved callback
			log.Println("Settings saved successfully")
		})
	})

	// View menu for quick theme toggle
	themeToggleItem := fyne.NewMenuItem("Toggle Theme", func() {
		if settings.Theme == "dark" {
			settings.Theme = "light"
		} else {
			settings.Theme = "dark"
		}
		settings.Save()
		ApplyTheme(myApp, settings.GetThemeVariant())
	})

	helpMenu := fyne.NewMenu("Help", aboutItem)
	fileMenu := fyne.NewMenu("File", settingsItem)
	viewMenu := fyne.NewMenu("View", themeToggleItem)
	mainMenu := fyne.NewMainMenu(fileMenu, viewMenu, helpMenu)
	myWindow.SetMainMenu(mainMenu)

	// Determine the number of CPU cores
	numCores, err := cpu.Counts(config.LogicalCores)
	if err != nil {
		log.Fatalf("Error getting CPU core count: %v", err)
	}

	cpuTiles := make([]*CoreTile, numCores)

	// Create CPU grid container
	cpuGrid := container.NewGridWithColumns(config.GridColumns)

	for i := 0; i < numCores; i++ {
		cpuTiles[i] = NewCoreTile()
		cpuGrid.Add(cpuTiles[i].GetContainer())
	}

	// Create memory tiles (overall memory + swap if available)
	memoryTiles := make([]*MemoryTile, 0)
	memoryTiles = append(memoryTiles, NewMemoryTile("System Memory"))

	// Create memory grid container
	memoryGrid := container.NewGridWithColumns(2)
	for _, tile := range memoryTiles {
		memoryGrid.Add(tile.GetContainer())
	}

	// Create tabs for CPU and Memory
	tabs := container.NewAppTabs(
		container.NewTabItem("CPU", cpuGrid),
		container.NewTabItem("Memory", memoryGrid),
	)

	myWindow.SetContent(tabs)

	// Update CPU info periodically
	go func() {
		for {
			UpdateCPUInfo(cpuTiles)
			time.Sleep(config.UpdateInterval)
		}
	}()

	// Update Memory info periodically
	go func() {
		for {
			UpdateMemoryInfo(memoryTiles, config.HistorySize)
			time.Sleep(config.UpdateInterval)
		}
	}()

	myWindow.ShowAndRun()
}
