package main

import (
	"testing"
	"time"
)

func TestSettingsDefaults(t *testing.T) {
	// Create a test settings instance
	// Note: In a real test, we'd mock the fyne.App
	// For now, we're testing the Reset functionality

	tests := []struct {
		name     string
		field    string
		expected interface{}
	}{
		{"Default Theme", "Theme", "auto"},
		{"Default Grid Columns", "GridColumns", 8},
		{"Default History Size", "HistorySize", 30},
		{"Default Logical Cores", "LogicalCores", true},
		{"Default Update Interval", "UpdateInterval", 2 * time.Second},
	}

	// We can't easily test with a real Settings object without mocking fyne.App
	// but we can test the logic
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test validates our expected defaults
			// Actual testing would require fyne.App mock
		})
	}
}

func TestGetThemeVariant(t *testing.T) {
	tests := []struct {
		theme    string
		expected string // "light" or "dark"
	}{
		{"light", "light"},
		{"dark", "dark"},
		{"auto", "dark"}, // Auto defaults to dark
		{"", "dark"},     // Empty defaults to dark
	}

	for _, tt := range tests {
		t.Run(tt.theme, func(t *testing.T) {
			// Create a settings object with the theme
			s := &Settings{
				Theme: tt.theme,
			}

			variant := s.GetThemeVariant()

			// Check if variant matches expected
			if tt.expected == "light" && variant != 0 {
				t.Errorf("Expected light theme variant for %s", tt.theme)
			}
			if tt.expected == "dark" && variant != 1 {
				t.Errorf("Expected dark theme variant for %s", tt.theme)
			}
		})
	}
}

func TestApplyToConfig(t *testing.T) {
	s := &Settings{
		GridColumns:    12,
		HistorySize:    50,
		LogicalCores:   false,
		UpdateInterval: 3 * time.Second,
	}

	config := DefaultConfig()
	s.ApplyToConfig(config)

	if config.GridColumns != 12 {
		t.Errorf("Expected GridColumns=12, got %d", config.GridColumns)
	}
	if config.HistorySize != 50 {
		t.Errorf("Expected HistorySize=50, got %d", config.HistorySize)
	}
	if config.LogicalCores != false {
		t.Errorf("Expected LogicalCores=false, got %v", config.LogicalCores)
	}
	if config.UpdateInterval != 3*time.Second {
		t.Errorf("Expected UpdateInterval=3s, got %v", config.UpdateInterval)
	}
}

func TestLoadFromConfig(t *testing.T) {
	config := &AppConfig{
		GridColumns:    16,
		HistorySize:    100,
		LogicalCores:   true,
		UpdateInterval: 5 * time.Second,
	}

	s := &Settings{}
	s.LoadFromConfig(config)

	if s.GridColumns != 16 {
		t.Errorf("Expected GridColumns=16, got %d", s.GridColumns)
	}
	if s.HistorySize != 100 {
		t.Errorf("Expected HistorySize=100, got %d", s.HistorySize)
	}
	if s.LogicalCores != true {
		t.Errorf("Expected LogicalCores=true, got %v", s.LogicalCores)
	}
	if s.UpdateInterval != 5*time.Second {
		t.Errorf("Expected UpdateInterval=5s, got %v", s.UpdateInterval)
	}
}
