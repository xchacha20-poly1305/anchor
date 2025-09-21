// Package route implements simple connection router.
package route

import (
	"context"
	"strings"

	M "github.com/sagernet/sing/common/metadata"
	N "github.com/sagernet/sing/common/network"
)

type routeInboundContext struct{}

// InboundContext stores metadata of inbound connection.
type InboundContext struct {
	Network             string
	Source, Destination M.Socksaddr
	Override            M.Socksaddr
}

// AppendInboundContext appends inboundContext to ctx.
func AppendInboundContext(ctx context.Context, inboundContext *InboundContext) context.Context {
	return context.WithValue(ctx, routeInboundContext{}, inboundContext)
}

// InboundContextFrom tries getting the InboundContext in ctx.
func InboundContextFrom(ctx context.Context) *InboundContext {
	value := ctx.Value(routeInboundContext{})
	if value == nil {
		return nil
	}
	return value.(*InboundContext) // Must
}

// Clone returns a copy of InboundContext.
func (i *InboundContext) Clone() *InboundContext {
	return &InboundContext{
		Network:     i.Network,
		Source:      i.Source,
		Destination: i.Destination,
		Override:    i.Override,
	}
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

// BypassICMP bypass the ICMP network.
func BypassICMP(dialer N.Dialer) Func {
	if dialer == nil {
		panic(nilDialer)
	}
	return func(ctx *InboundContext) N.Dialer {
		if !strings.Contains(ctx.Network, N.NetworkICMP) {
			return nil
		}
		return dialer
	}
}
