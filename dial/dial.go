package dial

import (
	"runtime"

	tun "github.com/sagernet/sing-tun"
	"github.com/sagernet/sing/common/control"
	M "github.com/sagernet/sing/common/metadata"
	N "github.com/sagernet/sing/common/network"
)

var _ N.Dialer = (*Dialer)(nil)

type Dialer struct {
	N.DefaultDialer
}

func New(finder control.InterfaceFinder, monitor tun.DefaultInterfaceMonitor, interfaceName string) (dialer *Dialer) {
	dialer = &Dialer{}
	if finder == nil {
		finder = control.NewDefaultInterfaceFinder()
	}
	var bindFunc control.Func
	if interfaceName == "" {
		bindFunc = control.BindToInterfaceFunc(finder, func(network string, address string) (name string, index int, err error) {
			remoteAddr := M.ParseSocksaddr(address).Addr
			switch runtime.GOOS {
			case "linux", "android":
				name, index = monitor.DefaultInterface(remoteAddr)
				if index == -1 {
					err = tun.ErrNoRoute
				}
			default:
				index = monitor.DefaultInterfaceIndex(remoteAddr)
				if index == -1 {
					err = tun.ErrNoRoute
				}
			}
			return
		})
	} else {
		bindFunc = control.BindToInterface(finder, interfaceName, -1)
	}
	dialer.Dialer.Control = control.Append(dialer.Dialer.Control, bindFunc)
	dialer.ListenConfig.Control = control.Append(dialer.ListenConfig.Control, bindFunc)
	return dialer
}
