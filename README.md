# cycles
Desktop CPU Monitor I threw together while trying to debug some issues with my computer. 

![image](https://github.com/TylerCode/cycles/assets/18288340/36332a79-6882-4204-ba6d-51d061798229)


## Overview
This application provides a real-time graphical representation of CPU utilization for each logical core. It displays the utilization percentage and frequency of each core with history going back 20 measurements. 

The ultimate goal is to have something more akin to the performance tab in Windows Task Manager. 

![image](https://github.com/TylerCode/cycles/assets/18288340/e768eee2-d7c8-407a-b680-229ee16a093a)
![image](https://github.com/TylerCode/cycles/assets/18288340/460582ca-6260-4148-a19a-587ae03dc87f)



## Features
- Displays CPU core utilization and frequency.
- Real-time (every 2 seconds) updates for each CPU core.
- Graphical representation of CPU utilization history for the last 20 measurements.
- Customizable UI with system color scheme compatibility.

## Installation

Option 1: 
- Download the latest release on the github release page. [Found here...](https://github.com/TylerCode/cycles/releases)
- Download the "cycles" file or the entire zip
- Make it executable if it isn't already
```
chmod +x cycles
```
- Run it! (double click or `./cycles`)

Option 2:
- Wait for snap package

Option 3:
- Build it yourself!


## Contributing

### Prerequisites
- Go (Golang) installed on your system.
- Fyne library for UI development in Go.
- `gopsutil` library for accessing system information.


### Setup

To set up the project on your local machine:

1. Clone the repository (I would make a fork and clone that but to play around):
```
git clone https://github.com/TylerCode/cycles
```
2. Add dependencies:
```
go get fyne.io/fyne/v2
go get github.com/shirou/gopsutil/cpu
```

### Contrib Notes

Contributions to this project are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature/YourFeature`).
3. Make your changes.
4. Commit your changes (`git commit -am 'Add some feature'`).
5. Push to the branch (`git push origin feature/YourFeature`).
6. Open a Pull Request.

## License
MIT

## Acknowledgments
- Fyne team for the wonderful UI toolkit.
- The `gopsutil` library for system information access.
