// Package dialers provides functional dialers.
package dialers

import (
	"context"

	"github.com/sagernet/sing/common"
	"github.com/sagernet/sing/common/control"
	N "github.com/sagernet/sing/common/network"
)

type Controller interface {
	ControlFunc(ctx context.Context) control.Func
}

func GetControlFunc(ctx context.Context, dialer N.Dialer) control.Func {
	if controller, isController := common.Cast[Controller](dialer); isController {
		return controller.ControlFunc(ctx)
	}
	return nil
}
