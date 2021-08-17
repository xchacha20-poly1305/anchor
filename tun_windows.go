package main

import (
	"fmt"
	"net"
)

func addRoute(name string, bypassLan bool) (cmd string, err error) {
	cmd, err = execShell("netsh", "interface", "set", "address", name, "static", PrivateVlan4Client, "30")
	if err != nil {
		return
	}

	cmd, err = execShell("netsh", "interface", "ipv6", "add", "address", name, "static", fmt.Sprintf("%s/126", PrivateVlan6Client))
	if err != nil {
		return
	}

	if bypassLan {
		for _, addr := range BypassPrivateRoute {
			_, inet, _ := net.ParseCIDR(addr)
			cmd, err = execShell("netsh", "routing", "ip", "add", "persistentroute", inet.IP.String(), inet.Mask.String(), name, PrivateVlan4Client)
			if err != nil {
				return
			}
		}
	} else {
		cmd, err = execShell("netsh", "routing", "ip", "add", "persistentroute", "0.0.0.0", "0", name, PrivateVlan4Client)
		if err != nil {
			return
		}
	}

	cmd, err = execShell("netsh", "interface", "ipv6", "add", "route", "::/0", name)
	return
}
