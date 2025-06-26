package dialers

import (
	"github.com/sagernet/sing-tun"
	"github.com/sagernet/sing/common/control"
	N "github.com/sagernet/sing/common/network"
)

var _ N.Dialer = (*bound)(nil)

type bound struct {
	N.DefaultDialer
}

// NewBound returns a dialer bound to interface, whose name is bindInterface.
// If bindInterface is empty, this will bind to default interface.
func NewBound(finder control.InterfaceFinder, monitor tun.DefaultInterfaceMonitor, bindInterface string) N.Dialer {
	if finder == nil {
		finder = control.NewDefaultInterfaceFinder()
	}
	var bindFunc control.Func
	if bindInterface == "" {
		bindFunc = control.BindToInterfaceFunc(finder, func(network string, address string) (name string, index int, err error) {
			defaultInterface := monitor.DefaultInterface()
			if defaultInterface == nil {
				err = tun.ErrNoRoute
				return
			}
			name = defaultInterface.Name
			index = defaultInterface.Index
			return
		})
	} else {
		bindFunc = control.BindToInterface(finder, bindInterface, -1)
	}
	dialer := &bound{}
	dialer.Dialer.Control = control.Append(dialer.Dialer.Control, bindFunc)
	dialer.ListenConfig.Control = control.Append(dialer.ListenConfig.Control, bindFunc)
	return dialer
}

func (r *bound) Upstream() any {
	return r.DefaultDialer
}
