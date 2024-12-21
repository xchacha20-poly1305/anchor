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

type RejectFunc func(source net.Addr, deviceName string) (shouldReject bool)

// Anchor implements a service that provides anchor protocol.
type Anchor struct {
	ctx          context.Context
	logger       logger.ContextLogger
	packetConn   net.PacketConn
	response     []byte
	shouldReject RejectFunc
	listen       *net.UDPAddr
	startOnce    sync.Once
}

func New(ctx context.Context, logger logger.ContextLogger, listen *net.UDPAddr, response *anchor.Response, shouldReject RejectFunc) *Anchor {
	listen.Port = anchor.Port
	return &Anchor{
		ctx:          ctx,
		logger:       logger,
		response:     common.Must1(response.MarshalBinary()),
		shouldReject: shouldReject,
		listen:       listen,
	}
}

func (a *Anchor) Start() (err error) {
	notFirstStart := true
	a.startOnce.Do(func() {
		notFirstStart = false
		a.packetConn, err = net.ListenUDP(N.NetworkUDP, a.listen)
		if err != nil {
			return
		}
		go a.loop()
	})
	if err != nil {
		return err
	}
	if notFirstStart {
		return os.ErrExist
	}
	return nil
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
	a.logger.DebugContext(a.ctx, "closing Anchor server")
	return a.packetConn.Close()
}
