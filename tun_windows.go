package main

import "fmt"

func addRoute(name string) (cmd string, err error) {
	cmd, err = execShell("netsh", "interface", "set", "address", name, "static", PrivateVlan4Client, "30")
	if err != nil {
		return
	}

	cmd, err = execShell("netsh", "interface", "ipv6", "add", "address", name, "static", fmt.Sprintf("%s/126", PrivateVlan6Client))
	if err != nil {
		return
	}

	cmd, err = execShell("netsh", "routing", "ip", "add", "persistentroute", "0.0.0.0", "0", name, PrivateVlan4Client)
	if err != nil {
		return
	}

	cmd, err = execShell("netsh", "interface", "ipv6", "add", "route", "::/0", name)
	return
}
