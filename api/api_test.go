package api

import (
	"testing"
)

func TestData(t *testing.T) {
	d := "my device"
	q := &Query{
		Version:    Version,
		DeviceName: d,
	}
	m, err := MakeQuery(q)
	if err != nil {
		t.Fatal(err)
	}
	q, err = ParseQuery(m)
	if err != nil {
		t.Fatal(err)
	}
	if q.Version != Version || q.DeviceName != d {
		t.Fatal("err parse query")
	}
	r := &Response{
		Version:    Version,
		SocksPort:  114,
		DnsPort:    514,
		DeviceName: d,
		Debug:      false,
		BypassLan:  true,
	}
	m, err = MakeResponse(r)
	if err != nil {
		t.Fatal(err)
	}
	r, err = ParseResponse(m)
	if err != nil {
		t.Fatal(err)
	}
	if r.Version != Version || r.SocksPort != 114 || r.DnsPort != 514 || r.DeviceName != d || r.Debug || !r.BypassLan {
		t.Fatal("err parse response")
	}
}
