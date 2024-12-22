// Package tun2dialer provides Tun2Dialer to handle Tun.
package tun2dialer

import (
	"context"
	"io"
	"net"
	"time"

	"github.com/sagernet/sing-tun"
	"github.com/sagernet/sing/common"
	"github.com/sagernet/sing/common/atomic"
	"github.com/sagernet/sing/common/bufio"
	"github.com/sagernet/sing/common/canceler"
	"github.com/sagernet/sing/common/control"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/logger"
	M "github.com/sagernet/sing/common/metadata"
	N "github.com/sagernet/sing/common/network"
)

var _ tun.Handler = (*Tun2Dialer)(nil)

// Tun2Dialer forwards tun connections to dialer.
type Tun2Dialer struct {
	logger               logger.ContextLogger
	dialer, directDialer N.Dialer
	bypassLan            bool
	tunInterface         tun.Tun
	stack                tun.Stack
	udpTimeout           time.Duration
}

// NewTun2Dialer returns Tun2Dialer which forward to dialer.
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
		logger:       ctxLogger,
		dialer:       dialer,
		directDialer: bypassedDialer,
		bypassLan:    tunOptions.BypassLAN,
		udpTimeout:   tunOptions.UDPTimeout,
	}
	t.tunInterface, err = tun.New(tunOptions.Options)
	if err != nil {
		return nil, E.Cause(err, "create tun")
	}
	t.stack, err = tun.NewStack(tunOptions.Stack, tun.StackOptions{
		Context:                ctx,
		Tun:                    t.tunInterface,
		TunOptions:             tunOptions.Options,
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

func (t *Tun2Dialer) Start() error {
	err := t.stack.Start()
	if err != nil {
		return E.Cause(err, "start tun stack")
	}
	err = t.tunInterface.Start()
	if err != nil {
		return E.Cause(err, "start tun interface")
	}
	return nil
}

func (t *Tun2Dialer) PrepareConnection(network string, source, destination M.Socksaddr) error {
	return nil
}

func (t *Tun2Dialer) NewConnectionEx(ctx context.Context, conn net.Conn, source, destination M.Socksaddr, onClose N.CloseHandlerFunc) {
	t.logger.InfoContext(ctx, "inbound connection from ", source)
	t.logger.InfoContext(ctx, "inbound connection to ", destination)

	t.routeConnection(ctx, conn, source, destination, onClose)
}

func (t *Tun2Dialer) routeConnection(ctx context.Context, conn net.Conn, source, destination M.Socksaddr, onClose N.CloseHandlerFunc) {
	var (
		remoteConn net.Conn
		err        error
	)
	if t.bypassLan && !N.IsPublicAddr(source.Addr) {
		remoteConn, err = t.directDialer.DialContext(ctx, N.NetworkTCP, destination)
	} else {
		remoteConn, err = t.dialer.DialContext(ctx, N.NetworkTCP, destination)
	}
	if err != nil {
		err = E.Cause(err, "open outbound connection")
		_ = N.CloseOnHandshakeFailure(conn, onClose, err)
		t.logger.ErrorContext(ctx, err)
		return
	}
	err = N.ReportConnHandshakeSuccess(conn, remoteConn)
	if err != nil {
		err = E.Cause(err, "report handshake success")
		remoteConn.Close()
		_ = N.CloseOnHandshakeFailure(conn, onClose, err)
		t.logger.ErrorContext(ctx, err)
		return
	}

	var done atomic.Bool
	go t.connectionCopy(ctx, conn, remoteConn, false, &done, onClose)
	go t.connectionCopy(ctx, remoteConn, conn, true, &done, onClose)
}

func (t *Tun2Dialer) connectionCopy(ctx context.Context, source io.Reader, destination io.Writer, direction bool, done *atomic.Bool, onClose N.CloseHandlerFunc) {
	originSource := source
	originDestination := destination
	var readCounters, writeCounters []N.CountFunc
	for {
		source, readCounters = N.UnwrapCountReader(source, readCounters)
		destination, writeCounters = N.UnwrapCountWriter(destination, writeCounters)
		if cachedSrc, isCached := source.(N.CachedReader); isCached {
			cachedBuffer := cachedSrc.ReadCached()
			if cachedBuffer != nil {
				dataLen := cachedBuffer.Len()
				_, err := destination.Write(cachedBuffer.Bytes())
				cachedBuffer.Release()
				if err != nil {
					if done.Swap(true) {
						onClose(err)
					}
					common.Close(originSource, originDestination)
					if !direction {
						t.logger.ErrorContext(ctx, "connection upload payload: ", err)
					} else {
						t.logger.ErrorContext(ctx, "connection download payload: ", err)
					}
					return
				}
				for _, counter := range readCounters {
					counter(int64(dataLen))
				}
				for _, counter := range writeCounters {
					counter(int64(dataLen))
				}
			}
			continue
		}
		break
	}
	_, err := bufio.CopyWithCounters(destination, source, originSource, readCounters, writeCounters)
	if err != nil {
		common.Close(originDestination)
	} else if duplexDst, isDuplex := destination.(N.WriteCloser); isDuplex {
		err = duplexDst.CloseWrite()
		if err != nil {
			common.Close(originSource, originDestination)
		}
	} else {
		common.Close(originDestination)
	}
	if done.Swap(true) {
		onClose(err)
		common.Close(originSource, originDestination)
	}
	if !direction {
		if err == nil {
			t.logger.DebugContext(ctx, "connection upload finished")
		} else if !E.IsClosedOrCanceled(err) {
			t.logger.ErrorContext(ctx, "connection upload closed: ", err)
		} else {
			t.logger.TraceContext(ctx, "connection upload closed")
		}
	} else {
		if err == nil {
			t.logger.DebugContext(ctx, "connection download finished")
		} else if !E.IsClosedOrCanceled(err) {
			t.logger.ErrorContext(ctx, "connection download closed: ", err)
		} else {
			t.logger.TraceContext(ctx, "connection download closed")
		}
	}
}

func (t *Tun2Dialer) NewPacketConnectionEx(ctx context.Context, conn N.PacketConn, source, destination M.Socksaddr, onClose N.CloseHandlerFunc) {
	t.logger.InfoContext(ctx, "inbound redirect connection from ", source)
	t.logger.InfoContext(ctx, "inbound connection to ", destination)

	t.routePacketConn(ctx, conn, source, destination, onClose)
}

func (t *Tun2Dialer) routePacketConn(ctx context.Context, conn N.PacketConn, source, destination M.Socksaddr, onClose N.CloseHandlerFunc) {
	var (
		remotePacketConn net.PacketConn
		err              error
	)
	if t.bypassLan && destination.Port != DNSPort && !N.IsPublicAddr(source.Addr) {
		remotePacketConn, err = t.directDialer.ListenPacket(ctx, destination)
	} else {
		remotePacketConn, err = t.dialer.ListenPacket(ctx, destination)
	}
	if err != nil {
		_ = N.CloseOnHandshakeFailure(conn, onClose, err)
		t.logger.ErrorContext(ctx, "listen outbound packet connection: ", err)
		return
	}
	err = N.ReportPacketConnHandshakeSuccess(conn, remotePacketConn)
	if err != nil {
		conn.Close()
		remotePacketConn.Close()
		t.logger.ErrorContext(ctx, "report handshake success: ", err)
		return
	}
	if t.udpTimeout > 0 {
		ctx, conn = canceler.NewPacketConn(ctx, conn, t.udpTimeout)
	}
	remoteConn := bufio.NewPacketConn(remotePacketConn)

	var done atomic.Bool
	go t.packetConnectionCopy(ctx, conn, remoteConn, false, &done, onClose)
	go t.packetConnectionCopy(ctx, remoteConn, conn, true, &done, onClose)
}

func (t *Tun2Dialer) packetConnectionCopy(ctx context.Context, source N.PacketReader, destination N.PacketWriter, direction bool, done *atomic.Bool, onClose N.CloseHandlerFunc) {
	_, err := bufio.CopyPacket(destination, source)
	if !direction {
		if E.IsClosedOrCanceled(err) {
			t.logger.TraceContext(ctx, "packet upload closed")
		} else {
			t.logger.DebugContext(ctx, "packet upload closed: ", err)
		}
	} else {
		if E.IsClosedOrCanceled(err) {
			t.logger.TraceContext(ctx, "packet download closed")
		} else {
			t.logger.DebugContext(ctx, "packet download closed: ", err)
		}
	}
	if !done.Swap(true) {
		onClose(err)
	}
	_ = common.Close(source, destination)
}
