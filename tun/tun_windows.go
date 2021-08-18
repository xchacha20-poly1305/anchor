package tun

import (
	"fmt"
	"github.com/sagernet/sagerconnect/core"
)

func AddRoute(name string, bypassLan bool) (cmd string, err error) {
	cmd, err = core.ExecShell(cmd, err, "netsh", "interface", "ip", "set", "address", name, "static", PrivateVlan4Client, "255.255.255.252")
	cmd, err = core.ExecShell(cmd, err, "netsh", "interface", "ipv6", "add", "address", name, fmt.Sprintf("%s/126", PrivateVlan6Client))
	if bypassLan {
		for _, addr := range BypassPrivateRoute {
			cmd, err = core.ExecShell(cmd, err, "netsh", "interface", "ipv4", "add", "route", addr, name, PrivateVlan4Client)
		}
	} else {
		cmd, err = core.ExecShell(cmd, err, "netsh", "interface", "ipv4", "add", "route", "0.0.0.0/0", name, PrivateVlan4Client)
	}
	cmd, err = core.ExecShell(cmd, err, "netsh", "interface", "ipv6", "add", "route", "::/0", name)
	return
}
