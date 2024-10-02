package main

import (
	"net"
	"time"

	E "github.com/sagernet/sing/common/exceptions"
	N "github.com/sagernet/sing/common/network"

	"github.com/xchacha20-poly1305/anchor"
	"github.com/xchacha20-poly1305/anchor/log"
)

type ScanResult struct {
	addr         *net.UDPAddr
	netInterface *net.Interface
	response     *anchor.Response
}

func (a *interfaceWithAddr) scan(query []byte) (*ScanResult, error) {
	result := &ScanResult{
		netInterface: a.netInterface,
	}
	conn, err := net.ListenUDP(N.NetworkUDP+"4", &net.UDPAddr{
		IP: a.addr,
	})
	if err != nil {
		return nil, E.Cause(err, "create multicast conn on if ", a.netInterface.Name)
	}
	_, err = conn.WriteTo(query, &net.UDPAddr{
		IP:   net.IPv4bcast,
		Port: anchor.Port,
	})
	if err != nil {
		return nil, E.Cause(err, "send scan query on ", a.netInterface.Name)
	}
	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	buffer := make([]byte, anchor.MaxResponseSize)
	length, addr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		return nil, err
	}
	result.addr = addr
	err = conn.Close()
	if err != nil {
		return nil, E.Cause(err, "close scan conn")
	}

	result.response, err = anchor.ParseResponse(buffer[:length])
	if err != nil {
		return nil, E.Cause(err, "parse response")
	}

	return result, nil
}

type interfaceWithAddr struct {
	netInterface *net.Interface
	addr         net.IP
}

func listInterfaceAddr(logger *log.Logger) (interfaceWithAddrs []interfaceWithAddr) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		logger.Fatal("get interfaces: ", err)
	}

	for _, netInterface := range netInterfaces {
		addrs, err := netInterface.Addrs()
		if err != nil {
			logger.Error("get the address of interface [", netInterface.Name, "]: ", err)
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			default:
				continue
			}
			if ip.To4() == nil || ip.IsLoopback() || ip.IsLinkLocalUnicast() {
				continue
			}
			interfaceWithAddrs = append(interfaceWithAddrs, interfaceWithAddr{
				netInterface: &netInterface,
				addr:         ip,
			})
			break
		}
	}

	return
}
