//go:build !linux

package main

import (
	"os"

	E "github.com/sagernet/sing/common/exceptions"
)

func checkCapacity() error {
	return os.ErrInvalid
}

func setCapacity() error {
	return E.Cause(os.ErrInvalid, "not linux")
}
