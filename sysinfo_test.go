package main

import (
	"strings"
	"testing"
)

const fixtureMemInfo = `MemTotal:       65985392 kB
MemFree:         4197632 kB
MemAvailable:   50000000 kB
Buffers:          389120 kB
Cached:         47398080 kB
SwapCached:            0 kB
SwapTotal:       8388604 kB
SwapFree:        7476000 kB
`

func TestParseMemInfo(t *testing.T) {
	mem, swap, err := parseMemInfo(strings.NewReader(fixtureMemInfo))
	if err != nil {
		t.Fatalf("parseMemInfo returned error: %v", err)
	}

	if mem.Total != 65985392 {
		t.Errorf("Total = %d; want 65985392", mem.Total)
	}
	if mem.Free != 4197632 {
		t.Errorf("Free = %d; want 4197632", mem.Free)
	}
	if mem.Cached != 47398080 {
		t.Errorf("Cached = %d; want 47398080", mem.Cached)
	}
	if mem.Buffers != 389120 {
		t.Errorf("Buffers = %d; want 389120", mem.Buffers)
	}

	wantUsed := mem.Total - mem.Free - mem.Cached - mem.Buffers
	if mem.Used != wantUsed {
		t.Errorf("Used = %d; want %d", mem.Used, wantUsed)
	}

	// Breakdown components should sum back to the total, since the memory
	// breakdown bars visualize them as parts of a whole.
	if sum := mem.Used + mem.Cached + mem.Buffers + mem.Free; sum != mem.Total {
		t.Errorf("Used+Cached+Buffers+Free = %d; want Total = %d", sum, mem.Total)
	}

	if swap.Total != 8388604 {
		t.Errorf("Swap.Total = %d; want 8388604", swap.Total)
	}
	if swap.Free != 7476000 {
		t.Errorf("Swap.Free = %d; want 7476000", swap.Free)
	}
	if swap.Used != 8388604-7476000 {
		t.Errorf("Swap.Used = %d; want %d", swap.Used, 8388604-7476000)
	}
}

func TestUtilizationStatus(t *testing.T) {
	tests := []struct {
		percent  float64
		expected string
	}{
		{0, "green"},
		{59.9, "green"},
		{60, "yellow"},
		{84.9, "yellow"},
		{85, "red"},
		{100, "red"},
	}

	for _, tt := range tests {
		result := UtilizationStatus(tt.percent)
		if result != tt.expected {
			t.Errorf("UtilizationStatus(%v) = %s; want %s", tt.percent, result, tt.expected)
		}
	}
}

func TestComputeCPUAggregateStats(t *testing.T) {
	percent := []float64{10, 20, 90.5, 5}
	freqs := []float64{1746, 1746, 4850, 3861}

	stats := ComputeCPUAggregateStats(percent, freqs)

	if stats.Threads != 4 {
		t.Errorf("Threads = %d; want 4", stats.Threads)
	}
	if stats.PeakCoreIndex != 2 {
		t.Errorf("PeakCoreIndex = %d; want 2", stats.PeakCoreIndex)
	}
	if stats.PeakCoreUtil != 90.5 {
		t.Errorf("PeakCoreUtil = %v; want 90.5", stats.PeakCoreUtil)
	}
	if stats.MaxClock != 4850 {
		t.Errorf("MaxClock = %v; want 4850", stats.MaxClock)
	}

	wantAvg := (10.0 + 20.0 + 90.5 + 5.0) / 4.0
	if stats.AvgUtil != wantAvg {
		t.Errorf("AvgUtil = %v; want %v", stats.AvgUtil, wantAvg)
	}
}

func TestComputeCPUAggregateStatsEmpty(t *testing.T) {
	stats := ComputeCPUAggregateStats(nil, nil)
	if stats.Threads != 0 || stats.AvgUtil != 0 || stats.MaxClock != 0 {
		t.Errorf("expected zero-value stats for empty input, got %+v", stats)
	}
}
