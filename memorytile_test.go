package main

import (
	"os"
	"testing"

	"fyne.io/fyne/v2/test"
)

// TestMain starts a headless Fyne test app so widget constructors that read
// theme colors (e.g. theme.BackgroundColor()) don't panic on a nil current app.
func TestMain(m *testing.M) {
	test.NewApp()
	os.Exit(m.Run())
}

func TestNewMemoryTile(t *testing.T) {
	tile := NewMemoryTile("Test Memory")

	if tile == nil {
		t.Fatal("NewMemoryTile returned nil")
	}

	if tile.TitleLabel == nil {
		t.Error("TitleLabel is nil")
	}

	if tile.TotalLabel == nil {
		t.Error("TotalLabel is nil")
	}

	if tile.UsedLabel == nil {
		t.Error("UsedLabel is nil")
	}

	if tile.FreeLabel == nil {
		t.Error("FreeLabel is nil")
	}

	if tile.CachedLabel == nil {
		t.Error("CachedLabel is nil")
	}

	if tile.PercentLabel == nil {
		t.Error("PercentLabel is nil")
	}

	if tile.GraphImg == nil {
		t.Error("GraphImg is nil")
	}

	if tile.GetContainer() == nil {
		t.Error("GetContainer returned nil")
	}

	// Check initial label text
	if tile.TitleLabel.Text != "Test Memory" {
		t.Errorf("Expected title 'Test Memory', got '%s'", tile.TitleLabel.Text)
	}
}

func TestFormatMemorySize(t *testing.T) {
	tests := []struct {
		kbytes   uint64
		expected string
	}{
		{1024 * 1024, "1.00 GB"},       // 1 GB
		{2 * 1024 * 1024, "2.00 GB"},   // 2 GB
		{512 * 1024, "512.00 MB"},      // 512 MB
		{1024, "1.00 MB"},              // 1 MB
		{8 * 1024 * 1024, "8.00 GB"},   // 8 GB
		{16 * 1024 * 1024, "16.00 GB"}, // 16 GB
	}

	for _, tt := range tests {
		result := formatMemorySize(tt.kbytes)
		if result != tt.expected {
			t.Errorf("formatMemorySize(%d) = %s; want %s", tt.kbytes, result, tt.expected)
		}
	}
}

func TestFormatMemoryPercent(t *testing.T) {
	tests := []struct {
		percent  float64
		expected string
	}{
		{50.0, "50.0%"},
		{75.5, "75.5%"},
		{0.0, "0.0%"},
		{100.0, "100.0%"},
		{33.333, "33.3%"},
	}

	for _, tt := range tests {
		result := formatMemoryPercent(tt.percent)
		if result != tt.expected {
			t.Errorf("formatMemoryPercent(%.3f) = %s; want %s", tt.percent, result, tt.expected)
		}
	}
}

func TestMemoryTileGetContainer(t *testing.T) {
	tile := NewMemoryTile("Test")
	container := tile.GetContainer()

	if container == nil {
		t.Error("GetContainer() returned nil")
	}
}
