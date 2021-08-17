package main

import (
	"golang.org/x/sys/windows"
)

func ExecSu() error {
	if windows.Getuid() != 0 {
		windows.SetErrorMode(uint32(windows.ERROR_ELEVATION_REQUIRED))
		windows.Exit(1)
	}
	return nil
}
