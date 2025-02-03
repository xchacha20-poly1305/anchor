package route

import (
	"context"

	M "github.com/sagernet/sing/common/metadata"
	N "github.com/sagernet/sing/common/network"
)

type routeInboundContext struct{}

type InboundContext struct {
	Network             string
	Source, Destination M.Socksaddr
	Override            M.Socksaddr
}

func AppendInboundContext(ctx context.Context, inboundContext *InboundContext) context.Context {
	return context.WithValue(ctx, routeInboundContext{}, inboundContext)
}

func InboundContextFrom(ctx context.Context) *InboundContext {
	value := ctx.Value(routeInboundContext{})
	if value == nil {
		return nil
	}
	return value.(*InboundContext) // Must
}

// Func returns a dialer according to ctx.
// If match the rule, it will return a not-nil N.Dialer.
type Func func(ctx *InboundContext) N.Dialer

const nilDialer = "nil dialer"

// Lan matches destination is Lan.
func Lan(dialer N.Dialer) Func {
	if dialer == nil {
		panic(nilDialer)
	}
	return func(ctx *InboundContext) N.Dialer {
		if N.IsPublicAddr(ctx.Destination.Addr) {
			return nil
		}
		return dialer
	}
}

// UdpDnsPort matches destination is UDP DNS port (53).
func UdpDnsPort(dialer N.Dialer) Func {
	if dialer == nil {
		panic(nilDialer)
	}
	return func(ctx *InboundContext) N.Dialer {
		const dnsPort = 53
		if ctx.Destination.Port != dnsPort {
			return nil
		}
		return dialer
	}
}
