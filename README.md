# cycles
Desktop CPU Monitor I threw together while trying to debug some issues with my computer. 

## Overview
This application provides a real-time graphical representation of CPU utilization for each core, including those with hyperthreading. It displays the utilization percentage and frequency of each core in a grid layout, with each core represented as a tile. The application also includes a history graph for the last 20 measurements of CPU utilization for each core.

![image](https://github.com/TylerCode/cycles/assets/18288340/e768eee2-d7c8-407a-b680-229ee16a093a)


## Features
- Displays CPU core utilization and frequency.
- Real-time updates for each CPU core.
- Graphical representation of CPU utilization history for the last 20 measurements.
- Customizable UI with system color scheme compatibility.

## Prerequisites
- Go (Golang) installed on your system.
- Fyne library for UI development in Go.
- `gopsutil` library for accessing system information.

## Installation

Snap package coming soon...

To set up the project on your local machine:

1. Clone the repository:
```
git clone https://github.com/TylerCode/cycles
```
2. Add dependencies:
```
go get fyne.io/fyne/v2
go get github.com/shirou/gopsutil/cpu
```

## Usage
To run the application:
```
go run main.go
```


Alternatively, you can build an executable:

```
go build
./cycles
```


## Contributing
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
