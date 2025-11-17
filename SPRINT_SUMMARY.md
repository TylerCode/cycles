# Cycles Overhaul Sprint Summary
## Quick Reference Guide

**Current Version:** 0.4.1 → **Target Version:** 1.2.0

---

## Sprint Quick Reference

### Sprint 1: Memory Monitoring (v0.5.0) - 2-3 days
**Focus:** Add memory monitoring UI
- Create `memorytile.go` and `memorytile_test.go`
- Add Memory tab with graphs
- Enhance memory functions in `sysinfo.go`
- Display RAM, Swap, Cache usage

**Key Deliverable:** Working memory monitoring tab

---

### Sprint 2: Settings System (v0.6.0) - 2-3 days
**Focus:** Persistent settings and theme management
- Create `settings.go` and `settingsui.go`
- Implement JSON-based settings persistence
- Add Settings dialog (theme, intervals, preferences)
- Dark/light theme toggle

**Key Deliverable:** Settings persist between sessions

---

### Sprint 3: Disk & Network (v0.7.0) - 3-4 days
**Focus:** Expand monitoring capabilities
- Create `diskinfo.go`, `disktile.go` (disk I/O monitoring)
- Create `netinfo.go`, `nettile.go` (network monitoring)
- Add Disk and Network tabs
- Per-device/interface breakdown

**Key Deliverable:** Disk and Network monitoring tabs

---

### Sprint 4: Performance (v0.7.5) - 2 days
**Focus:** Optimize rendering and memory
- Create `ringbuffer.go` (circular buffer for history)
- Optimize graph rendering (canvas reuse)
- Reduce memory allocations
- Benchmark improvements

**Key Deliverable:** 30%+ performance improvement

---

### Sprint 5: Cross-Platform (v0.8.0) - 4-5 days
**Focus:** Windows and macOS support
- Create `platform/` directory with abstraction layer
- Implement `platform/linux.go`, `windows.go`, `darwin.go`
- Platform-specific UI adjustments
- Cross-compilation setup

**Key Deliverable:** Builds and runs on Linux, Windows, macOS

---

### Sprint 6: Process Management (v0.9.0) - 3-4 days
**Focus:** Process list and management
- Create `processinfo.go`, `processtable.go`
- Sortable process table (PID, Name, CPU%, Memory)
- Process details dialog
- End process functionality

**Key Deliverable:** Working process manager tab

---

### Sprint 7: Export & Logging (v0.9.5) - 2-3 days
**Focus:** Data export and historical logging
- Create `export.go`, `logger.go`, `exportui.go`
- Export to CSV/JSON/XML
- Continuous logging with rotation
- Historical data viewer

**Key Deliverable:** Export and logging functionality

---

### Sprint 8: GPU Monitoring (v1.0.0) - 3-4 days
**Focus:** GPU support (NVIDIA, AMD, Intel)
- Create `gpuinfo.go`, `gputile.go`
- NVIDIA NVML support
- AMD ROCm/sysfs support
- Intel integrated GPU support

**Key Deliverable:** GPU monitoring tab (v1.0 milestone!)

---

### Sprint 9: Advanced UI (v1.1.0) - 3-4 days
**Focus:** UI enhancements and customization
- Resizable tiles
- Custom layouts (`layouts.go`)
- Mini mode (`minimode.go`)
- Enhanced visualizations
- Always-on-top, system tray

**Key Deliverable:** Highly customizable UI

---

### Sprint 10: Polish & Release (v1.2.0) - 2-3 days
**Focus:** Production readiness
- Bug fixes and final polish
- Complete documentation
- Installers for all platforms (.deb, .rpm, .msi, .dmg)
- Performance final pass
- i18n preparation

**Key Deliverable:** Production-ready v1.2.0 release

---

## Total Timeline: 26-35 days (full-time development)

---

## File Creation Summary

### Sprint 1 (Memory)
- `memorytile.go` (60-80 lines)
- `memorytile_test.go`

### Sprint 2 (Settings)
- `settings.go` (100-120 lines)
- `settingsui.go` (150-200 lines)
- `settings_test.go`, `settingsui_test.go`

### Sprint 3 (Disk/Network)
- `diskinfo.go` (150-180 lines)
- `disktile.go` (80-100 lines)
- `netinfo.go` (150-180 lines)
- `nettile.go` (80-100 lines)
- Test files

### Sprint 4 (Performance)
- `ringbuffer.go` (100-120 lines)
- `ringbuffer_test.go`
- `benchmark_test.go`

### Sprint 5 (Cross-Platform)
- `platform/interface.go`
- `platform/linux.go`
- `platform/windows.go`
- `platform/darwin.go`
- Platform-specific tests

### Sprint 6 (Processes)
- `processinfo.go` (200-250 lines)
- `processtable.go` (250-300 lines)
- `processdetails.go` (150-200 lines)
- Test files

### Sprint 7 (Export)
- `export.go` (150-180 lines)
- `logger.go` (120-150 lines)
- `exportui.go` (100-120 lines)
- `historyviewer.go` (200-250 lines)
- Test files

### Sprint 8 (GPU)
- `gpuinfo.go` (100-120 lines)
- `gpuinfo_nvidia.go` (150-200 lines)
- `gpuinfo_amd.go` (150-200 lines)
- `gpuinfo_intel.go` (100-150 lines)
- `gputile.go` (100-120 lines)
- Test files

### Sprint 9 (Advanced UI)
- `layouts.go` (150-200 lines)
- `minimode.go` (100-120 lines)
- Test files

### Sprint 10 (Polish)
- Documentation files
- Installer scripts
- i18n files

---

## Architecture Evolution

### Before (v0.4.1) - Flat Structure
```
cycles/
├── main.go
├── config.go
├── theme.go
├── tile.go
├── sysinfo.go
├── graphics.go
└── info.go
```

### After (v1.2.0) - Modular Structure
```
cycles/
├── cmd/cycles/main.go
├── internal/
│   ├── app/
│   ├── config/
│   ├── ui/
│   │   ├── components/
│   │   ├── dialogs/
│   │   └── views/
│   ├── platform/
│   ├── monitoring/
│   ├── data/
│   └── utils/
├── docs/
├── installers/
└── i18n/
```

---

## Key Features by Version

| Version | Key Features |
|---------|-------------|
| 0.4.1 | CPU monitoring, command-line config |
| 0.5.0 | + Memory monitoring |
| 0.6.0 | + Persistent settings, themes |
| 0.7.0 | + Disk & Network monitoring |
| 0.7.5 | + Performance optimizations |
| 0.8.0 | + Windows & macOS support |
| 0.9.0 | + Process management |
| 0.9.5 | + Export & logging |
| 1.0.0 | + GPU monitoring ⭐ |
| 1.1.0 | + Advanced UI features |
| 1.2.0 | + Production polish 🎉 |

---

## Sprint Dependencies

```
Sprint 1 (Memory)
    ↓
Sprint 2 (Settings) ←─── Can run parallel with Sprint 1
    ↓
Sprint 3 (Disk/Net) ←─── Depends on Sprint 2 for settings
    ↓
Sprint 4 (Performance) ←─ Depends on Sprints 1-3 data structures
    ↓
Sprint 5 (Cross-Platform) ←─ Refactors all previous work
    ↓
Sprint 6 (Processes) ←─── Independent, can be earlier
    ↓
Sprint 7 (Export) ←────── Depends on all monitoring features
    ↓
Sprint 8 (GPU) ←────────  Independent, can be earlier
    ↓
Sprint 9 (Advanced UI) ←─ Depends on all UI components
    ↓
Sprint 10 (Polish) ←───── Depends on everything
```

---

## Testing Requirements per Sprint

| Sprint | Unit Tests | Integration Tests | Platform Tests | Performance Tests |
|--------|-----------|-------------------|----------------|-------------------|
| 1 | ✓ | - | - | - |
| 2 | ✓ | ✓ | - | - |
| 3 | ✓ | ✓ | - | ✓ |
| 4 | ✓ | - | - | ✓✓ |
| 5 | ✓ | ✓ | ✓✓ | - |
| 6 | ✓ | ✓ | - | ✓ |
| 7 | ✓ | ✓ | - | - |
| 8 | ✓ | - | ✓ | - |
| 9 | ✓ | ✓ | - | ✓ |
| 10 | ✓✓ | ✓✓ | ✓✓ | ✓✓ |

---

## Priority Ranking (if resources limited)

### Must Have (P0)
1. Sprint 1: Memory Monitoring
2. Sprint 2: Settings System
3. Sprint 5: Cross-Platform Support
4. Sprint 10: Polish & Release

### Should Have (P1)
5. Sprint 3: Disk & Network
6. Sprint 4: Performance
7. Sprint 6: Process Management

### Nice to Have (P2)
8. Sprint 7: Export & Logging
9. Sprint 8: GPU Monitoring
10. Sprint 9: Advanced UI

---

## Getting Started

### To Start Sprint 1:
```bash
# Create new files
touch memorytile.go memorytile_test.go

# Run existing tests to ensure baseline
go test -v ./...

# Begin implementation
# See OVERHAUL_PLAN_V2.md Sprint 1 section for details
```

### To Track Progress:
- Use GitHub Issues for each sprint
- Create milestones for each version
- Tag commits with sprint numbers
- Update CHANGELOG.md after each sprint

---

*For detailed task breakdowns, see OVERHAUL_PLAN_V2.md*
