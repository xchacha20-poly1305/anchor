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

	"github.com/xchacha20-poly1305/anchor"
)

type RejectFunc func(source net.Addr, deviceName string) bool

// Anchor implements a service that provides anchor protocol.
// What's more, it is reusable.
type Anchor struct {
	ctx          context.Context
	logger       logger.ContextLogger
	access       sync.RWMutex
	packetConn   net.PacketConn
	response     []byte
	shouldReject RejectFunc
	startOnce    sync.Once
}

func New(ctx context.Context, logger logger.ContextLogger, response *anchor.Response, shouldReject RejectFunc) *Anchor {
	return &Anchor{
		ctx:          ctx,
		logger:       logger,
		response:     common.Must1(response.MarshalBinary()),
		shouldReject: shouldReject,
	}
}

func (a *Anchor) Start(packetConn net.PacketConn) error {
	notFirstStart := true
	a.startOnce.Do(func() {
		notFirstStart = false
		a.packetConn = packetConn
		a.access = common.DefaultValue[sync.RWMutex]()
		go a.loop()
	})
	if notFirstStart {
		return os.ErrExist
	}
	return nil
}

func (a *Anchor) UpdateResponse(response *anchor.Response, shouldReject RejectFunc) {
	a.access.Lock()
	defer a.access.Unlock()
	if response != nil {
		a.response = common.Must1(response.MarshalBinary())
	}
	a.shouldReject = shouldReject
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
			continue
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

	a.access.RLock()
	defer a.access.RUnlock()
	_, err = a.packetConn.WriteTo(a.response, source)
	if err != nil {
		a.logger.InfoContext(a.ctx, "send response: ", err)
	}
}

func (a *Anchor) Close() error {
	a.logger.DebugContext(a.ctx, "closing Anchor server")
	a.startOnce = common.DefaultValue[sync.Once]() // reusable
	return a.packetConn.Close()
}
