package main

import (
	"context"
	"fmt"
	"github.com/miekg/dns"
	"github.com/xjasonlyu/tun2socks/component/dialer"
	"github.com/xjasonlyu/tun2socks/constant"
	"github.com/xjasonlyu/tun2socks/core"
	"github.com/xjasonlyu/tun2socks/core/device"
	"github.com/xjasonlyu/tun2socks/core/stack"
	"github.com/xjasonlyu/tun2socks/log"
	"github.com/xjasonlyu/tun2socks/proxy"
	"github.com/xjasonlyu/tun2socks/tunnel"
	"net"
	"sync"
	"time"
)

type Tun2socks struct {
	access  sync.Mutex
	stack   *stack.Stack
	proxy   *proxy.Socks5
	device  *device.Device
	dns     string
	dnsAddr *net.UDPAddr
}

type proxyTunnel struct{}

func (*proxyTunnel) Add(conn core.TCPConn) {
	tunnel.Add(conn)
}
func (*proxyTunnel) AddPacket(packet core.UDPPacket) {
	tunnel.AddPacket(packet)
}

func NewTun2socks(name string, addr net.Addr, socksPort int, dnsPort int, debug bool) (*Tun2socks, error) {

	tun, err := openDevice(name)
	if err != nil {
		return nil, err
	}

	gvisor, err := stack.New(tun, &proxyTunnel{}, stack.WithDefault())

	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	socks5Proxy, err := proxy.NewSocks5(fmt.Sprintf("127.0.0.1:%d", socksPort), "", "")
	if err != nil {
		return nil, err
	}

	dnsAddrStr := fmt.Sprintf("127.0.0.1:%d", dnsPort)
	dnsAddr, err := net.ResolveUDPAddr("udp", dnsAddrStr)
	if err != nil {
		return nil, err
	}
	tunnel.SetUDPTimeout(5 * 60)

	return &Tun2socks{
		stack:   gvisor,
		device:  &tun,
		proxy:   socks5Proxy,
		dns:     dnsAddrStr,
		dnsAddr: dnsAddr,
	}, nil
}

func (t *Tun2socks) Start() {
	t.access.Lock()
	defer t.access.Unlock()

	proxy.SetDialer(t)
}

func (t *Tun2socks) Close() {
	t.access.Lock()
	defer t.access.Unlock()

	_ = (*t.device).Close()
	t.stack.Close()
}

func (t *Tun2socks) DialContext(ctx context.Context, metadata *constant.Metadata) (net.Conn, error) {
	if metadata.DstPort == 53 {
		return dialer.DialContext(ctx, "tcp", t.dns)
	}
	return t.proxy.DialContext(ctx, metadata)
}

func (t *Tun2socks) DialUDP(metadata *constant.Metadata) (net.PacketConn, error) {
	if metadata.DstPort == 53 {
		return t.newDnsPacketConn(metadata)
	} else {
		return t.proxy.DialUDP(metadata)
	}
}

func (t *Tun2socks) newDnsPacketConn(metadata *constant.Metadata) (conn net.PacketConn, err error) {
	conn, err = dialer.ListenPacket("udp", "")
	if err == nil {
		conn = &dnsPacketConn{conn: conn, dnsAddr: t.dnsAddr, realAddr: metadata.UDPAddr()}
	}
	return
}

type dnsPacketConn struct {
	conn     net.PacketConn
	notDns   bool
	dnsAddr  *net.UDPAddr
	realAddr net.Addr
}

func (pc *dnsPacketConn) WriteTo(b []byte, addr net.Addr) (int, error) {
	pc.realAddr = addr
	if !pc.notDns {
		req := new(dns.Msg)
		err := req.Unpack(b)
		if err == nil && !req.Response {
			if len(req.Question) > 0 {
				log.Debugf("new dns query: %s", req.Question[0].Name)
			}
			return pc.conn.WriteTo(b, pc.dnsAddr)
		} else {
			pc.notDns = true
		}
	}
	return pc.conn.WriteTo(b, addr)
}

func (pc *dnsPacketConn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	n, realAddr, err := pc.conn.ReadFrom(p)
	if pc.realAddr != nil {
		return n, pc.realAddr, err
	} else {
		return n, realAddr, err
	}
}

func (pc *dnsPacketConn) Close() error {
	return pc.conn.Close()
}

func (pc *dnsPacketConn) LocalAddr() net.Addr {
	return pc.conn.LocalAddr()
}

func (pc *dnsPacketConn) SetDeadline(t time.Time) error {
	return pc.conn.SetDeadline(t)
}

func (pc *dnsPacketConn) SetReadDeadline(t time.Time) error {
	return pc.conn.SetReadDeadline(t)
}

func (pc *dnsPacketConn) SetWriteDeadline(t time.Time) error {
	return pc.conn.SetWriteDeadline(t)
}
