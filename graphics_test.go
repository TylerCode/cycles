package main

import (
	"image"
	"testing"
)

func TestAbs(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{5, 5},
		{-5, 5},
		{0, 0},
		{-100, 100},
		{100, 100},
	}

	for _, tt := range tests {
		result := abs(tt.input)
		if result != tt.expected {
			t.Errorf("abs(%d) = %d; want %d", tt.input, result, tt.expected)
		}
	}
}

func TestFormatCoreLabel(t *testing.T) {
	tests := []struct {
		coreNum  int
		expected string
	}{
		{0, "Core 0"},
		{1, "Core 1"},
		{15, "Core 15"},
	}

	for _, tt := range tests {
		result := formatCoreLabel(tt.coreNum)
		if result != tt.expected {
			t.Errorf("formatCoreLabel(%d) = %s; want %s", tt.coreNum, result, tt.expected)
		}
	}
}

func TestFormatUtilLabel(t *testing.T) {
	tests := []struct {
		util     float64
		expected string
	}{
		{0.0, "0.0%"},
		{50.5, "50.5%"},
		{100.0, "100.0%"},
		{99.99, "100.0%"},
	}

	for _, tt := range tests {
		result := formatUtilLabel(tt.util)
		if result != tt.expected {
			t.Errorf("formatUtilLabel(%f) = %s; want %s", tt.util, result, tt.expected)
		}
	}
}

func TestFormatClockLabel(t *testing.T) {
	tests := []struct {
		freq     float64
		expected string
	}{
		{1000.0, "1000 MHz"},
		{2500.5, "2500 MHz"},
		{3600.99, "3601 MHz"},
	}

	for _, tt := range tests {
		result := formatClockLabel(tt.freq)
		if result != tt.expected {
			t.Errorf("formatClockLabel(%f) = %s; want %s", tt.freq, result, tt.expected)
		}
	}
}

func TestFormatClockNumber(t *testing.T) {
	if got := formatClockNumber(1746); got != "1746" {
		t.Errorf("formatClockNumber(1746) = %s; want 1746", got)
	}
}

func TestFormatThreadsValue(t *testing.T) {
	if got := formatThreadsValue(32); got != "32" {
		t.Errorf("formatThreadsValue(32) = %s; want 32", got)
	}
}

func TestFormatPeakCoreValue(t *testing.T) {
	if got := formatPeakCoreValue(16); got != "Core 16" {
		t.Errorf("formatPeakCoreValue(16) = %s; want %q", got, "Core 16")
	}
}

func TestDrawSparklineSizesToRequestedDimensions(t *testing.T) {
	img := DrawSparkline(40, 20, []float64{10, 50, 90, 20})
	bounds := img.Bounds()
	if bounds.Dx() != 40 || bounds.Dy() != 20 {
		t.Errorf("DrawSparkline size = %dx%d; want 40x20", bounds.Dx(), bounds.Dy())
	}
}

func TestDrawSparklineHandlesShortData(t *testing.T) {
	img := DrawSparkline(40, 20, []float64{10})
	if img.Bounds().Dx() != 40 {
		t.Errorf("DrawSparkline should still return correctly sized image for short data")
	}
}

func TestDrawAreaChartSizesToRequestedDimensions(t *testing.T) {
	img := DrawAreaChart(80, 40, []float64{10, 40, 70}, []float64{5, 10, 15})
	bounds := img.Bounds()
	if bounds.Dx() != 80 || bounds.Dy() != 40 {
		t.Errorf("DrawAreaChart size = %dx%d; want 80x40", bounds.Dx(), bounds.Dy())
	}
}

func TestDrawRadialGaugeSizesToRequestedDimensions(t *testing.T) {
	img := DrawRadialGauge(100, 100, 42, GetSeriesColor("blue"), GetSeriesColor("gray"))
	bounds := img.Bounds()
	if bounds.Dx() != 100 || bounds.Dy() != 100 {
		t.Errorf("DrawRadialGauge size = %dx%d; want 100x100", bounds.Dx(), bounds.Dy())
	}
}

func TestDrawRadialGaugeZeroSize(t *testing.T) {
	img := DrawRadialGauge(0, 0, 50, GetSeriesColor("blue"), GetSeriesColor("gray"))
	if _, ok := img.(*image.RGBA); !ok {
		t.Fatal("expected an RGBA image even for zero size")
	}
}
