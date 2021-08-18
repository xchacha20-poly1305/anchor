package tun

import (
	"fmt"
	"github.com/sagernet/sagerconnect/core"
	"net"
)

func AddRoute(name string, bypassLan bool) (cmd string, err error) {
	cmd, err = core.ExecShell(cmd, err, "netsh", "interface", "set", "address", name, "static", PrivateVlan4Client, "30")
	cmd, err = core.ExecShell(cmd, err, "netsh", "interface", "ipv6", "add", "address", name, "static", fmt.Sprintf("%s/126", PrivateVlan6Client))
	if bypassLan {
		for _, addr := range BypassPrivateRoute {
			_, inet, _ := net.ParseCIDR(addr)
			cmd, err = core.ExecShell(cmd, err, "netsh", "routing", "ip", "add", "persistentroute", inet.IP.String(), inet.Mask.String(), name, PrivateVlan4Client)
		}
	} else {
		cmd, err = core.ExecShell(cmd, err, "netsh", "routing", "ip", "add", "persistentroute", "0.0.0.0", "0", name, PrivateVlan4Client)
	}
	cmd, err = core.ExecShell(cmd, err, "netsh", "interface", "ipv6", "add", "route", "::/0", name)
	return
}
