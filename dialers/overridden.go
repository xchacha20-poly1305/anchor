package dialers

import (
	"context"
	"net"

	M "github.com/sagernet/sing/common/metadata"
	N "github.com/sagernet/sing/common/network"

	"github.com/xchacha20-poly1305/anchor/route"
)

var _ N.Dialer = (*Overridden)(nil)

// Overridden overrides destination.
type Overridden struct {
	dialer   N.Dialer
	override func(destination M.Socksaddr) M.Socksaddr
}

// NewOverridden returns a new Overridden.
// dialer is optional while override must not nil.
func NewOverridden(dialer N.Dialer, override func(destination M.Socksaddr) M.Socksaddr) *Overridden {
	if override == nil {
		override(M.Socksaddr{}) // Make it panic
	}
	if dialer == nil {
		dialer = &N.DefaultDialer{}
	}
	return &Overridden{dialer: dialer, override: override}
}

func (o *Overridden) DialContext(ctx context.Context, network string, destination M.Socksaddr) (net.Conn, error) {
	destination = o.override(destination)
	tryUpdateInboundContext(ctx, destination)
	return o.dialer.DialContext(ctx, network, destination)
}

func (o *Overridden) ListenPacket(ctx context.Context, destination M.Socksaddr) (net.PacketConn, error) {
	destination = o.override(destination)
	tryUpdateInboundContext(ctx, destination)
	return o.dialer.ListenPacket(ctx, destination)
}

func tryUpdateInboundContext(ctx context.Context, destination M.Socksaddr) {
	inboundContext := route.InboundContextFrom(ctx)
	if inboundContext == nil {
		return
	}
	inboundContext.Destination = destination
}
