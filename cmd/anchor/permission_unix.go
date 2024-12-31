//go:build unix

package main

import (
	"errors"
	"os"
)

func checkPermission() error {
	const rootPID = 0
	if os.Geteuid() == rootPID {
		return nil
	}
	err := checkCapacity()
	if err == nil {
		return nil
	}
	if errors.Is(err, os.ErrInvalid) {
		return os.ErrPermission
	}
	if !errors.Is(err, os.ErrPermission) {
		return err
	}
	err = setCapacity()
	if err == nil {
		_, _ = os.Stderr.WriteString("Set capacity, please restart anchor!")
		os.Exit(0)
	}
	return err
}
