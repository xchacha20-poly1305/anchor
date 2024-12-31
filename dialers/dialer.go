// Package dialers provides implementations for N.Dialer.
package dialers

import (
	"github.com/sagernet/sing-tun"
	"github.com/sagernet/sing/common/control"
	N "github.com/sagernet/sing/common/network"
)

var _ N.Dialer = (*Dialer)(nil)

type Dialer struct {
	N.DefaultDialer
}

func New() *Dialer {
	return &Dialer{DefaultDialer: N.DefaultDialer{}}
}

// Bind binds dialer to default interface if bindInterface name is empty
// or binds dialer to the interfaceName you set.
// finder is optional.
func (d *Dialer) Bind(finder control.InterfaceFinder, monitor tun.DefaultInterfaceMonitor, bindInterface string) (dialer *Dialer) {
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
	dialer.Dialer.Control = control.Append(d.Dialer.Control, bindFunc)
	dialer.ListenConfig.Control = control.Append(d.ListenConfig.Control, bindFunc)
	return dialer
}
