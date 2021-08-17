package main

import (
	"fmt"
	"os/exec"
)

func addRoute(name string) (cmd string, err error) {
	_, ipRoute2NotFound := exec.LookPath("ip")
	if ipRoute2NotFound == nil {
		cmd, err = execShell("ip", "addr", "add", fmt.Sprintf("%s/30", PrivateVlan4Client), "dev", name)
		if err != nil {
			return
		}
		cmd, err = execShell("ip", "addr", "add", fmt.Sprintf("%s/126", PrivateVlan6Client), "dev", name)
		if err != nil {
			return
		}

		cmd, err = execShell("ip", "link", "set", "dev", name, "up")
		if err != nil {
			return
		}

		cmd, err = execShell("ip", "route", "add", "default", "dev", name)
		if err != nil {
			return
		}

		cmd, err = execShell("ip", "route", "flush", "cache")

		return
	} else {
		cmd, err = execShell("ifconfig", name, PrivateVlan4Client, "netmask", "30")
		if err != nil {
			return
		}

		cmd, err = execShell("ifconfig", name, "add", fmt.Sprintf("%s/126", PrivateVlan6Client))
		if err != nil {
			return
		}

		cmd, err = execShell("route", "add", "-net", "0.0.0.0", "netmask", "0", name)
		if err != nil {
			return
		}

		cmd, err = execShell("route", "add", "-A", "inet6", "::/0", name)
	}
	return
}
