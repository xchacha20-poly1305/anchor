package main

import (
	"io"
	"os"
)

const (
	Stdout = "stdout"
	Stderr = "stderr"
)

func parseLogOutput(name string) (writer io.Writer, err error) {
	switch name {
	case Stdout:
		return os.Stdout, nil
	case Stderr:
		return os.Stderr, nil
	}
	file, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}
	return file, nil
}
