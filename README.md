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
- Customizable update interval (default: every 2 seconds)
- Utilization graphs showing historical data
- Command-line flags for customization:
  - `--columns`: Number of columns in the grid layout (default: 4)
  - `--interval`: Update interval (default: 2s)
  - `--history`: Number of historical data points to keep (default: 30)
  - `--logical`: Show logical cores vs physical cores (default: true)
- Improved error handling and logging
- Modular code structure for easier maintenance


## Installation
Currently, this application is only available on x86 machines running linux. ARM64 is available via snap with the `--edge` flag but it's untested. 


### Snap Store
- Install the snap package with
```bash
sudo snap install tylercode-cycles
```
- A "Release" build will be out once I've tested it on a few other machines.


### FlatPak
- Coming in 0.5


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


### Setup
To set up the project on your local machine:

1. Clone the repository (I would make a fork and clone that to contribute):
```bash
git clone https://github.com/TylerCode/cycles
cd cycles
```

2. Install system dependencies (Ubuntu/Debian):
```bash
sudo apt-get install libgl1-mesa-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev libglfw3-dev libxxf86vm-dev
```

3. Install Go dependencies:
```bash
go mod tidy
```

4. Build the application:
```bash
go build -o cycles
```

5. Run it:
```bash
./cycles
```

6. Run tests:
```bash
go test -v ./...
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
