package core

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func Must(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func Must0(_ interface{}, err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
func Must1(_ interface{}, _ interface{}, err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func ExecShell(cmdIn string, errIn error, name string, arg ...string) (cmd string, err error) {
	if err != nil {
		return cmdIn, err
	}
	cmd = strings.Join([]string{name, strings.Join(arg, " ")}, " ")
	shell := exec.Command(name, arg...)
	shell.Stdin = os.Stdin
	shell.Stdout = os.Stdout
	shell.Stderr = os.Stderr
	err = shell.Start()
	if err == nil {
		err = shell.Wait()
	}
	return
}

func ExecProc(name string, arg []string) error {
	shell := exec.Command(name, arg...)
	shell.Stdin = os.Stdin
	shell.Stdout = os.Stdout
	shell.Stderr = os.Stderr
	err := shell.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	status, err := shell.Process.Wait()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(status.ExitCode())
	return nil
}
