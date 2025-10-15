# Changelog

All notable changes to the Cycles project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
