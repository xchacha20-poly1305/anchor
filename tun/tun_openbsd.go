package tun

import "github.com/sagernet/sagerconnect"

func AddRoute(name string, bypassLan bool) (cmd string, err error) {
	cmd, err = main.execShell("ifconfig", name, PrivateVlan4Client, "netmask", "30")
	if err != nil {
		return
	}

	cmd, err = main.execShell("ifconfig", name, "inet6", PrivateVlan6Client, "prefixlen", "126")
	if err != nil {
		return
	}

	if bypassLan {
		for _, addr := range BypassPrivateRoute {
			cmd, err = main.execShell("route", "add", addr, PrivateVlan4Client)
			if err != nil {
				return
			}
		}
	} else {
		cmd, err = main.execShell("route", "add", "0.0.0.0/0", PrivateVlan4Client)
		if err != nil {
			return
		}
	}

	cmd, err = main.execShell("route", "add", "-inet6", "::", PrivateVlan6Client, "-prefixlen", "0")
	return
}
