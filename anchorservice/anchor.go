// Package anchorservice implements a service that provides anchor protocol.
package anchorservice

import (
	"context"
	"net"
	"os"
	"sync"

	"github.com/sagernet/sing/common"
	"github.com/sagernet/sing/common/buf"
	"github.com/sagernet/sing/common/logger"
	N "github.com/sagernet/sing/common/network"

	"github.com/xchacha20-poly1305/anchor"
)

// RejectFunc decides whether to reject the client query.
type RejectFunc func(source net.Addr, deviceName string) (shouldReject bool)

// ListenPacket returns a net.PacketConn.
type ListenPacket func(ctx context.Context) (net.PacketConn, error)

// Anchor implements a service that provides anchor protocol.
type Anchor struct {
	ctx          context.Context
	logger       logger.ContextLogger
	packetConn   net.PacketConn
	response     []byte
	shouldReject RejectFunc
	listen       ListenPacket
	startOnce    sync.Once
}

// New returns a new Anchor service instance.
// contextLogger, listen and shouldReject are optional.
func New(ctx context.Context, contextLogger logger.ContextLogger, listen ListenPacket, response *anchor.Response, shouldReject RejectFunc) *Anchor {
	if contextLogger == nil {
		contextLogger = logger.NOP()
	}
	if listen == nil {
		listen = func(ctx context.Context) (net.PacketConn, error) {
			return net.ListenUDP(N.NetworkUDP, &net.UDPAddr{
				IP:   net.IPv4zero,
				Port: anchor.Port,
			})
		}
	}
	return &Anchor{
		ctx:          ctx,
		logger:       contextLogger,
		response:     common.Must1(response.MarshalBinary()),
		shouldReject: shouldReject,
		listen:       listen,
	}
}

func (a *Anchor) Start() (err error) {
	notFirstStart := true
	a.startOnce.Do(func() {
		notFirstStart = false
		var packetConn net.PacketConn
		packetConn, err = a.listen(a.ctx)
		if err != nil {
			return
		}
		a.packetConn = packetConn
		go a.loop()
	})
	if notFirstStart {
		return os.ErrExist
	}
	return
}

func (a *Anchor) loop() {
	for {
		select {
		case <-a.ctx.Done():
			return
		default:
		}

		buffer := buf.NewSize(anchor.MaxQuerySize)
		_, source, err := buffer.ReadPacketFrom(a.packetConn)
		if err != nil {
			buffer.Release()
			a.logger.WarnContext(a.ctx, "stop loop because: ", err)
			return
		}
		a.logger.DebugContext(a.ctx, "new packet from: ", source)
		go a.handle(source, buffer)
	}
}

func (a *Anchor) handle(source net.Addr, buffer *buf.Buffer) {
	defer buffer.Release()

	query, err := anchor.ParseQuery(buffer)
	if err != nil {
		a.logger.InfoContext(a.ctx, "parse query: ", err)
		return
	}
	if a.shouldReject != nil && a.shouldReject(source, query.DeviceName) {
		a.logger.InfoContext(a.ctx, "reject query from: ", query.DeviceName)
		return
	}
	a.logger.DebugContext(a.ctx, "received query from: ", query.DeviceName)

	_, err = a.packetConn.WriteTo(a.response, source)
	if err != nil {
		a.logger.InfoContext(a.ctx, "send response: ", err)
	}
}

func (a *Anchor) Close() error {
	if a == nil || a.packetConn == nil {
		return os.ErrInvalid
	}
	a.logger.DebugContext(a.ctx, "closing Anchor server")
	buf.Put(a.response)
	a.response = nil
	return a.packetConn.Close()
}

func (a *Anchor) Upstream() any {
	return a.packetConn
}
