package main

import (
	"runtime"
	"strings"
	"testing"
)

func TestGetSystemInfo(t *testing.T) {
	info := GetSystemInfo()

	if !strings.Contains(info, runtime.GOOS) {
		t.Errorf("Expected system info to contain OS: %s", runtime.GOOS)
	}

	if !strings.Contains(info, runtime.GOARCH) {
		t.Errorf("Expected system info to contain architecture: %s", runtime.GOARCH)
	}

	if !strings.Contains(info, "Go:") {
		t.Error("Expected system info to contain Go version")
	}
}

func TestGetAppInfo(t *testing.T) {
	info := GetAppInfo()

	if !strings.Contains(info, "Cycles") {
		t.Error("Expected app info to contain 'Cycles'")
	}

	if !strings.Contains(info, "0.4.1") {
		t.Error("Expected app info to contain version 0.4.1")
	}

	if !strings.Contains(info, "MIT") {
		t.Error("Expected app info to contain license information")
	}
}
