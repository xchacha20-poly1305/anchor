// +build linux

package main

import (
	"github.com/xjasonlyu/tun2socks/core/device"
	"github.com/xjasonlyu/tun2socks/core/device/tun"
)

func openDevice(name string) (device.Device, error) {
	return tun.Open(tun.WithName(name), tun.WithMTU(1500))
}
