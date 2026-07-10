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

func TestNewMemoryDashboard(t *testing.T) {
	dashboard := NewMemoryDashboard()

	if dashboard == nil {
		t.Fatal("NewMemoryDashboard returned nil")
	}

	if dashboard.GetContainer() == nil {
		t.Error("GetContainer returned nil")
	}
}

func TestMemoryDashboardUpdate(t *testing.T) {
	dashboard := NewMemoryDashboard()

	mem := MemoryInfo{Total: 1000, Used: 200, Cached: 700, Buffers: 50, Free: 50}
	swap := SwapInfo{Total: 100, Used: 10, Free: 90}

	// Update should not panic and should record the percentages passed in.
	dashboard.Update(mem, swap, 20.0, 10.0)

	if dashboard.gaugePercent != 20.0 {
		t.Errorf("Expected gaugePercent=20.0, got %v", dashboard.gaugePercent)
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
