// Package sockshttp mixed socks and http inbound.
package sockshttp

import (
	"bufio"
	"context"
	"net"
	"os"

	"github.com/sagernet/sing/common"
	"github.com/sagernet/sing/common/auth"
	"github.com/sagernet/sing/common/logger"
	M "github.com/sagernet/sing/common/metadata"
	N "github.com/sagernet/sing/common/network"
	"github.com/sagernet/sing/protocol/http"
	"github.com/sagernet/sing/protocol/socks"
	"github.com/sagernet/sing/protocol/socks/socks4"
	"github.com/sagernet/sing/protocol/socks/socks5"
)

type HandlerEx interface {
	N.TCPConnectionHandlerEx
	N.UDPConnectionHandlerEx
}

type SocksHttp struct {
	listener net.Listener
	auth     *auth.Authenticator
	ctx      context.Context
	logger   logger.ContextLogger
	detour   HandlerEx
}

func New(ctx context.Context, logger logger.ContextLogger, detour HandlerEx, listen string, authenticator *auth.Authenticator) (*SocksHttp, error) {
	listener, err := net.Listen(N.NetworkTCP, listen)
	if err != nil {
		return nil, err
	}
	return &SocksHttp{
		listener: listener,
		auth:     authenticator,
		ctx:      ctx,
		logger:   logger,
		detour:   detour,
	}, nil
}

func (s *SocksHttp) Start() error {
	if common.Done(s.ctx) {
		return os.ErrClosed
	}
	go s.loop()
	return nil
}

func (s *SocksHttp) loop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.logger.DebugContext(s.ctx, "Sockshttp: exit because: ", err)
			return
		}
		go s.handle(conn)
	}
}

func (s *SocksHttp) handle(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	headByte, err := reader.Peek(1)
	if err != nil {
		s.logger.ErrorContext(s.ctx, "Sockshttp: read first byte: ", err)
		return
	}
	done := make(chan struct{})
	switch headByte[0] {
	case socks4.Version, socks5.Version:
		_ = socks.HandleConnectionEx(s.ctx, conn, reader, s.auth, s.detour, M.SocksaddrFromNet(conn.RemoteAddr()), func(_ error) {
			close(done)
		})
	default:
		_ = http.HandleConnectionEx(s.ctx, conn, reader, s.auth, s.detour, M.SocksaddrFromNet(conn.RemoteAddr()), func(_ error) {
			close(done)
		})
	}
	select {
	case <-s.ctx.Done():
	case <-done:
	}
}

func (s *SocksHttp) Close() error {
	return common.Close(s.listener)
}
