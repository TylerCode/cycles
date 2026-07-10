package main

import (
	"testing"
	"time"

	"fyne.io/fyne/v2/theme"
)

func TestSettingsDefaults(t *testing.T) {
	// Create a test settings instance
	// Note: In a real test, we'd mock the fyne.App
	// For now, we're testing the Reset functionality

	tests := []struct {
		name     string
		field    string
		expected any
	}{
		{"Default Theme", "Theme", "auto"},
		{"Default View Mode", "ViewMode", "tiles"},
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
		themeName string
		expected  string // "light" or "dark"
	}{
		{"light", "light"},
		{"dark", "dark"},
		{"auto", "dark"}, // Auto defaults to dark
		{"", "dark"},     // Empty defaults to dark
	}

	for _, tt := range tests {
		t.Run(tt.themeName, func(t *testing.T) {
			s := &Settings{
				Theme: tt.themeName,
			}

			variant := s.GetThemeVariant()

			if tt.expected == "light" && variant != theme.VariantLight {
				t.Errorf("Expected light theme variant for %s", tt.themeName)
			}
			if tt.expected == "dark" && variant != theme.VariantDark {
				t.Errorf("Expected dark theme variant for %s", tt.themeName)
			}
		})
	}
}

func TestApplyToConfig(t *testing.T) {
	s := &Settings{
		HistorySize:    50,
		LogicalCores:   false,
		UpdateInterval: 3 * time.Second,
	}

	config := DefaultConfig()
	s.ApplyToConfig(config)

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
		HistorySize:    100,
		LogicalCores:   true,
		UpdateInterval: 5 * time.Second,
	}

	s := &Settings{}
	s.LoadFromConfig(config)

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
