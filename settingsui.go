package main

import (
	"fmt"
	"image/color"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// settingsSectionHeader renders a small section label matching the style
// used for the CPU stats strip and Memory chart legend elsewhere in the app
// (bold, slightly larger, theme foreground), rather than the plain form
// labels a bare widget.Form would otherwise produce.
func settingsSectionHeader(text string) fyne.CanvasObject {
	label := canvas.NewText(text, theme.ForegroundColor())
	label.TextStyle = fyne.TextStyle{Bold: true}
	label.TextSize = 13
	return container.NewPadded(label)
}

// settingsRow lays a fixed-width label alongside a control, matching the
// label-left/control-right rows used throughout the CPU list view and
// Memory breakdown panel.
func settingsRow(label string, control fyne.CanvasObject) fyne.CanvasObject {
	labelText := canvas.NewText(label, theme.ForegroundColor())
	labelText.TextSize = 13
	labelBox := container.NewGridWrap(fyne.NewSize(170, 30), labelText)
	return container.NewBorder(nil, nil, labelBox, nil, container.NewVBox(container.NewPadded(control)))
}

// ShowSettingsDialog displays the settings configuration dialog
func ShowSettingsDialog(settings *Settings, window fyne.Window, onSave func()) {
	// Theme selection
	themeSelect := widget.NewSelect([]string{"Auto", "Light", "Dark"}, func(value string) {
		// Callback handled on save
	})
	switch settings.Theme {
	case "light":
		themeSelect.SetSelected("Light")
	case "dark":
		themeSelect.SetSelected("Dark")
	default:
		themeSelect.SetSelected("Auto")
	}

	// History size slider, with its current value inline in the label
	// (rather than as a separate row) so it reads as one control.
	historyLabelText := canvas.NewText(fmt.Sprintf("History Size: %d", settings.HistorySize), theme.ForegroundColor())
	historyLabelText.TextSize = 13
	historySlider := widget.NewSlider(10, 100)
	historySlider.SetValue(float64(settings.HistorySize))
	historySlider.OnChanged = func(value float64) {
		historyLabelText.Text = fmt.Sprintf("History Size: %d", int(value))
		historyLabelText.Refresh()
	}

	// Update interval entry
	intervalEntry := widget.NewEntry()
	intervalEntry.SetText(fmt.Sprintf("%.0f", settings.UpdateInterval.Seconds()))
	intervalEntry.SetPlaceHolder("Seconds")

	// Logical cores checkbox
	logicalCoresCheck := widget.NewCheck("Show Logical Cores", nil)
	logicalCoresCheck.SetChecked(settings.LogicalCores)

	divider := func() fyne.CanvasObject {
		d := canvas.NewRectangle(theme.ShadowColor())
		d.SetMinSize(fyne.NewSize(0, 1))
		return d
	}

	body := container.NewVBox(
		settingsSectionHeader("Appearance"),
		settingsRow("Theme", themeSelect),

		container.NewPadded(divider()),

		settingsSectionHeader("Monitoring"),
		settingsRow("History", container.NewVBox(historyLabelText, historySlider)),
		settingsRow("Update Interval (s)", intervalEntry),
		settingsRow("CPU Cores", logicalCoresCheck),
	)

	// settingsDialog is built without Fyne's own Confirm/Dismiss button row
	// (dialog.NewCustomConfirm) since that stacks a second button row below
	// the content - here Reset/Cancel/Save are laid out together as a single
	// footer row instead, which is what a normal dialog footer looks like.
	var settingsDialog *dialog.CustomDialog

	resetButton := widget.NewButton("Reset to Defaults", func() {
		confirm := dialog.NewConfirm("Reset Settings",
			"Are you sure you want to reset all settings to default values?",
			func(confirmed bool) {
				if confirmed {
					settings.Reset()
					// Reload UI with default values
					themeSelect.SetSelected("Auto")
					historySlider.SetValue(30)
					historyLabelText.Text = "History Size: 30"
					historyLabelText.Refresh()
					intervalEntry.SetText("2")
					logicalCoresCheck.SetChecked(true)

					if onSave != nil {
						onSave()
					}
				}
			}, window)
		confirm.Show()
	})

	applyAndSave := func() {
		switch themeSelect.Selected {
		case "Light":
			settings.Theme = "light"
		case "Dark":
			settings.Theme = "dark"
		default:
			settings.Theme = "auto"
		}

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

		settingsDialog.Hide()

		// Show restart notification
		dialog.ShowInformation("Settings Saved",
			"The CPU Cores setting will take effect after restarting the application.",
			window)

		if onSave != nil {
			onSave()
		}
	}

	cancelButton := widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() {
		settingsDialog.Hide()
	})
	saveButton := widget.NewButtonWithIcon("Save", theme.ConfirmIcon(), applyAndSave)
	saveButton.Importance = widget.HighImportance

	footer := container.NewBorder(nil, nil, resetButton, container.NewHBox(cancelButton, saveButton))

	content := container.NewBorder(
		nil,
		container.NewPadded(footer),
		nil,
		nil,
		container.NewPadded(body),
	)

	settingsDialog = dialog.NewCustomWithoutButtons("Settings", content, window)
	settingsDialog.Resize(fyne.NewSize(460, 480))
	settingsDialog.Show()
}

// CustomTheme is a theme that wraps the default theme with custom variant
type CustomTheme struct {
	variant fyne.ThemeVariant
}

func (t *CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, t.variant)
}

func (t *CustomTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *CustomTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *CustomTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
