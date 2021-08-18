package tun

import (
	"fmt"
	"github.com/sagernet/sagerconnect/core"
	"os/exec"
)

func AddRoute(name string, bypassLan bool) (cmd string, err error) {
	_, ipRoute2NotFound := exec.LookPath("ip")
	if ipRoute2NotFound == nil {
		cmd, err = core.ExecShell(cmd, err, "ip", "addr", "add", fmt.Sprintf("%s/30", PrivateVlan4Client), "dev", name)
		cmd, err = core.ExecShell(cmd, err, "ip", "addr", "add", fmt.Sprintf("%s/126", PrivateVlan6Client), "dev", name)
		cmd, err = core.ExecShell(cmd, err, "ip", "link", "set", "dev", name, "up")
		if bypassLan {
			for _, addr := range BypassPrivateRoute {
				cmd, err = core.ExecShell(cmd, err, "ip", "route", "add", addr, "dev", name)
			}
		} else {
			cmd, err = core.ExecShell(cmd, err, "ip", "route", "add", "0.0.0.0/0", "dev", name)
		}
		cmd, err = core.ExecShell(cmd, err, "ip", "route", "add", "::/0", "dev", name)
		cmd, err = core.ExecShell(cmd, err, "ip", "route", "flush", "cache")
	} else {
		cmd, err = core.ExecShell(cmd, err, "ifconfig", name, PrivateVlan4Client, "netmask", "255.255.255.252")
		cmd, err = core.ExecShell(cmd, err, "ifconfig", name, "add", fmt.Sprintf("%s/126", PrivateVlan6Client))
		if bypassLan {
			for _, addr := range BypassPrivateRoute {
				cmd, err = core.ExecShell(cmd, err, "route", "add", "-net", addr, "netmask", "255", name)
			}
		} else {
			cmd, err = core.ExecShell(cmd, err, "route", "add", "-net", "0.0.0.0/0", "netmask", "255", name)
		}
		cmd, err = core.ExecShell(cmd, err, "route", "add", "-A", "inet6", "::/0", name)
	}
	return
}
