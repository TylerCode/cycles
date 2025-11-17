package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/cpu"
)

// MemoryInfo represents memory statistics
type MemoryInfo struct {
	Total uint64
	Used  uint64
	Free  uint64
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

// UpdateCPUInfo updates the CPU information for all tiles
func UpdateCPUInfo(tiles []*CoreTile) {
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

	for i, tile := range tiles {
		if i >= len(percent) || i >= len(freqs) {
			continue
		}

		// Update labels
		tile.CoreLabel.SetText(formatCoreLabel(i))
		tile.UtilLabel.SetText(formatUtilLabel(percent[i]))
		tile.ClockLabel.SetText(formatClockLabel(freqs[i]))

		// Update utilization history
		tile.UtilHistory = append(tile.UtilHistory, percent[i])
		if len(tile.UtilHistory) > 30 { // Keep only the last 30 measurements
			tile.UtilHistory = tile.UtilHistory[1:]
		}

		// Draw graph
		DrawGraph(tile.GraphImg, tile.UtilHistory)
	}
}

// GetMemoryInfoDetailed returns detailed memory information including cached memory
func GetMemoryInfoDetailed() (MemoryInfo, uint64, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return MemoryInfo{}, 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var total, free, available, cached, buffers uint64
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
		} else if strings.HasPrefix(line, "MemAvailable:") {
			parts := strings.Fields(line)
			if len(parts) == 3 {
				available, _ = strconv.ParseUint(parts[1], 10, 64)
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
		}
	}

	if err := scanner.Err(); err != nil {
		return MemoryInfo{}, 0, err
	}

	// Calculate used memory (total - available gives a more accurate "used" value)
	used := total - available

	return MemoryInfo{
		Total: total,
		Used:  used,
		Free:  free,
	}, cached + buffers, nil
}

// UpdateMemoryInfo updates the memory information for memory tiles
func UpdateMemoryInfo(tiles []*MemoryTile, historySize int) {
	memInfo, cached, err := GetMemoryInfoDetailed()
	if err != nil {
		log.Printf("Error getting memory info: %v", err)
		return
	}

	for _, tile := range tiles {
		// Calculate usage percentage
		usagePercent := 0.0
		if memInfo.Total > 0 {
			usagePercent = float64(memInfo.Used) / float64(memInfo.Total) * 100.0
		}

		// Update labels
		tile.TotalLabel.SetText("Total: " + formatMemorySize(memInfo.Total))
		tile.UsedLabel.SetText("Used: " + formatMemorySize(memInfo.Used))
		tile.FreeLabel.SetText("Free: " + formatMemorySize(memInfo.Free))
		tile.CachedLabel.SetText("Cached: " + formatMemorySize(cached))
		tile.PercentLabel.SetText("Usage: " + formatMemoryPercent(usagePercent))

		// Update usage history
		tile.UsageHistory = append(tile.UsageHistory, usagePercent)
		if len(tile.UsageHistory) > historySize {
			tile.UsageHistory = tile.UsageHistory[1:]
		}

		// Draw graph
		DrawGraph(tile.GraphImg, tile.UsageHistory)
	}
}
