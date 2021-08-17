package main

func addRoute(name string, bypassLan bool) (cmd string, err error) {
	cmd, err = execShell("ifconfig", name, PrivateVlan4Client, "netmask", "30")
	if err != nil {
		return
	}

	cmd, err = execShell("ifconfig", name, "inet6", PrivateVlan6Client, "prefixlen", "126")
	if err != nil {
		return
	}

	if bypassLan {
		for _, addr := range BypassPrivateRoute {
			cmd, err = execShell("route", "add", addr, PrivateVlan4Client)
			if err != nil {
				return
			}
		}
	} else {
		cmd, err = execShell("route", "add", "0.0.0.0/0", PrivateVlan4Client)
		if err != nil {
			return
		}
	}

	cmd, err = execShell("route", "add", "-inet6", "::", PrivateVlan6Client, "-prefixlen", "0")
	return
}
