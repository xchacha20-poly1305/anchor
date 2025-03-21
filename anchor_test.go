package anchor

import (
	"bytes"
	"testing"

	"github.com/sagernet/sing/common/buf"
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
		buf.Put(m)
		t.Fatal(err)
	}
	if q.Version != Version || q.DeviceName != deviceName {
		buf.Put(m)
		t.Fatal("err parse query")
	}
	r := &Response{
		Version:    Version,
		DnsPort:    6450,
		DeviceName: deviceName,
		SocksPort:  2080,
	}
	buf.Put(m)
	m, err = r.MarshalBinary()
	if err != nil {
		buf.Put(m)
		t.Fatal(err)
	}
	parsedResponse, err := ParseResponse(bytes.NewReader(m))
	if err != nil {
		buf.Put(m)
		t.Fatal(err)
	}
	if parsedResponse.Version != Version ||
		parsedResponse.SocksPort != r.SocksPort ||
		parsedResponse.DnsPort != r.DnsPort ||
		parsedResponse.DeviceName != deviceName {
		buf.Put(m)
		t.Fatal("err parse response")
	}
	buf.Put(m)
}
