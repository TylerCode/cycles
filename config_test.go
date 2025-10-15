package main

import (
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.Version != "0.4.1" {
		t.Errorf("Expected version 0.4.1, got %s", config.Version)
	}

	if config.GridColumns != 4 {
		t.Errorf("Expected 4 grid columns, got %d", config.GridColumns)
	}

	if config.UpdateInterval != 2*time.Second {
		t.Errorf("Expected update interval of 2 seconds, got %v", config.UpdateInterval)
	}

	if config.HistorySize != 30 {
		t.Errorf("Expected history size of 30, got %d", config.HistorySize)
	}

	if !config.LogicalCores {
		t.Error("Expected LogicalCores to be true")
	}
}
