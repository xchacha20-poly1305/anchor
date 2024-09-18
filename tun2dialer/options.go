package tun2dialer

import (
	"time"

	tun "github.com/sagernet/sing-tun"
)

// Options is option for Tun2Dialer
type Options struct {
	tun.Options
	Stack                  string
	EndPointIndependentNat bool
	UDPTimeout             int64
	IncludeAllNetworks     bool
	BypassLAN              bool
}

const (
	UDPTimeout = 5 * time.Minute
	MTU        = 9000
	Stack      = "mixed"
)
