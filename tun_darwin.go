package main

func addRoute(name string) (cmd string, err error) {
	cmd, err = execShell("ifconfig", name, PrivateVlan4Client, "netmask", "30")
	if err != nil {
		return
	}

	cmd, err = execShell("ifconfig", name, "inet6", PrivateVlan6Client, "prefixlen", "126")
	if err != nil {
		return
	}

	cmd, err = execShell("route", "add", "0.0.0.0/0", "-interface", name)
	if err != nil {
		return
	}

	cmd, err = execShell("route", "add", "-inet6", "::/0", "-interface", name)
	return
}
