package main

import (
	"net/netip"

	"github.com/sagernet/sing-tun"
	"github.com/sagernet/sing/common/auth"
	"github.com/sagernet/sing/common/json/badoption"
	"github.com/sagernet/sing/common/logger"
	"github.com/xchacha20-poly1305/anchor/tun2dialer"
)

type Options struct {
	BindInterface            string                           `json:"bind_interface,omitempty"`
	InterfaceName            string                           `json:"interface_name,omitempty"`
	Stack                    string                           `json:"stack,omitempty"`
	MTU                      uint32                           `json:"mtu,omitempty"`
	GSO                      bool                             `json:"gso,omitempty"`
	AutoRoute                *bool                            `json:"auto_route,omitempty"`
	StrictRoute              bool                             `json:"strict_route,omitempty"`
	UDPTimeout               badoption.Duration               `json:"udp_timeout,omitempty"` // Second
	IncludeAllNetworks       bool                             `json:"include_all_networks,omitempty"`
	Inet4Address             badoption.Listable[netip.Prefix] `json:"inet4_address,omitempty"`
	Inet6Address             badoption.Listable[netip.Prefix] `json:"inet6_address,omitempty"`
	Inet4RouteAddress        badoption.Listable[netip.Prefix] `json:"inet4_route_address,omitempty"`
	Inet6RouteAddress        badoption.Listable[netip.Prefix] `json:"inet6_route_address,omitempty"`
	Inet4RouteExcludeAddress badoption.Listable[netip.Prefix] `json:"inet4_route_exclude_address,omitempty"`
	Inet6RouteExcludeAddress badoption.Listable[netip.Prefix] `json:"inet6_route_exclude_address,omitempty"`
	IncludeInterface         badoption.Listable[string]       `json:"include_interface,omitempty"`
	ExcludeInterface         badoption.Listable[string]       `json:"exclude_interface,omitempty"`
	BypassLAN                bool                             `json:"bypass_lan,omitempty"`
	DNS                      badoption.Listable[string]       `json:"dns,omitempty"`
	AutoRedirect             bool                             `json:"auto_redirect"`

	MixedInbound MixedInboundOptions `json:"mixed_inbound,omitempty"`
}

func (o *Options) ForTun2Dialer(ctxLogger logger.ContextLogger, interfaceMonitor tun.DefaultInterfaceMonitor) (tun2dialer.Options, error) {
	dnsServers := make([]netip.Addr, 0, len(o.DNS))
	for _, dns := range o.DNS {
		addr, err := netip.ParseAddr(dns)
		if err != nil {
			return tun2dialer.Options{}, err
		}
		dnsServers = append(dnsServers, addr)
	}
	var (
		autoRedirectInputMark  uint32
		autoRedirectOutputMark uint32
	)
	if o.AutoRedirect {
		autoRedirectInputMark = tun.DefaultAutoRedirectInputMark
		autoRedirectOutputMark = tun.DefaultAutoRedirectOutputMark
	}
	return tun2dialer.Options{
		Options: tun.Options{
			Name:                     o.InterfaceName,
			Inet4Address:             o.Inet4Address,
			Inet6Address:             o.Inet6Address,
			MTU:                      o.MTU,
			GSO:                      o.GSO,
			AutoRoute:                *o.AutoRoute,
			DNSServers:               dnsServers,
			StrictRoute:              o.StrictRoute,
			Inet4RouteAddress:        o.Inet4RouteAddress,
			Inet6RouteAddress:        o.Inet6RouteAddress,
			Inet4RouteExcludeAddress: o.Inet4RouteExcludeAddress,
			Inet6RouteExcludeAddress: o.Inet6RouteExcludeAddress,
			IncludeInterface:         o.IncludeInterface,
			ExcludeInterface:         o.ExcludeInterface,
			IncludeUID:               nil,
			ExcludeUID:               nil,
			InterfaceMonitor:         interfaceMonitor,
			Logger:                   ctxLogger,
			AutoRedirectInputMark:    autoRedirectInputMark,
			AutoRedirectOutputMark:   autoRedirectOutputMark,
		},
		Stack:              o.Stack,
		UDPTimeout:         o.UDPTimeout.Build(),
		IncludeAllNetworks: o.IncludeAllNetworks,
		BypassLAN:          o.BypassLAN,
	}, nil
}

func (o *Options) ApplyDefault() {
	if o.MTU <= 0 {
		o.MTU = tun2dialer.MTU
	}
	if o.Stack == "" {
		o.Stack = tun2dialer.Stack
	}
	if o.InterfaceName == "" {
		o.InterfaceName = tun.CalculateInterfaceName("anchor")
	}
	if len(o.Inet4Address) == 0 {
		o.Inet4Address = []netip.Prefix{netip.MustParsePrefix("172.19.0.1/30")}
	}
	/*if len(o.Inet6Address) == 0 {
		o.Inet6Address = []netip.Prefix{netip.MustParsePrefix("fdfe:dcba:9876::1/126")}
	}*/
	if o.UDPTimeout <= 0 {
		o.UDPTimeout = badoption.Duration(tun2dialer.UDPTimeout)
	}
	if o.AutoRoute == nil {
		b := true
		o.AutoRoute = &b
	}
	if len(o.DNS) == 0 {
		o.DNS = []string{tun2dialer.DNSServer}
	}
}

type MixedInboundOptions struct {
	Listen string                        `json:"listen,omitempty"`
	Users  badoption.Listable[auth.User] `json:"users,omitempty"`
}
