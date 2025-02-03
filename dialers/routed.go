package dialers

import (
	"context"
	"net"

	M "github.com/sagernet/sing/common/metadata"
	N "github.com/sagernet/sing/common/network"

	"github.com/xchacha20-poly1305/anchor/route"
)

var _ N.Dialer = (*Routed)(nil)

// Routed is a dialer with custom route rules.
type Routed struct {
	dialer N.Dialer
	routes []route.Func
}

// NewRouted returns a new Routed.
// dialer is optional.
func NewRouted(dialer N.Dialer, routes ...route.Func) *Routed {
	if dialer == nil {
		dialer = &N.DefaultDialer{}
	}
	return &Routed{dialer: dialer, routes: routes}
}

func (r *Routed) DialContext(ctx context.Context, network string, destination M.Socksaddr) (net.Conn, error) {
	return r.matchDialer(ctx).DialContext(ctx, network, destination)
}

func (r *Routed) ListenPacket(ctx context.Context, destination M.Socksaddr) (net.PacketConn, error) {
	return r.matchDialer(ctx).ListenPacket(ctx, destination)
}

// AppendRule appends new rule.
func (r *Routed) AppendRule(routeFunc route.Func) {
	if routeFunc == nil {
		return
	}
	r.routes = append(r.routes, routeFunc)
}

func (r *Routed) Upstream() any {
	return r.dialer
}

func (r *Routed) matchDialer(ctx context.Context) N.Dialer {
	inboundContext := route.InboundContextFrom(ctx)
	if inboundContext == nil {
		return r.dialer
	}
	for _, routeFunc := range r.routes {
		dialer := routeFunc(inboundContext)
		if dialer != nil {
			return dialer
		}
	}
	return r.dialer
}
