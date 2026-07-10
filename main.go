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

	helpMenu := fyne.NewMenu("Help", aboutItem)
	fileMenu := fyne.NewMenu("File", settingsItem)
	mainMenu := fyne.NewMainMenu(fileMenu, helpMenu)
	myWindow.SetMainMenu(mainMenu)

	// Determine the number of CPU cores
	numCores, err := cpu.Counts(config.LogicalCores)
	if err != nil {
		log.Fatalf("Error getting CPU core count: %v", err)
	}

	cpuTab := NewCPUTab(numCores, settings.ViewMode, func(view string) {
		settings.ViewMode = view
		settings.Save()
	})

	memoryDashboard := NewMemoryDashboard()

	// Create tabs for CPU and Memory
	tabs := container.NewAppTabs(
		container.NewTabItem("CPU", cpuTab.Content),
		container.NewTabItem("Memory", memoryDashboard.GetContainer()),
	)

	myWindow.SetContent(tabs)
	myWindow.Resize(fyne.NewSize(1280, 800))

	// canvas.Text/Rectangle primitives (used throughout the CPU/Memory tabs
	// for custom-colored drawing) only read theme colors once, at
	// construction — unlike built-in widgets, they don't repaint themselves
	// when the app theme changes live. Listen for theme changes and push a
	// refresh through both tabs so switching themes doesn't require a
	// restart.
	themeChanges := make(chan fyne.Settings)
	myApp.Settings().AddChangeListener(themeChanges)
	go func() {
		for range themeChanges {
			cpuTab.RefreshTheme()
			memoryDashboard.RefreshTheme()
		}
	}()

	// Update CPU info periodically
	go func() {
		for {
			UpdateCPUInfo(cpuTab.Cores, config.HistorySize, cpuTab.UpdateStats)
			time.Sleep(config.UpdateInterval)
		}
	}()

	// Update Memory info periodically
	go func() {
		for {
			UpdateMemoryInfo(memoryDashboard, config.HistorySize)
			time.Sleep(config.UpdateInterval)
		}
	}()

	myWindow.ShowAndRun()
}
