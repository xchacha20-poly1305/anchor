package anchor

import (
	"bytes"
	"testing"
)

func TestData(t *testing.T) {
	const deviceName = "my device"
	q := &Query{
		Version:    Version,
		DeviceName: deviceName,
	}
	m, _ := q.MarshalBinary()
	q, err := ParseQuery(bytes.NewReader(m))
	if err != nil {
		t.Fatal(err)
	}
	if q.Version != Version || q.DeviceName != deviceName {
		t.Fatal("err parse query")
	}
	r := &Response{
		Version:    Version,
		DnsPort:    6450,
		DeviceName: deviceName,
		SocksPort:  2080,
	}
	m, err = r.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	parsedResponse, err := ParseResponse(bytes.NewReader(m))
	if err != nil {
		t.Fatal(err)
	}
	if parsedResponse.Version != Version ||
		parsedResponse.SocksPort != r.SocksPort ||
		parsedResponse.DnsPort != r.DnsPort ||
		parsedResponse.DeviceName != deviceName {
		t.Fatal("err parse response")
	}
}
