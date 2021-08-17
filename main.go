package main

import (
	"github.com/sagernet/sagerconnect/api"
	"github.com/xjasonlyu/tun2socks/log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var tun2socks *Tun2socks
var eventListener *api.EventListener
var sigCh chan os.Signal

func main() {
	log.SetLevel(log.InfoLevel)

	conn, err := net.ListenUDP("udp", nil)
	must(err)

	eventListener, err := api.NewEventListener(func(event *api.Event) {
		switch event.EventType {
		case api.EventClose:
			log.Infof("close signal received")
			sigCh <- syscall.SIGINT
		}
	})

	deviceName, err := os.Hostname()
	must(err)

	query, err := api.MakeQuery(api.Query{DeviceName: deviceName, ListenerPort: eventListener.LocalPort()})
	must(err)

	_, err = api.ParseQuery(query)
	must(err)

	_, err = conn.WriteTo(query, &net.UDPAddr{
		IP:   net.IPv4bcast,
		Port: 11451,
	})
	must(err)

	buffer := make([]byte, 2048)
	must(conn.SetReadDeadline(time.Now().Add(5 * time.Second)))
	length, addr, err := conn.ReadFromUDP(buffer)
	if err != nil && strings.Contains(err.Error(), "timeout") {
		log.Fatalf("no device found")
	}
	must(err)

	response, err := api.ParseResponse(buffer[:length])
	must(err)

	log.Infof("found %s (%s)", response.DeviceName, addr.IP.String())

	tunName := "tun0"
	if len(os.Args) > 1 {
		tunName = os.Args[1]
	}

	log.Infof("socks port: %d", response.SocksPort)
	log.Infof("dns port: : %d", response.DnsPort)
	log.Infof("enable log: %v", response.Debug)

	tun2socks, err = NewTun2socks(tunName, addr.IP.String(), int(response.SocksPort), int(response.DnsPort), response.Debug)
	must(err)
	tun2socks.Start()

	cmd, err := addRoute(tunName)
	if err != nil {
		tun2socks.Close()
		log.Fatalf("add route failed: %s: %v\n", cmd, err)
	}

	log.Infof("%s started", tunName)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	destroy()
	log.Infof("closed")
}

func execShell(name string, arg ...string) (cmd string, err error) {
	cmd = strings.Join([]string{name, strings.Join(arg, " ")}, " ")
	shell := exec.Command(name, arg...)
	shell.Stdin = os.Stdin
	shell.Stdout = os.Stdout
	shell.Stderr = os.Stderr
	err = shell.Start()
	if err == nil {
		err = shell.Wait()
	}
	return
}

func must(err error) {
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func destroy() {
	if tun2socks != nil {
		tun2socks.Close()
	}
	if eventListener != nil {
		eventListener.Close()
	}
}
