package main

import (
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
		{0, "Core #0"},
		{1, "Core #1"},
		{15, "Core #15"},
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
		{0.0, "Util: 0.00%"},
		{50.5, "Util: 50.50%"},
		{100.0, "Util: 100.00%"},
		{99.99, "Util: 99.99%"},
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
		{1000.0, "Clock: 1000.00 MHz"},
		{2500.5, "Clock: 2500.50 MHz"},
		{3600.99, "Clock: 3600.99 MHz"},
	}

	for _, tt := range tests {
		result := formatClockLabel(tt.freq)
		if result != tt.expected {
			t.Errorf("formatClockLabel(%f) = %s; want %s", tt.freq, result, tt.expected)
		}
	}
}
