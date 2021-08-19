package tun

import "github.com/sagernet/sagerconnect/core"

func AddRoute(name string, bypassLan bool) (cmd string, err error) {
	cmd, err = core.ExecShell(cmd, err, "ifconfig", name, PrivateVlan4Client, "netmask", "30")
	cmd, err = core.ExecShell(cmd, err, "ifconfig", name, "inet6", PrivateVlan6Client, "prefixlen", "126")
	if bypassLan {
		for _, addr := range BypassPrivateRoute {
			cmd, err = core.ExecShell(cmd, err, "route", "add", addr, PrivateVlan4Client)
		}
	} else {
		cmd, err = core.ExecShell(cmd, err, "route", "add", "0.0.0.0/0", PrivateVlan4Client)
	}

	cmd, err = core.ExecShell(cmd, err, "route", "add", "-inet6", "::", PrivateVlan6Client, "-prefixlen", "0")
	return
}
