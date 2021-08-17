package main

import (
	"github.com/sagernet/sagerconnect/api"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	log.SetFlags(0)

	conn, err := net.ListenUDP("udp", nil)
	must(err)

	deviceName, err := os.Hostname()
	must(err)

	query, err := api.MakeQuery(api.Query{DeviceName: deviceName})
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
	length, addr, err := conn.ReadFrom(buffer)
	if err != nil && strings.Contains(err.Error(), "timeout") {
		log.Fatalln("no device found")
	}
	must(err)

	response, err := api.ParseResponse(buffer[:length])
	must(err)

	log.Printf("connect to %s (%s)\n", response.DeviceName, addr.String())

	tunName := "tun0"
	if len(os.Args) > 0 {
		tunName = os.Args[0]
	}

	log.Println("open ", tunName)

	tun2socks, err := NewTun2socks(tunName, addr, int(response.SocksPort), int(response.DnsPort), response.Debug)
	must(err)
	tun2socks.Start()

	defer tun2socks.Close()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}

func must(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
