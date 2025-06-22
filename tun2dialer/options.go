package tun2dialer

import (
	"time"

	"github.com/sagernet/sing-tun"
)

// Options is option for Tun2Dialer
type Options struct {
	tun.Options
	Stack              string
	UDPTimeout         time.Duration
	IncludeAllNetworks bool
	BypassLAN          bool
}

const (
	UDPTimeout = 5 * time.Minute
	MTU        = 9000
	Stack      = "mixed"
	DNSServer  = "8.8.8.8"
)
