# Changelog

All notable changes to the Cycles project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
