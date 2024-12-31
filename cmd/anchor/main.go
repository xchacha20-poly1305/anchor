package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/sagernet/sing-tun"
	"github.com/sagernet/sing/common"
	"github.com/sagernet/sing/common/control"
	E "github.com/sagernet/sing/common/exceptions"
	F "github.com/sagernet/sing/common/format"
	"github.com/sagernet/sing/common/json"
	M "github.com/sagernet/sing/common/metadata"
	"github.com/sagernet/sing/protocol/socks"

	"github.com/xchacha20-poly1305/anchor"
	"github.com/xchacha20-poly1305/anchor/dialers"
	"github.com/xchacha20-poly1305/anchor/log"
	"github.com/xchacha20-poly1305/anchor/route"
	"github.com/xchacha20-poly1305/anchor/tun2dialer"
	"go.uber.org/zap/zapcore"
)

//go:generate goversioninfo --platform-specific

const VERSION = "v0.4.1"

func main() {
	fs := flag.NewFlagSet("anchor", flag.ExitOnError)
	showVersion := fs.Bool("v", false, "Show version")
	logLevel := fs.String("l", zapcore.WarnLevel.String(), "Log level")
	logOutput := fs.String("o", Stderr, "Log output.")
	configPath := fs.String("c", "", "Configuration file path")
	selectedIndex := fs.Int("d", -1, "selected device index (skip select)")
	immediately := fs.Bool("i", false, "skip waiting when the first device is found")
	remoteIp := fs.String("a", "", "remote ip address (skip scan)")
	socksPort := fs.Int("socks", 2080, "remote socks port (skip scan)")
	dnsPort := fs.Int("dns", 6450, "remote dns port (skip scan)")
	_ = fs.Parse(os.Args[1:])

	if *showVersion || isShowVersion(fs.Arg(0)) {
		_, _ = os.Stdout.WriteString(F.ToString("Anchor: ", VERSION))
		os.Exit(0)
	}
	err := checkPermission()
	if err != nil {
		log.Fatal(err)
	}

	level, err := zapcore.ParseLevel(*logLevel)
	if err != nil {
		log.Fatal(E.Cause(err, "parse log level"))
	}
	output, err := parseLogOutput(*logOutput)
	if err != nil {
		log.Fatal(E.Cause(err, "parse log output"))
	}
	defer common.Close(output)
	logger := log.New(output, level)
	defer logger.Sync()

	config := &Options{}
	for {
		if *configPath == "" {
			break
		}
		content, err := os.Open(*configPath)
		if err != nil {
			logger.Fatal(E.Cause(err, "read config file"))
		}
		decoder := json.NewDecoder(content)
		decoder.DisallowUnknownFields()
		err = decoder.Decode(config)
		if err != nil {
			logger.Fatal(err)
		}
		break
	}
	config.ApplyDefault()

	var devices []*ScanResult
	if *remoteIp == "" {
		// scan devices

		deviceName, err := os.Hostname()
		if err != nil {
			logger.Fatal("get host name: ", err)
		}

		query, err := anchor.Query{Version: anchor.Version, DeviceName: deviceName}.MarshalBinary()
		if err != nil {
			logger.Fatal("make scan query: ", err)
		}

		ifAddrs := listInterfaceAddr(logger)
		if len(ifAddrs) == 0 {
			logger.Fatal("failed to get available network interfaces")
		}

		for _, addr := range ifAddrs {
			result, err := addr.scan(query)
			if err != nil {
				if !E.IsTimeout(err) {
					logger.Error(err)
				}
				continue
			}
			devices = append(devices, result)
			logger.Info("Found ", len(devices), ". ", result.response.DeviceName, " (", result.addr.IP, ")")
			if *immediately {
				break
			}
		}
	} else {
		ip := net.ParseIP(*remoteIp)
		if ip == nil {
			logger.Fatal("Failed to parse remote address: ", *remoteIp)
		}

		devices = append(devices, &ScanResult{
			addr: &net.UDPAddr{
				IP:   ip,
				Port: -1,
			},
			response: &anchor.Response{
				Version:   anchor.Version,
				SocksPort: uint16(*socksPort),
				DnsPort:   uint16(*dnsPort),
			},
		})
	}

	deviceSize := len(devices)
	var selected *ScanResult
	if deviceSize == 0 {
		logger.Fatal("no devices found")
		return
	} else if deviceSize > 1 {
		if *selectedIndex != -1 {
			if deviceSize < *selectedIndex {
				logger.Fatal("Invalid device selected: ", *selectedIndex)
				return
			}
			selected = devices[*selectedIndex-1]
		} else {
			for {
				_, _ = os.Stdout.WriteString("Select device to connect: ")
				var line string
				_, err = fmt.Scanln(&line)
				if err != nil {
					logger.Fatal("failed to read selection: ", err)
				}
				index, err := strconv.ParseUint(line, 10, 8)
				if err != nil || index <= 0 || int(index) > deviceSize {
					logger.Error("Invalid device selected: ", line)
					continue
				}
				selected = devices[index-1]
				break
			}
		}
	} else {
		selected = devices[0]
	}

	if *remoteIp == "" {
		logger.Info("Selected ", selected.response.DeviceName, " (", selected.addr.IP, ")")
	}

	if *socksPort != 2080 {
		selected.response.SocksPort = uint16(*socksPort)
	}

	if *socksPort != 6450 {
		selected.response.DnsPort = uint16(*dnsPort)
	}

	logger.Info("SOCKS port: ", selected.response.SocksPort)
	logger.Info("DNS port: ", selected.response.DnsPort)

	interfaceFinder := control.NewDefaultInterfaceFinder()
	networkMonitor, err := tun.NewNetworkUpdateMonitor(logger)
	if err != nil {
		logger.Fatal("Create network update monitor: ", err)
	}
	err = networkMonitor.Start()
	if err != nil {
		logger.Fatal("Start network monitor: ", err)
	}
	interfaceMonitor, err := tun.NewDefaultInterfaceMonitor(networkMonitor, logger, tun.DefaultInterfaceMonitorOptions{
		InterfaceFinder:       interfaceFinder,
		OverrideAndroidVPN:    false,
		UnderNetworkExtension: false,
	})
	if err != nil {
		logger.Fatal("Create default interface monitor: ", err)
	}
	err = interfaceMonitor.Start()
	if err != nil {
		logger.Fatal("Start interface monitor: ", err)
	}
	directDialer := dialers.NewBound(interfaceFinder, interfaceMonitor, config.BindInterface)
	serverAddr := M.SocksaddrFromNet(selected.addr)
	serverAddr.Port = selected.response.SocksPort
	socksDialer := socks.NewClient(
		directDialer,
		serverAddr,
		socks.Version5,
		"",
		"",
	)

	routedDialer := dialers.NewRouted(socksDialer)
	routedDialer.AppendRule(route.UdpDnsPort(dialers.NewOverridden(socksDialer, func(destination M.Socksaddr) M.Socksaddr {
		oldFqdn := destination.Fqdn
		destination = M.SocksaddrFromNetIP(selected.addr.AddrPort())
		destination.Fqdn = oldFqdn
		return destination
	})))
	routedDialer.AppendRule(route.Lan(directDialer))

	tunOption, err := config.ForTun2Dialer(logger, interfaceMonitor)
	if err != nil {
		logger.Fatal("Build Tun config: ", err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	instance, err := tun2dialer.NewTun2Dialer(ctx, logger, tunOption, interfaceFinder, routedDialer)
	if err != nil {
		logger.Fatal("Create tun2dialer instance: ", err)
	}
	err = instance.Start()
	if err != nil {
		logger.Fatal("try start tun2dialer: ", err)
	}
	logger.Warn("Started")

	<-ctx.Done()
	cancel()
	err = common.Close(instance, interfaceMonitor, networkMonitor)
	if err != nil {
		log.Fatal("Close got error: ", err)
	}
	logger.Info("Exit")
}

func isShowVersion(str string) bool {
	switch strings.ToLower(str) {
	case "v", "version":
		return true
	default:
		return false
	}
}
