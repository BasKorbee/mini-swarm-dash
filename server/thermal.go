package main

import (
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// readNodeTemp reads the average temperature (°C) across all thermal zones
// from /sys/class/thermal/thermal_zone*/temp. Returns nil if unavailable.
func readNodeTemp() *float64 {
	zones, err := filepath.Glob("/sys/class/thermal/thermal_zone*/temp")
	if err != nil || len(zones) == 0 {
		return nil
	}
	var sum float64
	var count int
	for _, path := range zones {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		millideg, err := strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
		if err != nil {
			continue
		}
		sum += float64(millideg) / 1000.0
		count++
	}
	if count == 0 {
		return nil
	}
	avg := math.Round(sum/float64(count)*10) / 10
	return &avg
}
