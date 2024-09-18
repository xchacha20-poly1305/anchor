// Package tun2dialer provides Tun2Dialer to handle Tun.
package tun2dialer

import (
	"context"
	"net"

	"github.com/sagernet/sing-tun"
	"github.com/sagernet/sing/common"
	"github.com/sagernet/sing/common/bufio"
	"github.com/sagernet/sing/common/control"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/logger"
	M "github.com/sagernet/sing/common/metadata"
	N "github.com/sagernet/sing/common/network"
)

var _ tun.Handler = (*Tun2Dialer)(nil)

// Tun2Dialer forwards tun connections to dialer.
type Tun2Dialer struct {
	logger                 logger.ContextLogger
	dialer, bypassedDialer N.Dialer
	bypassLan              bool
	tunInterface           tun.Tun
	stack                  tun.Stack
}

// NewTun2Dialer returns Tun2Dialer which forward to dialer.
//
// ctxLogger is optional.
func NewTun2Dialer(ctx context.Context,
	ctxLogger logger.ContextLogger,
	tunOptions Options,
	interfaceFinder control.InterfaceFinder,
	dialer, bypassedDialer N.Dialer,
) (*Tun2Dialer, error) {
	var err error
	if dialer == nil {
		return nil, E.New("missing dialer")
	}
	if ctxLogger == nil {
		ctxLogger = logger.NOP()
	}
	t := &Tun2Dialer{
		logger:         ctxLogger,
		dialer:         dialer,
		bypassedDialer: bypassedDialer,
		bypassLan:      tunOptions.BypassLAN,
	}
	t.tunInterface, err = tun.New(tunOptions.Options)
	if err != nil {
		return nil, E.Cause(err, "create tun")
	}
	t.stack, err = tun.NewStack(tunOptions.Stack, tun.StackOptions{
		Context:                ctx,
		Tun:                    t.tunInterface,
		TunOptions:             tunOptions.Options,
		EndpointIndependentNat: tunOptions.EndPointIndependentNat,
		UDPTimeout:             tunOptions.UDPTimeout,
		Handler:                t,
		Logger:                 ctxLogger,
		ForwarderBindInterface: false,
		IncludeAllNetworks:     tunOptions.IncludeAllNetworks,
		InterfaceFinder:        interfaceFinder,
	})
	if err != nil {
		return nil, E.Cause(err, "create stack")
	}
	return t, nil
}

func (t *Tun2Dialer) NewConnection(ctx context.Context, conn net.Conn, metadata M.Metadata) (err error) {
	if metadata.Protocol != "" {
		t.logger.InfoContext(ctx, "inbound ", metadata.Protocol, " connection from: ", metadata.Source)
	} else {
		t.logger.InfoContext(ctx, "inbound connection from: ", metadata.Source)
	}
	t.logger.InfoContext(ctx, "inbound connection to: ", metadata.Destination)
	var newConn net.Conn
	if t.bypassLan && !N.IsPublicAddr(metadata.Destination.Addr) {
		newConn, err = t.bypassedDialer.DialContext(ctx, N.NetworkTCP, metadata.Destination)
	} else {
		newConn, err = t.dialer.DialContext(ctx, N.NetworkTCP, metadata.Destination)
	}
	if err != nil {
		t.logger.WarnContext(ctx, "failed to dial: ", err)
		return err
	}
	defer newConn.Close()
	return bufio.CopyConn(ctx, conn, newConn)
}

func (t *Tun2Dialer) NewPacketConnection(ctx context.Context, conn N.PacketConn, metadata M.Metadata) (err error) {
	if metadata.Protocol != "" {
		t.logger.InfoContext(ctx, "inbound ", metadata.Protocol, " packet from: ", metadata.Source)
	} else {
		t.logger.InfoContext(ctx, "inbound packet from: ", metadata.Source)
	}
	t.logger.InfoContext(ctx, "inbound packet to: ", metadata.Destination)
	var packetConn net.PacketConn
	if t.bypassLan && !N.IsPublicAddr(metadata.Destination.Addr) {
		packetConn, err = t.bypassedDialer.ListenPacket(ctx, metadata.Destination)
	} else {
		packetConn, err = t.dialer.ListenPacket(ctx, metadata.Destination)
	}
	if err != nil {
		return E.Cause(err, "listen packet")
	}
	defer packetConn.Close()
	wrappedConn := bufio.NewPacketConn(packetConn)
	return bufio.CopyPacketConn(ctx, wrappedConn, conn)
}

func (t *Tun2Dialer) NewError(ctx context.Context, err error) {
	_ = common.Close(err)
	if E.IsClosedOrCanceled(err) {
		t.logger.DebugContext(ctx, "connection closed: ", err)
		return
	}
	t.logger.ErrorContext(ctx, err)
}

func (t *Tun2Dialer) Start() error {
	return t.stack.Start()
}

func (t *Tun2Dialer) Close() error {
	return common.Close(t.stack, t.tunInterface)
}
