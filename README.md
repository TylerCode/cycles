[![Build Snap](https://snapcraft.io/tylercode-cycles/badge.svg)](https://snapcraft.io/tylercode-cycles)
[![Build AppImage](https://github.com/TylerCode/cycles/actions/workflows/appimage.yml/badge.svg)](https://github.com/TylerCode/cycles/actions/workflows/appimage.yml)

# cycles
Desktop CPU Monitor I threw together while trying to debug some issues with my computer. Have not been dedicating a lot of time to this lately since I'm not on a machine with Snap access anymore. 

![image](https://github.com/TylerCode/cycles/assets/18288340/36332a79-6882-4204-ba6d-51d061798229)


## Overview
This application provides a real-time graphical representation of CPU utilization for each logical core. It displays the utilization percentage and frequency of each core with history going back 30 measurements. 

The ultimate goal is to have something more akin to the performance tab in Windows Task Manager. 

![image](https://github.com/TylerCode/cycles/assets/18288340/460582ca-6260-4148-a19a-587ae03dc87f)



## Features
- Displays CPU core utilization and frequency in real-time
- Memory tab showing total/used/free/cached memory with a usage graph
- Settings dialog (File > Preferences) for grid columns, history size, update interval, and logical/physical core display, persisted between runs
- Light/dark/auto theme, with a quick toggle in the View menu
- Customizable update interval (default: every 2 seconds)
- Utilization graphs showing historical data
- Command-line flags for customization (override saved settings):
  - `--columns`: Number of columns in the grid layout (default: 8)
  - `--interval`: Update interval (default: 2s)
  - `--history`: Number of historical data points to keep (default: 30)
  - `--logical`: Show logical cores vs physical cores (default: true)


## Installation
Currently, this application is only available on x86 machines running linux. ARM64 is available via snap with the `--edge` flag but it's untested. 


### Snap Store
- Install the snap package with
```bash
sudo snap install tylercode-cycles
```
- A "Release" build will be out once I've tested it on a few other machines.


### AppImage
- Download the latest .AppImage from the [releases page.](https://github.com/TylerCode/cycles/releases)
- Run the AppImage


### Plain old executable
- Download the latest release on the github release page. [Found here...](https://github.com/TylerCode/cycles/releases)
- Download the "cycles" file or the entire zip
- Make it executable if it isn't already
```
chmod +x cycles
```
- Run it! (double click or `./cycles`)


## Contributing

### Prerequisites
- Go (Golang) installed on your system.
- Fyne library for UI development in Go.
- `gopsutil` library for accessing system information.


### Quick Start (Automated Setup)

The easiest way to get started developing:

1. Clone the repository (I would make a fork and clone that to contribute):
```bash
git clone https://github.com/TylerCode/cycles
cd cycles
```

2. Run automated setup (detects your OS and installs all dependencies):
```bash
make setup
```

3. Build and run:
```bash
make run
```

That's it! The `make setup` command will:
- Detect your Linux distribution (Ubuntu/Debian/Zorin/Pop, Fedora/RHEL, Arch)
- Install all required system dependencies
- Download and verify Go dependencies
- Test that everything builds correctly

**Note:** If your specific distribution isn't recognized, the script will automatically detect your package manager (apt-get, dnf, or pacman) and use the appropriate installation method.

### Manual Setup

If you prefer manual installation or the automated setup doesn't work:

1. Clone the repository:
```bash
git clone https://github.com/TylerCode/cycles
cd cycles
```

2. Install system dependencies:

**Ubuntu/Debian:**
```bash
sudo apt-get install libgl1-mesa-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev libglfw3-dev libxxf86vm-dev pkg-config gcc make
```

**Fedora/RHEL:**
```bash
sudo dnf install mesa-libGL-devel libXcursor-devel libXrandr-devel libXinerama-devel libXi-devel glfw-devel libXxf86vm-devel pkg-config gcc gcc-c++ make
```

**Arch Linux:**
```bash
sudo pacman -S mesa libxcursor libxrandr libxinerama libxi glfw-x11 pkg-config gcc make
```

3. Install Go dependencies:
```bash
make install-deps
```

4. Build the application:
```bash
make build
```

5. Run it:
```bash
make run
```

6. Run tests:
```bash
make test
```

### Build System Commands

```bash
make build      # Build optimized binary
make run        # Build and run
make test       # Run all tests
make check      # Format, vet, and test
make clean      # Remove build artifacts
make help       # Show all available commands
```

### Command-Line Options
Cycles supports several command-line flags for customization:
```bash
./cycles --columns 8 --interval 1s --history 60 --logical=false
```


### Contrib Notes

Contributions to this project will be welcome probably after 0.6!

Please follow these steps:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature/YourFeature`).
3. Make your changes.
4. Commit your changes (`git commit -am 'Add some feature'`).
5. Push to the branch (`git push origin feature/YourFeature`).
6. Open a Pull Request.


## License
MIT


## Acknowledgments
- Fyne team for an incredible cross-platform ui kit.
- The `gopsutil` library for system information access.
