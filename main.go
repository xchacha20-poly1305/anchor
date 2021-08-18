package main

import (
	"flag"
	"github.com/sagernet/sagerconnect/api"
	"github.com/sagernet/sagerconnect/core"
	"github.com/sagernet/sagerconnect/tun"
	"github.com/xjasonlyu/tun2socks/log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const Version = "0.1.0"

func main() {
	flags := flag.NewFlagSet("SagerConnect", flag.ContinueOnError)

	showHelp := flags.Bool("help", false, "show help and exit")
	showVersion := flags.Bool("version", false, "show version and exit")

	//headless := flags.Bool("h", false, "don't show gui")

	err := flags.Parse(os.Args)
	if err != nil {
		println(err)
		flags.Usage()
		os.Exit(1)
	}

	if *showHelp {
		flags.Usage()
		return
	}

	if *showVersion {
		println(Version)
		return
	}

	simple()

	/*if *headless {
		return
	}*/

	/*applicationInit()

	mainApp.Run()*/
}

func simple() {
	log.SetLevel(log.InfoLevel)
	core.Must(core.ExecSu())

	conn, err := net.ListenUDP("udp", nil)
	core.Must(err)

	deviceName, err := os.Hostname()
	core.Must(err)

	query, err := api.MakeQuery(api.Query{Version: api.Version, DeviceName: deviceName})
	core.Must(err)

	//core.Must0(api.ParseQuery(query))

	_, err = conn.WriteTo(query, &net.UDPAddr{
		IP:   net.IPv4bcast,
		Port: 11451,
	})
	core.Must(err)

	buffer := make([]byte, 2048)
	core.Must(conn.SetReadDeadline(time.Now().Add(5 * time.Second)))
	length, addr, err := conn.ReadFromUDP(buffer)
	if err != nil && strings.Contains(err.Error(), "timeout") {
		log.Fatalf("no device found")
	}
	core.Must(err)
	core.Must(conn.Close())

	response, err := api.ParseResponse(buffer[:length])
	core.Must(err)

	log.Infof("found %s (%s)", response.DeviceName, addr.IP.String())

	tunName := "tun0"
	if len(os.Args) > 1 {
		tunName = os.Args[1]
	}

	log.Infof("socks port: %d", response.SocksPort)
	log.Infof("dns port: %d", response.DnsPort)
	log.Infof("enable log: %v", response.Debug)

	tun2socks, err := tun.NewTun2socks(tunName, addr.IP.String(), int(response.SocksPort), int(response.DnsPort), response.Debug)
	core.Must(err)
	tun2socks.Start()

	cmd, err := tun.AddRoute(tunName, response.BypassLan)
	if err != nil {
		tun2socks.Close()
		log.Fatalf("add route failed: %s: %v\n", cmd, err)
	}

	log.Infof("%s started", tunName)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.SetLevel(log.InfoLevel)
	tun2socks.Close()
	log.Infof("closed")
}
