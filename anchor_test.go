package anchor

import (
	"testing"
)

func TestData(t *testing.T) {
	const deviceName = "my device"
	q := &Query{
		Version:    Version,
		DeviceName: deviceName,
	}
	m, _ := q.MarshalBinary()
	q, err := ParseQuery(m)
	if err != nil {
		t.Fatal(err)
	}
	if q.Version != Version || q.DeviceName != deviceName {
		t.Fatal("err parse query")
	}
	r := &Response{
		Version:       Version,
		DnsPort:       6450,
		DeviceName:    deviceName,
		SocksPort:     2080,
		SocksUser:     "invalid",
		SocksPassword: "sekai",
	}
	m, err = r.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	parsedResponse, err := ParseResponse(m)
	if err != nil {
		t.Fatal(err)
	}
	if parsedResponse.Version != Version ||
		parsedResponse.SocksPort != r.SocksPort ||
		parsedResponse.DnsPort != r.DnsPort ||
		parsedResponse.DeviceName != deviceName ||
		parsedResponse.SocksUser != r.SocksUser ||
		parsedResponse.SocksPassword != r.SocksPassword {
		t.Fatal("err parse response")
	}
}
