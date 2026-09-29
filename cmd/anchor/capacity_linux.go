package main

import (
	"os"

	"github.com/sagernet/sing/common"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/shell"

	"golang.org/x/sys/unix"
)

func checkCapacity() error {
	header := unix.CapUserHeader{Version: unix.LINUX_CAPABILITY_VERSION_3}
	var data [2]unix.CapUserData
	err := unix.Capget(&header, &data[0])
	if err != nil {
		return E.Cause(err, "capget")
	}
	const capability = unix.CAP_NET_ADMIN
	const index = capability / 32
	const mask = uint32(1) << (capability % 32)
	for _, set := range []uint32{data[index].Effective, data[index].Permitted} {
		if set&mask == 0 {
			return os.ErrPermission
		}
	}
	return nil
}

func setCapacity() error {
	programPath := common.Must1(os.Executable())
	return shell.Exec("sudo", "setcap", "cap_net_admin=ep", programPath).Attach().Run()
}
