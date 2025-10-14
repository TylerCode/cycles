package main

import (
	"fmt"
	"runtime"
)

// GetSystemInfo returns a formatted string with system information
func GetSystemInfo() string {
	return fmt.Sprintf("OS: %s | Arch: %s | Go: %s",
		runtime.GOOS,
		runtime.GOARCH,
		runtime.Version())
}

// GetAppInfo returns formatted application information
func GetAppInfo() string {
	return fmt.Sprintf("Cycles v%s\nCPU Monitor for Linux\n\n%s\n\nLicense: MIT\nAuthor: Tyler C",
		DefaultConfig().Version,
		GetSystemInfo())
}
