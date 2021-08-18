package tun

import (
	"context"
	"errors"
	"fmt"
	"github.com/miekg/dns"
	"github.com/xjasonlyu/tun2socks/component/dialer"
	"github.com/xjasonlyu/tun2socks/constant"
	"github.com/xjasonlyu/tun2socks/core"
	"github.com/xjasonlyu/tun2socks/core/device"
	"github.com/xjasonlyu/tun2socks/core/device/tun"
	"github.com/xjasonlyu/tun2socks/core/stack"
	"github.com/xjasonlyu/tun2socks/log"
	"github.com/xjasonlyu/tun2socks/proxy"
	"github.com/xjasonlyu/tun2socks/tunnel"
	"net"
	"sync"
	"time"
)

const (
	PrivateVlan4Client = "172.19.0.1"
	PrivateVlan6Client = "fdfe:dcba:9876::1"
)

var BypassPrivateRoute []string

func init() {
	BypassPrivateRoute = []string{
		"1.0.0.0/8",
		"2.0.0.0/7",
		"4.0.0.0/6",
		"8.0.0.0/7",
		"11.0.0.0/8",
		"12.0.0.0/6",
		"16.0.0.0/4",
		"32.0.0.0/3",
		"64.0.0.0/3",
		"96.0.0.0/6",
		"100.0.0.0/10",
		"100.128.0.0/9",
		"101.0.0.0/8",
		"102.0.0.0/7",
		"104.0.0.0/5",
		"112.0.0.0/10",
		"112.64.0.0/11",
		"112.96.0.0/12",
		"112.112.0.0/13",
		"112.120.0.0/14",
		"112.124.0.0/19",
		"112.124.32.0/21",
		"112.124.40.0/22",
		"112.124.44.0/23",
		"112.124.46.0/24",
		"112.124.48.0/20",
		"112.124.64.0/18",
		"112.124.128.0/17",
		"112.125.0.0/16",
		"112.126.0.0/15",
		"112.128.0.0/9",
		"113.0.0.0/8",
		"114.0.0.0/10",
		"114.64.0.0/11",
		"114.96.0.0/12",
		"114.112.0.0/15",
		"114.114.0.0/18",
		"114.114.64.0/19",
		"114.114.96.0/20",
		"114.114.112.0/23",
		"114.114.115.0/24",
		"114.114.116.0/22",
		"114.114.120.0/21",
		"114.114.128.0/17",
		"114.115.0.0/16",
		"114.116.0.0/14",
		"114.120.0.0/13",
		"114.128.0.0/9",
		"115.0.0.0/8",
		"116.0.0.0/6",
		"120.0.0.0/6",
		"124.0.0.0/7",
		"126.0.0.0/8",
		"128.0.0.0/3",
		"160.0.0.0/5",
		"168.0.0.0/8",
		"169.0.0.0/9",
		"169.128.0.0/10",
		"169.192.0.0/11",
		"169.224.0.0/12",
		"169.240.0.0/13",
		"169.248.0.0/14",
		"169.252.0.0/15",
		"169.255.0.0/16",
		"170.0.0.0/7",
		"172.0.0.0/12",
		"172.32.0.0/11",
		"172.64.0.0/10",
		"172.128.0.0/9",
		"173.0.0.0/8",
		"174.0.0.0/7",
		"176.0.0.0/4",
		"192.0.0.8/29",
		"192.0.0.16/28",
		"192.0.0.32/27",
		"192.0.0.64/26",
		"192.0.0.128/25",
		"192.0.1.0/24",
		"192.0.3.0/24",
		"192.0.4.0/22",
		"192.0.8.0/21",
		"192.0.16.0/20",
		"192.0.32.0/19",
		"192.0.64.0/18",
		"192.0.128.0/17",
		"192.1.0.0/16",
		"192.2.0.0/15",
		"192.4.0.0/14",
		"192.8.0.0/13",
		"192.16.0.0/12",
		"192.32.0.0/11",
		"192.64.0.0/12",
		"192.80.0.0/13",
		"192.88.0.0/18",
		"192.88.64.0/19",
		"192.88.96.0/23",
		"192.88.98.0/24",
		"192.88.100.0/22",
		"192.88.104.0/21",
		"192.88.112.0/20",
		"192.88.128.0/17",
		"192.89.0.0/16",
		"192.90.0.0/15",
		"192.92.0.0/14",
		"192.96.0.0/11",
		"192.128.0.0/11",
		"192.160.0.0/13",
		"192.169.0.0/16",
		"192.170.0.0/15",
		"192.172.0.0/14",
		"192.176.0.0/12",
		"192.192.0.0/10",
		"193.0.0.0/8",
		"194.0.0.0/7",
		"196.0.0.0/7",
		"198.0.0.0/12",
		"198.16.0.0/15",
		"198.20.0.0/14",
		"198.24.0.0/13",
		"198.32.0.0/12",
		"198.48.0.0/15",
		"198.50.0.0/16",
		"198.51.0.0/18",
		"198.51.64.0/19",
		"198.51.96.0/22",
		"198.51.101.0/24",
		"198.51.102.0/23",
		"198.51.104.0/21",
		"198.51.112.0/20",
		"198.51.128.0/17",
		"198.52.0.0/14",
		"198.56.0.0/13",
		"198.64.0.0/10",
		"198.128.0.0/9",
		"199.0.0.0/8",
		"200.0.0.0/7",
		"202.0.0.0/8",
		"203.0.0.0/18",
		"203.0.64.0/19",
		"203.0.96.0/20",
		"203.0.112.0/24",
		"203.0.114.0/23",
		"203.0.116.0/22",
		"203.0.120.0/21",
		"203.0.128.0/17",
		"203.1.0.0/16",
		"203.2.0.0/15",
		"203.4.0.0/14",
		"203.8.0.0/13",
		"203.16.0.0/12",
		"203.32.0.0/11",
		"203.64.0.0/10",
		"203.128.0.0/9",
		"204.0.0.0/6",
		"208.0.0.0/4",
	}
}

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

func NewTun2socks(name string, addr string, socksPort int, dnsPort int, debug bool) (*Tun2socks, error) {

	device, err := tun.Open(tun.WithName(name), tun.WithMTU(1500))
	if err != nil {
		return nil, err
	}

	gvisor, err := stack.New(device, &proxyTunnel{}, stack.WithDefault())

	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	socks5Proxy, err := proxy.NewSocks5(fmt.Sprintf("%s:%d", addr, socksPort), "", "")
	if err != nil {
		return nil, err
	}

	dnsAddrStr := fmt.Sprintf("%s:%d", addr, dnsPort)
	dnsAddr, err := net.ResolveUDPAddr("udp", dnsAddrStr)
	if err != nil {
		return nil, err
	}
	tunnel.SetUDPTimeout(5 * 60)

	return &Tun2socks{
		stack:   gvisor,
		device:  &device,
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
		return dialer.DialContext(ctx, metadata.Network(), t.dns)
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
	if err != nil {
		return
	}
	conn = &dnsPacketConn{conn: conn, dnsAddr: t.dnsAddr, realAddr: metadata.UDPAddr(), proxyConstructor: func() (net.PacketConn, error) {
		return t.proxy.DialUDP(metadata)
	}}
	return
}

type dnsPacketConn struct {
	conn             net.PacketConn
	proxyConn        net.PacketConn
	proxyConstructor func() (net.PacketConn, error)
	notDns           bool
	dnsAddr          *net.UDPAddr
	realAddr         net.Addr
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
			conn, err := pc.proxyConstructor()
			if err != nil {
				return 0, err
			}
			pc.proxyConn = conn
			log.Debugf("not dns query conn to %s", addr.String())
		}
	}
	return pc.proxyConn.WriteTo(b, addr)
}

func (pc *dnsPacketConn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	if pc.notDns {
		if pc.proxyConn == nil {
			return 0, nil, errors.New("unexpected")
		}
		n, addr, err = pc.proxyConn.ReadFrom(p)
	} else {
		n, addr, err = pc.conn.ReadFrom(p)
	}
	if pc.realAddr != nil {
		addr = pc.realAddr
	}

	return
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
