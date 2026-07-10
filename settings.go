package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// Settings represents user preferences for the application
type Settings struct {
	app fyne.App

	// Display settings
	Theme        string // "auto", "light", "dark"
	ViewMode     string // "tiles", "list" — CPU tab view
	HistorySize  int
	LogicalCores bool

	// Update settings
	UpdateInterval time.Duration
}

// NewSettings creates a new Settings instance with the Fyne app
func NewSettings(app fyne.App) *Settings {
	s := &Settings{
		app: app,
	}
	s.Load()
	return s
}

// Load loads settings from Fyne preferences
func (s *Settings) Load() {
	prefs := s.app.Preferences()

	// Load theme setting (default: "auto")
	s.Theme = prefs.StringWithFallback("theme", "auto")

	// Load CPU tab view mode (default: "tiles")
	s.ViewMode = prefs.StringWithFallback("view_mode", "tiles")

	// Load history size (default: 30)
	s.HistorySize = prefs.IntWithFallback("history_size", 30)

	// Load logical cores setting (default: true)
	s.LogicalCores = prefs.BoolWithFallback("logical_cores", true)

	// Load update interval in seconds (default: 2)
	intervalSecs := prefs.IntWithFallback("update_interval_secs", 2)
	s.UpdateInterval = time.Duration(intervalSecs) * time.Second
}

// Save saves settings to Fyne preferences
func (s *Settings) Save() {
	prefs := s.app.Preferences()

	prefs.SetString("theme", s.Theme)
	prefs.SetString("view_mode", s.ViewMode)
	prefs.SetInt("history_size", s.HistorySize)
	prefs.SetBool("logical_cores", s.LogicalCores)

	// Store update interval as seconds
	intervalSecs := int(s.UpdateInterval.Seconds())
	prefs.SetInt("update_interval_secs", intervalSecs)
}

// Reset resets all settings to default values
func (s *Settings) Reset() {
	s.Theme = "auto"
	s.ViewMode = "tiles"
	s.HistorySize = 30
	s.LogicalCores = true
	s.UpdateInterval = 2 * time.Second
	s.Save()
}

// ApplyToConfig applies settings to an AppConfig
func (s *Settings) ApplyToConfig(config *AppConfig) {
	config.HistorySize = s.HistorySize
	config.LogicalCores = s.LogicalCores
	config.UpdateInterval = s.UpdateInterval
}

// LoadFromConfig loads settings from an AppConfig
// This is used to populate settings from command-line flags
func (s *Settings) LoadFromConfig(config *AppConfig) {
	s.HistorySize = config.HistorySize
	s.LogicalCores = config.LogicalCores
	s.UpdateInterval = config.UpdateInterval
}

// GetThemeVariant returns the Fyne theme variant based on settings
func (s *Settings) GetThemeVariant() fyne.ThemeVariant {
	switch s.Theme {
	case "light":
		return theme.VariantLight
	case "dark":
		return theme.VariantDark
	default: // "auto"
		return theme.VariantDark // Default to dark, can be enhanced to detect system theme
	}
}
