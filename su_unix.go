// +build !windows

package main

import (
	"golang.org/x/sys/unix"
	"os"
	"os/exec"
)

func ExecSu() error {
	if unix.Getuid() == 0 {
		return nil
	}

	_, sudoNotExists := exec.LookPath("sudo")
	if sudoNotExists == nil {
		execProc("sudo", os.Args)
	} else {
		args := []string{"-c"}
		for _, arg := range os.Args {
			args = append(args, arg)
		}
		execProc("su", args)
	}

	return nil
}
