# Changelog

All notable changes to the Cycles project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.8.0] - 2026-07-10

### Added
- **CPU tab redesign**: replaced the fixed-column tile grid with a
  responsive `GridWrap` layout that reflows column count as the window
  resizes, plus a new List view (compact per-core rows with inline
  utilization bars) toggleable from an aggregate stats strip
  (Threads/Avg util/Peak core/Max clock).
- **Memory tab redesign**: replaced the single memory tile with a full
  dashboard — usage gauge, a used/cached/buffers/free/swap breakdown, and a
  persistent area chart. Swap is now read from `/proc/meminfo`
  (`SwapTotal`/`SwapFree`) and rendered for the first time.
- Yellow utilization/warning band (previously defined in `theme.go` but
  never rendered) is now wired up across CPU bars and the memory
  breakdown's swap row.
- Live theme switching now repaints custom `canvas.Text`/`canvas.Rectangle`
  elements on the CPU and Memory tabs, which don't repaint themselves
  automatically the way built-in widgets do.
- `bar.go`: shared `Bar` component (a `fyne.Layout`-driven progress bar with
  a fixed fill color) used by both the CPU list view and the Memory
  breakdown rows.

### Fixed
- Build break from `fyne.Color`/`fyne.VariantLight`/`fyne.VariantDark`
  (not real Fyne v2.4.2 types) in `settings.go`/`settingsui.go`.
- Preferences silently not persisting — `main.go` now calls
  `app.NewWithID("us.tylerc.cycles")` instead of `app.New()`, which Fyne's
  Preferences API requires.
- Infinite recursion / stack overflow on startup in `CustomTheme`, caused by
  delegating color lookups to the app's current theme (itself) instead of
  `theme.DefaultTheme()`.
- Version-string drift between `config.go`, `snap/snapcraft.yaml`, and the
  AppImage workflow — the AppImage workflow now derives its version from
  `config.go` at build time instead of hardcoding it.
- `snap/snapcraft.yaml` migrated from `base: core20` (no longer supported by
  current Launchpad snapcraft) to `base: core22` with the `gpu` extension
  replacing the old hand-rolled graphics content-interface plugs.

### Changed
- README and internal docs reconciled with actual current features (Memory
  tab, Settings dialog, theme toggle); stale planning/retrospective docs
  removed from repo root.

## [0.6.0] - 2025-11-17

### Added
- **Persistent Settings System** (`settings.go`):
  - Settings automatically save and load using Fyne preferences
  - Theme preference (Auto, Light, Dark)
  - Grid columns, history size, update interval
  - Logical/physical cores preference
- **Settings UI Dialog** (`settingsui.go`):
  - Comprehensive settings dialog with sliders and controls
  - Theme selector dropdown
  - Grid columns slider (1-16)
  - History size slider (10-100)
  - Update interval input
  - Logical cores checkbox
  - Reset to defaults button with confirmation
- **Enhanced Theme System** (`theme.go`):
  - Dynamic theme switching without restart
  - Custom theme implementation
  - `ApplyTheme()` function for programmatic theme changes
  - `isDarkTheme()` helper function
- **Improved Menu System** (update `main.go`):
  - File menu with "Preferences..." option
  - View menu with "Toggle Theme" quick action
  - Reorganized menu structure (File | View | Help)
  - Settings accessible via File → Preferences
- **Command-line Flag Priority**:
  - Command-line flags now override saved settings
  - Seamless integration between CLI and persistent settings

### Changed
- Application now loads saved preferences on startup
- Theme applies immediately on selection
- Settings persist across application restarts
- Menu structure reorganized for better UX

### Technical Improvements
- Settings stored using Fyne's preferences API
- Platform-independent settings storage
- Type-safe settings management
- Comprehensive settings validation
- Unit tests for settings functionality

### Fixed
- Theme switching now works reliably
- Settings apply correctly on save

## [0.5.0] - 2025-11-17

### Added
- **Memory Monitoring Tab**: New dedicated tab for system memory monitoring
  - Real-time memory usage display (Total, Used, Free, Cached)
  - Memory usage percentage with historical graph
  - Automatic history tracking
- **Tabbed Interface**: Switched from single-view to tabbed layout
  - CPU tab: Existing CPU core monitoring
  - Memory tab: New memory monitoring view
- **Memory Tile Component** (`memorytile.go`):
  - Displays comprehensive memory statistics
  - Shows usage graph over time
  - Formatted memory sizes (GB/MB)
- **Enhanced System Information** (`sysinfo.go`):
  - `GetMemoryInfoDetailed()`: Reads detailed memory info from /proc/meminfo
  - `UpdateMemoryInfo()`: Updates memory tiles with current data
  - Tracks MemTotal, MemFree, MemAvailable, Cached, and Buffers
- **Memory Utility Functions**:
  - `formatMemorySize()`: Human-readable memory sizes
  - `formatMemoryPercent()`: Formatted percentage display
- **Unit Tests** (`memorytile_test.go`):
  - Tests for memory tile creation
  - Tests for memory formatting functions
  - Validation of component initialization

### Changed
- Main window now uses tabbed layout instead of single grid
- CPU tiles renamed to `cpuTiles` for clarity
- Separate update goroutines for CPU and Memory monitoring
- Window adapts to tabbed content structure

### Technical Improvements
- Modular memory monitoring component following existing tile pattern
- Independent update loops for different metric types
- Reuses existing graph rendering infrastructure
- Follows established code organization patterns

## [0.4.1] - 2025-10-15

### Changed
- Updated version number consistently across all project files

## [0.4.0] - 2025-10-14

### Added
- Command-line flags for customization:
  - `--columns`: Configure number of columns in grid layout (default: 4)
  - `--interval`: Set update interval (default: 2s)
  - `--history`: Configure number of historical data points (default: 30)
  - `--logical`: Toggle between logical and physical cores (default: true)
- Unit tests for core functionality
- Proper error handling and logging throughout the application
- Configuration system for managing application settings

### Changed
- Refactored codebase into multiple files for better organization:
  - `config.go`: Configuration management
  - `theme.go`: Theme and color management
  - `tile.go`: CoreTile structure and methods
  - `sysinfo.go`: System information retrieval
  - `graphics.go`: Graph rendering utilities
  - `main.go`: Application orchestration
- Improved error messages with proper logging
- Updated version number from 0.3.4 to 0.4.0
- Enhanced code documentation and comments
- Updated README with comprehensive setup instructions

### Fixed
- Improved theme detection logic for graph colors
- Better error handling for icon loading
- More robust CPU frequency reading

### Technical Improvements
- Separated concerns for better code maintainability
- Added test coverage for formatting functions and utilities
- Cleaner main function with configuration-driven behavior
- Removed unused code and comments

## [0.3.4] - Previous Version

### Features
- Basic CPU monitoring for each core
- Real-time utilization graphs
- Frequency display for each core
- Fixed 4-column grid layout
- 2-second update interval
