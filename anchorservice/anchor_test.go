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
	service := New(context.Background(), logger.NOP(), response, nil)

	listenAddr := &net.UDPAddr{
		IP: net.IP{127, 0, 0, 0},
	}
	packetConn := common.Must1(net.ListenUDP(N.NetworkUDP, listenAddr)) // For test, we don't set standard anchor port to listen.
	defer packetConn.Close()
	err := service.Start(packetConn)
	if err != nil {
		t.Fatalf("start service: %v", err)
	}
	defer service.Close()
	err = service.Start(packetConn)
	if err == nil {
		t.Fatal("try start again but not got any error")
	}

	query := anchor.Query{
		Version:    anchor.Version,
		DeviceName: "common",
	}
	const RejectDevice = "rejectable"
	rejectQuery := anchor.Query{
		Version:    anchor.Version,
		DeviceName: RejectDevice,
	}
	testQuery := func(name string, wantReject bool) {
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

		resp, err := connect(common.Must1(query.MarshalBinary()), packetConn.LocalAddr())
		if err != nil {
			t.Errorf("[%s] failed: %v", name, err)
			return
		}
		if !reflect.DeepEqual(resp, response) {
			t.Errorf("[%s] not same", name)
			return
		}

		rejectResp, err := connect(common.Must1(rejectQuery.MarshalBinary()), packetConn.LocalAddr())
		if wantReject {
			if err == nil {
				t.Errorf("[%s] wants reject but failed: %v", name, err)
			}
			return
		}
		if err != nil {
			t.Errorf("[%s] reject failed: %v", name, err)
			return
		}
		if !reflect.DeepEqual(rejectResp, response) {
			t.Errorf("[%s] reject not same", name)
			return
		}
	}
	testQuery("No Reject", false)

	service.UpdateResponse(response, func(_ net.Addr, deviceName string) bool {
		return deviceName == RejectDevice
	})
	testQuery("Reject", true)

	err = service.Close()
	if err != nil {
		t.Fatalf("close: %v", err)
	}

	packetConn = common.Must1(net.ListenUDP(N.NetworkUDP, listenAddr))
	err = service.Start(packetConn)
	if err != nil {
		t.Fatalf("try reuse instance: %v", err)
	}
	testQuery("Reject + reuse", true)
	err = service.Close()
	if err != nil {
		t.Fatalf("close again: %v", err)
	}
}
