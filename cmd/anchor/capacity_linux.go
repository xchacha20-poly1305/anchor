package main

import (
	"os"

	"github.com/sagernet/sing/common"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/shell"
	"kernel.org/pub/linux/libs/security/libcap/cap"
)

func checkCapacity() error {
	set := cap.GetProc()
	for _, flag := range []cap.Flag{cap.Effective, cap.Permitted} {
		hasNetAdmin, err := set.GetFlag(flag, cap.NET_ADMIN)
		if err != nil {
			return E.Cause(err, "GetFlag")
		}
		if !hasNetAdmin {
			return os.ErrPermission
		}
	}
	return nil
}

func setCapacity() error {
	programPath := common.Must1(os.Executable())
	return shell.Exec("sudo", "setcap", "cap_net_admin=ep", programPath).Attach().Run()
}
