package main

import (
	"flag"
	"time"
)

// AppConfig holds the application configuration
type AppConfig struct {
	Version        string
	GridColumns    int
	UpdateInterval time.Duration
	HistorySize    int
	LogicalCores   bool
}

// DefaultConfig returns the default configuration
func DefaultConfig() *AppConfig {
	return &AppConfig{
		Version:        "0.4.0",
		GridColumns:    8,
		UpdateInterval: 2 * time.Second,
		HistorySize:    30,
		LogicalCores:   true,
	}
}

// ParseFlags parses command-line flags and updates the configuration
func (c *AppConfig) ParseFlags() {
	flag.IntVar(&c.GridColumns, "columns", c.GridColumns, "Number of columns in the grid layout")
	flag.DurationVar(&c.UpdateInterval, "interval", c.UpdateInterval, "Update interval for CPU monitoring")
	flag.IntVar(&c.HistorySize, "history", c.HistorySize, "Number of historical data points to keep")
	flag.BoolVar(&c.LogicalCores, "logical", c.LogicalCores, "Show logical cores (true) or physical cores (false)")
	flag.Parse()
}
