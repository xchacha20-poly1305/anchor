package anchorservice

import (
	"context"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/sagernet/sing/common"
	"github.com/sagernet/sing/common/buf"
	"github.com/sagernet/sing/common/logger"
	N "github.com/sagernet/sing/common/network"

	"github.com/xchacha20-poly1305/anchor"
)

func Test_Service(t *testing.T) {
	response := &anchor.Response{
		Version:    anchor.Version,
		DnsPort:    53,
		DeviceName: "test_service",
		SocksPort:  1080,
	}
	listenAddr := &net.UDPAddr{
		IP:   net.IP{127, 0, 0, 0},
		Port: anchor.Port,
	}
	const RejectDevice = "rejectable"
	service := New(context.Background(), logger.NOP(), listenAddr, response, func(_ net.Addr, deviceName string) bool {
		return deviceName == RejectDevice
	})

	err := service.Start()
	if err != nil {
		t.Fatalf("start service: %v", err)
	}
	defer service.Close()
	err = service.Start()
	if err == nil {
		t.Fatal("try start again but not got any error")
	}

	connect := func(query []byte, target net.Addr) (*anchor.Response, error) {
		localConn := common.Must1(net.ListenUDP(N.NetworkUDP, nil))
		defer localConn.Close()
		common.Must1(localConn.WriteTo(query, target))
		const waitTimeout = 5 * time.Second
		common.Must(localConn.SetReadDeadline(time.Now().Add(waitTimeout)))
		buffer := buf.NewSize(anchor.MaxResponseSize)
		defer buffer.Release()
		_, _, err := buffer.ReadPacketFrom(localConn)
		if err != nil {
			return nil, err
		}
		return anchor.ParseResponse(buffer)
	}

	query := anchor.Query{
		Version:    anchor.Version,
		DeviceName: "common",
	}
	rejectQuery := anchor.Query{
		Version:    anchor.Version,
		DeviceName: RejectDevice,
	}
	resp, err := connect(common.Must1(query.MarshalBinary()), listenAddr)
	if err != nil {
		t.Fatalf("try send query: %v", err)
	}
	if !reflect.DeepEqual(resp, response) {
		t.Fatalf("receive invalid response")
	}
	_, err = connect(common.Must1(rejectQuery.MarshalBinary()), listenAddr)
	if err == nil {
		t.Fatalf("want reject but not")
	}

	err = service.Close()
	if err != nil {
		t.Fatalf("close: %v", err)
	}
}
