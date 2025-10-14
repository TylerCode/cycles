# Developer Guide

This guide helps developers understand the Cycles codebase structure and how to contribute effectively.

## Project Structure

```
cycles/
├── main.go              # Application entry point and orchestration
├── config.go            # Configuration management and command-line flags
├── theme.go             # Theme and color definitions for graphs
├── tile.go              # CoreTile UI component
├── sysinfo.go           # System information retrieval (CPU, memory)
├── graphics.go          # Graph rendering and drawing utilities
├── info.go              # Application and system information display
├── *_test.go            # Unit tests for each module
├── README.md            # User-facing documentation
├── CHANGELOG.md         # Version history and changes
├── OVERHAUL_SUMMARY.md  # Detailed refactoring documentation
└── vendor/              # Vendored dependencies
```

## Module Responsibilities

### main.go
- Application initialization
- Window and menu setup
- Orchestrates other modules
- Starts update goroutine

**Key Functions:**
- `main()` - Entry point, sets up UI and starts monitoring

### config.go
- Configuration structure definition
- Default values
- Command-line flag parsing

**Key Types:**
- `AppConfig` - Main configuration struct

**Key Functions:**
- `DefaultConfig()` - Returns default configuration
- `ParseFlags()` - Parses command-line arguments

### theme.go
- Color definitions for light/dark themes
- Graph color logic based on utilization

**Key Functions:**
- `GetGraphLineColor(status)` - Returns appropriate color for CPU utilization level

### tile.go
- CoreTile UI component
- Individual CPU core display

**Key Types:**
- `CoreTile` - Represents one CPU core's display

**Key Functions:**
- `NewCoreTile()` - Creates a new tile
- `GetContainer()` - Returns the Fyne container

### sysinfo.go
- System information retrieval
- CPU and memory data collection

**Key Types:**
- `MemoryInfo` - Memory statistics

**Key Functions:**
- `GetCPUFrequencies()` - Reads CPU frequencies from /proc/cpuinfo
- `GetMemoryInfo()` - Reads memory info from /proc/meminfo
- `UpdateCPUInfo(tiles)` - Updates all tiles with latest data

### graphics.go
- Graph rendering
- Drawing utilities
- Label formatting

**Key Functions:**
- `DrawGraph(img, data)` - Renders utilization graph
- `drawLine(img, x1, y1, x2, y2, color)` - Bresenham's line algorithm
- `formatCoreLabel(num)` - Formats core number
- `formatUtilLabel(util)` - Formats utilization percentage
- `formatClockLabel(freq)` - Formats clock frequency

### info.go
- Application metadata
- System information formatting

**Key Functions:**
- `GetSystemInfo()` - Returns OS, arch, Go version
- `GetAppInfo()` - Returns app version and info

## Adding New Features

### Adding a Configuration Option

1. Add field to `AppConfig` in `config.go`:
```go
type AppConfig struct {
    // ... existing fields
    NewOption string
}
```

2. Update `DefaultConfig()`:
```go
func DefaultConfig() *AppConfig {
    return &AppConfig{
        // ... existing fields
        NewOption: "default",
    }
}
```

3. Add flag in `ParseFlags()`:
```go
func (c *AppConfig) ParseFlags() {
    // ... existing flags
    flag.StringVar(&c.NewOption, "newoption", c.NewOption, "Description")
    flag.Parse()
}
```

4. Use in `main.go`:
```go
config.NewOption // Access the value
```

### Adding a Test

1. Create or edit `*_test.go` file
2. Write test function:
```go
func TestNewFeature(t *testing.T) {
    result := NewFeature()
    if result != expected {
        t.Errorf("Expected %v, got %v", expected, result)
    }
}
```

3. Run tests:
```bash
go test -v ./...
```

### Adding System Information

1. Add function to `sysinfo.go`:
```go
func GetNewSystemInfo() (InfoType, error) {
    // Read from /proc or use gopsutil
    return info, nil
}
```

2. Update `UpdateCPUInfo()` or create new update function
3. Add to UI in `main.go` or create new tile type

## Development Workflow

### Setup
```bash
# Clone repository
git clone https://github.com/TylerCode/cycles
cd cycles

# Install dependencies
sudo apt-get install libgl1-mesa-dev libxcursor-dev libxrandr-dev \
                     libxinerama-dev libxi-dev libglfw3-dev libxxf86vm-dev

# Get Go dependencies
go mod tidy
go mod vendor
```

### Development
```bash
# Build
go build -o cycles

# Run
./cycles

# Run with options
./cycles --columns 8 --interval 1s

# Format code
go fmt ./...

# Check for issues
go vet ./...

# Run tests
go test -v ./...
```

### Testing Checklist
- [ ] All tests pass (`go test -v ./...`)
- [ ] No vet warnings (`go vet ./...`)
- [ ] Code formatted (`go fmt ./...`)
- [ ] Application builds (`go build`)
- [ ] Application runs without errors
- [ ] New features tested manually
- [ ] Documentation updated

## Code Style Guidelines

1. **Error Handling**: Always check and handle errors
```go
// Good
result, err := someFunction()
if err != nil {
    log.Printf("Error: %v", err)
    return err
}

// Bad
result, _ := someFunction()
```

2. **Naming Conventions**:
   - Exported functions: `PascalCase`
   - Unexported functions: `camelCase`
   - Constants: `PascalCase`
   - Variables: `camelCase`

3. **Comments**:
   - Document all exported functions
   - Use `///` for function documentation (matches existing style)
   - Explain complex algorithms with references

4. **Function Size**:
   - Keep functions under 50 lines when possible
   - Extract complex logic into helper functions
   - One function, one purpose

5. **Testing**:
   - Test exported functions
   - Test edge cases
   - Use table-driven tests for multiple cases

## Common Tasks

### Updating Version
1. Update version in `config.go`:
```go
Version: "0.5.0",
```

2. Update `CHANGELOG.md` with new version section

3. Update `.github/workflows/appimage.yml` with new version number

### Adding Dependencies
```bash
# Add dependency
go get github.com/new/dependency

# Update vendor
go mod tidy
go mod vendor

# Verify build
go build
```

### Debugging
```bash
# Run with race detector
go run -race main.go

# Build with debug symbols
go build -gcflags="all=-N -l" -o cycles

# Use delve debugger
dlv debug
```

## Platform-Specific Notes

### Linux
- Reads CPU info from `/proc/cpuinfo`
- Reads memory info from `/proc/meminfo`
- Requires X11 libraries for GUI

### Future Platforms (TODO)
- **Windows**: Need alternative to `/proc` filesystem
- **macOS**: Different system info APIs
- Consider using `gopsutil` more extensively for cross-platform support

## Architecture Decisions

### Why vendor directory?
- Ensures reproducible builds
- Works in environments without internet
- Faster CI/CD builds

### Why not use gopsutil for CPU frequency?
- Direct `/proc/cpuinfo` reading is faster
- More control over parsing
- Less external dependencies for core functionality

### Why split into so many files?
- Single Responsibility Principle
- Easier testing of individual components
- Better code organization
- Simpler navigation

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes following the style guide
4. Add tests for new functionality
5. Ensure all tests pass
6. Update documentation
7. Commit: `git commit -m 'Add amazing feature'`
8. Push: `git push origin feature/amazing-feature`
9. Open a Pull Request

## Resources

- [Fyne Documentation](https://developer.fyne.io/)
- [gopsutil Library](https://github.com/shirou/gopsutil)
- [Go Testing](https://golang.org/pkg/testing/)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

## Getting Help

- Open an issue on GitHub
- Check existing issues for similar problems
- Review the OVERHAUL_SUMMARY.md for detailed implementation notes
- Consult the CHANGELOG.md for recent changes
