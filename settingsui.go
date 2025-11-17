package main

import (
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// ShowSettingsDialog displays the settings configuration dialog
func ShowSettingsDialog(settings *Settings, window fyne.Window, onSave func()) {
	// Theme selection
	themeSelect := widget.NewSelect([]string{"Auto", "Light", "Dark"}, func(value string) {
		// Callback handled on save
	})
	// Set current theme
	switch settings.Theme {
	case "light":
		themeSelect.SetSelected("Light")
	case "dark":
		themeSelect.SetSelected("Dark")
	default:
		themeSelect.SetSelected("Auto")
	}

	// Grid columns slider
	columnsLabel := widget.NewLabel(fmt.Sprintf("Grid Columns: %d", settings.GridColumns))
	columnsSlider := widget.NewSlider(1, 16)
	columnsSlider.SetValue(float64(settings.GridColumns))
	columnsSlider.OnChanged = func(value float64) {
		columnsLabel.SetText(fmt.Sprintf("Grid Columns: %d", int(value)))
	}

	// History size slider
	historyLabel := widget.NewLabel(fmt.Sprintf("History Size: %d", settings.HistorySize))
	historySlider := widget.NewSlider(10, 100)
	historySlider.SetValue(float64(settings.HistorySize))
	historySlider.OnChanged = func(value float64) {
		historyLabel.SetText(fmt.Sprintf("History Size: %d", int(value)))
	}

	// Update interval entry
	intervalEntry := widget.NewEntry()
	intervalEntry.SetText(fmt.Sprintf("%.0f", settings.UpdateInterval.Seconds()))
	intervalEntry.SetPlaceHolder("Update interval in seconds")

	// Logical cores checkbox
	logicalCoresCheck := widget.NewCheck("Show Logical Cores", nil)
	logicalCoresCheck.SetChecked(settings.LogicalCores)

	// Create form
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Theme", Widget: themeSelect},
			{Text: "", Widget: columnsLabel},
			{Text: "", Widget: columnsSlider},
			{Text: "", Widget: historyLabel},
			{Text: "", Widget: historySlider},
			{Text: "Update Interval (seconds)", Widget: intervalEntry},
			{Text: "CPU Cores", Widget: logicalCoresCheck},
		},
	}

	// Reset button
	resetButton := widget.NewButton("Reset to Defaults", func() {
		confirm := dialog.NewConfirm("Reset Settings",
			"Are you sure you want to reset all settings to default values?",
			func(confirmed bool) {
				if confirmed {
					settings.Reset()
					// Reload UI with default values
					themeSelect.SetSelected("Auto")
					columnsSlider.SetValue(8)
					columnsLabel.SetText("Grid Columns: 8")
					historySlider.SetValue(30)
					historyLabel.SetText("History Size: 30")
					intervalEntry.SetText("2")
					logicalCoresCheck.SetChecked(true)

					if onSave != nil {
						onSave()
					}
				}
			}, window)
		confirm.Show()
	})

	// Create dialog content with form and reset button
	content := container.NewBorder(
		nil,
		container.NewHBox(resetButton),
		nil,
		nil,
		form,
	)

	// Create dialog
	settingsDialog := dialog.NewCustomConfirm("Settings", "Save", "Cancel", content,
		func(save bool) {
			if save {
				// Apply settings from UI
				switch themeSelect.Selected {
				case "Light":
					settings.Theme = "light"
				case "Dark":
					settings.Theme = "dark"
				default:
					settings.Theme = "auto"
				}

				settings.GridColumns = int(columnsSlider.Value)
				settings.HistorySize = int(historySlider.Value)
				settings.LogicalCores = logicalCoresCheck.Checked

				// Parse update interval
				if intervalSecs, err := strconv.ParseFloat(intervalEntry.Text, 64); err == nil && intervalSecs > 0 {
					settings.UpdateInterval = time.Duration(intervalSecs) * time.Second
				}

				// Save settings
				settings.Save()

				// Apply theme immediately
				fyne.CurrentApp().Settings().SetTheme(&CustomTheme{variant: settings.GetThemeVariant()})

				// Show restart notification
				dialog.ShowInformation("Settings Saved",
					"Some settings (like Grid Columns and CPU Cores) will take effect after restarting the application.",
					window)

				if onSave != nil {
					onSave()
				}
			}
		}, window)

	settingsDialog.Resize(fyne.NewSize(400, 500))
	settingsDialog.Show()
}

// CustomTheme is a theme that wraps the default theme with custom variant
type CustomTheme struct {
	variant fyne.ThemeVariant
}

func (t *CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) fyne.Color {
	return fyne.CurrentApp().Settings().Theme().Color(name, t.variant)
}

func (t *CustomTheme) Font(style fyne.TextStyle) fyne.Resource {
	return fyne.CurrentApp().Settings().Theme().Font(style)
}

func (t *CustomTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return fyne.CurrentApp().Settings().Theme().Icon(name)
}

func (t *CustomTheme) Size(name fyne.ThemeSizeName) float32 {
	return fyne.CurrentApp().Settings().Theme().Size(name)
}
