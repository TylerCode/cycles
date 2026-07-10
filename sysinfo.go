package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/cpu"
)

// MemoryInfo represents memory statistics
type MemoryInfo struct {
	Total   uint64
	Used    uint64
	Free    uint64
	Cached  uint64
	Buffers uint64
}

// SwapInfo represents swap statistics (in kB, matching /proc/meminfo units)
type SwapInfo struct {
	Total uint64
	Used  uint64
	Free  uint64
}

// CPUAggregateStats summarizes per-core stats into the CPU tab's stats strip
type CPUAggregateStats struct {
	Threads       int
	AvgUtil       float64
	PeakCoreIndex int
	PeakCoreUtil  float64
	MaxClock      float64
}

// ComputeCPUAggregateStats computes aggregate stats from per-core utilization
// percentages and clock frequencies (as returned by cpu.Percent and
// GetCPUFrequencies).
func ComputeCPUAggregateStats(percent, freqs []float64) CPUAggregateStats {
	stats := CPUAggregateStats{Threads: len(percent)}
	if len(percent) == 0 {
		return stats
	}

	var total float64
	for i, p := range percent {
		total += p
		if p > stats.PeakCoreUtil {
			stats.PeakCoreUtil = p
			stats.PeakCoreIndex = i
		}
	}
	stats.AvgUtil = total / float64(len(percent))

	for _, f := range freqs {
		if f > stats.MaxClock {
			stats.MaxClock = f
		}
	}

	return stats
}

// GetCPUFrequencies reads CPU frequencies from /proc/cpuinfo
func GetCPUFrequencies() ([]float64, error) {
	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var freqs []float64
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu MHz") {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				freq, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
				if err == nil {
					freqs = append(freqs, freq)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return freqs, nil
}

// GetMemoryInfo returns a MemoryInfo struct with the total, used, and free memory
func GetMemoryInfo() (MemoryInfo, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return MemoryInfo{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var total, free uint64
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			parts := strings.Fields(line)
			if len(parts) == 3 {
				total, _ = strconv.ParseUint(parts[1], 10, 64)
			}
		} else if strings.HasPrefix(line, "MemFree:") {
			parts := strings.Fields(line)
			if len(parts) == 3 {
				free, _ = strconv.ParseUint(parts[1], 10, 64)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return MemoryInfo{}, err
	}

	used := total - free
	return MemoryInfo{
		Total: total,
		Used:  used,
		Free:  free,
	}, nil
}

// UpdateCPUInfo updates every core's state (feeding both the tile and list
// representations) and the CPU tab's aggregate stats strip.
func UpdateCPUInfo(cores []*CoreState, historySize int, onStats func(CPUAggregateStats)) {
	percent, err := cpu.Percent(0, true)
	if err != nil {
		log.Printf("Error getting CPU percent: %v", err)
		return
	}

	freqs, err := GetCPUFrequencies()
	if err != nil {
		log.Printf("Error getting CPU frequencies: %v", err)
		return
	}

	for i, core := range cores {
		if i >= len(percent) || i >= len(freqs) {
			continue
		}

		core.Update(percent[i], freqs[i], historySize)
	}

	if onStats != nil {
		onStats(ComputeCPUAggregateStats(percent, freqs))
	}
}

// GetMemoryInfoDetailed returns detailed memory information (including
// cached/buffers, broken out separately) and swap information, both read
// from /proc/meminfo.
//
// Used is computed as total-free-cached-buffers (rather than total-available)
// so that Used+Cached+Buffers+Free sums to Total — required for the memory
// breakdown bars, which visualize those four components as parts of a whole.
func GetMemoryInfoDetailed() (MemoryInfo, SwapInfo, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return MemoryInfo{}, SwapInfo{}, err
	}
	defer file.Close()

	return parseMemInfo(file)
}

// parseMemInfo parses /proc/meminfo-formatted content into memory and swap
// stats. Split out from GetMemoryInfoDetailed so the parsing logic can be
// tested against a fixture without touching the filesystem.
func parseMemInfo(r io.Reader) (MemoryInfo, SwapInfo, error) {
	scanner := bufio.NewScanner(r)
	var total, free, cached, buffers, swapTotal, swapFree uint64
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			parts := strings.Fields(line)
			if len(parts) == 3 {
				total, _ = strconv.ParseUint(parts[1], 10, 64)
			}
		} else if strings.HasPrefix(line, "MemFree:") {
			parts := strings.Fields(line)
			if len(parts) == 3 {
				free, _ = strconv.ParseUint(parts[1], 10, 64)
			}
		} else if strings.HasPrefix(line, "Cached:") {
			parts := strings.Fields(line)
			if len(parts) == 3 {
				cached, _ = strconv.ParseUint(parts[1], 10, 64)
			}
		} else if strings.HasPrefix(line, "Buffers:") {
			parts := strings.Fields(line)
			if len(parts) == 3 {
				buffers, _ = strconv.ParseUint(parts[1], 10, 64)
			}
		} else if strings.HasPrefix(line, "SwapTotal:") {
			parts := strings.Fields(line)
			if len(parts) == 3 {
				swapTotal, _ = strconv.ParseUint(parts[1], 10, 64)
			}
		} else if strings.HasPrefix(line, "SwapFree:") {
			parts := strings.Fields(line)
			if len(parts) == 3 {
				swapFree, _ = strconv.ParseUint(parts[1], 10, 64)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return MemoryInfo{}, SwapInfo{}, err
	}

	used := total - free - cached - buffers

	memInfo := MemoryInfo{
		Total:   total,
		Used:    used,
		Free:    free,
		Cached:  cached,
		Buffers: buffers,
	}
	swapInfo := SwapInfo{
		Total: swapTotal,
		Free:  swapFree,
		Used:  swapTotal - swapFree,
	}

	return memInfo, swapInfo, nil
}

// UpdateMemoryInfo updates the memory dashboard with the latest memory and
// swap statistics.
func UpdateMemoryInfo(dashboard *MemoryDashboard, historySize int) {
	memInfo, swapInfo, err := GetMemoryInfoDetailed()
	if err != nil {
		log.Printf("Error getting memory info: %v", err)
		return
	}

	usagePercent := 0.0
	if memInfo.Total > 0 {
		usagePercent = float64(memInfo.Used) / float64(memInfo.Total) * 100.0
	}

	swapPercent := 0.0
	if swapInfo.Total > 0 {
		swapPercent = float64(swapInfo.Used) / float64(swapInfo.Total) * 100.0
	}

	dashboard.UsageHistory = append(dashboard.UsageHistory, usagePercent)
	if len(dashboard.UsageHistory) > historySize {
		dashboard.UsageHistory = dashboard.UsageHistory[1:]
	}

	dashboard.SwapHistory = append(dashboard.SwapHistory, swapPercent)
	if len(dashboard.SwapHistory) > historySize {
		dashboard.SwapHistory = dashboard.SwapHistory[1:]
	}

	dashboard.Update(memInfo, swapInfo, usagePercent, swapPercent)
}
