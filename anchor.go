// Package anchor provides an API
// that allow switching the proxy info.
package anchor

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"io"

	"github.com/sagernet/sing/common"
	E "github.com/sagernet/sing/common/exceptions"
)

const (
	Version = 0x02
	Port    = 45947

	MaxDeviceName   = 128
	MinQuerySize    = 1 + 1 + 0                     // 2
	MaxQuerySize    = 1 + 1 + 128                   // 130
	MinResponseSize = 1 + 2 + 1 + 0 + 2             // 6
	MaxResponseSize = 1 + 2 + 1 + MaxDeviceName + 2 // 134
)

var _ encoding.BinaryMarshaler = (*Query)(nil)

type Query struct {
	Version    uint8
	DeviceName string
}

func ParseQuery(reader io.Reader) (*Query, error) {
	query := &Query{}
	err := binary.Read(reader, binary.BigEndian, &query.Version)
	if err != nil {
		return nil, E.Cause(err, "read version")
	}
	if query.Version < Version {
		return nil, E.New("remote version: ", query.Version, " is less than current version")
	} else if query.Version > Version {
		return nil, E.New("remote version: ", query.Version, " is greater than current version")
	}
	var strLen uint8
	err = binary.Read(reader, binary.BigEndian, &strLen)
	if err != nil {
		return nil, E.Cause(err, "read device name length")
	}
	if strLen > 0 {
		strBytes := make([]byte, strLen)
		_, err = reader.Read(strBytes)
		if err != nil {
			return nil, E.Cause(err, "read device name")
		}
		query.DeviceName = string(strBytes)
	}
	return query, nil
}

func (q Query) Length() (length int) {
	length += 1 // Version
	length += 1 // Device Name Length
	length += len(q.DeviceName)
	return
}

func (q Query) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, q.Length()))
	common.Must(binary.Write(buf, binary.BigEndian, q.Version))
	deviceName := []byte(q.DeviceName)
	if len(deviceName) > MaxDeviceName {
		deviceName = deviceName[:MaxDeviceName]
	}
	common.Must(binary.Write(buf, binary.BigEndian, uint8(len(deviceName))))
	common.Must1(buf.Write(deviceName))
	return buf.Bytes(), nil
}

var _ encoding.BinaryMarshaler = (*Response)(nil)

type Response struct {
	Version    uint8
	DnsPort    uint16
	DeviceName string
	SocksPort  uint16
}

func ParseResponse(reader io.Reader) (*Response, error) {
	response := &Response{}
	err := binary.Read(reader, binary.BigEndian, &response.Version)
	if err != nil {
		return nil, E.Cause(err, "read version")
	}
	if response.Version < Version {
		return nil, E.New("remote version: ", response.Version, " is less than current version")
	} else if response.Version > Version {
		return nil, E.New("remote version: ", response.Version, " is greater than current version")
	}
	err = binary.Read(reader, binary.BigEndian, &response.DnsPort)
	if err != nil {
		return nil, E.Cause(err, "read dns port")
	}
	var deviceNameLen uint8
	err = binary.Read(reader, binary.BigEndian, &deviceNameLen)
	if err != nil {
		return nil, E.Cause(err, "read device name length")
	}
	deviceName := make([]byte, deviceNameLen)
	_, err = reader.Read(deviceName)
	if err != nil {
		return nil, E.Cause(err, "read device name")
	}
	response.DeviceName = string(deviceName)
	err = binary.Read(reader, binary.BigEndian, &response.SocksPort)
	if err != nil {
		return nil, E.Cause(err, "read socks port")
	}

	return response, nil
}

func (r Response) Length() (length int) {
	nameLength := len(r.DeviceName)
	if nameLength > MaxDeviceName {
		nameLength = MaxDeviceName
	}

	length += 1 // Version
	length += 2 // DNS Port
	length += 1 // Device Name Length
	length += nameLength
	length += 2 // Socks Port
	return
}

func (r Response) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, r.Length()))
	common.Must(binary.Write(buf, binary.BigEndian, r.Version))
	common.Must(binary.Write(buf, binary.BigEndian, r.DnsPort))
	deviceName := []byte(r.DeviceName)
	if len(deviceName) > MaxDeviceName {
		deviceName = deviceName[:MaxDeviceName]
	}
	common.Must(binary.Write(buf, binary.BigEndian, uint8(len(deviceName))))
	common.Must1(buf.Write(deviceName))
	common.Must(binary.Write(buf, binary.BigEndian, r.SocksPort))

	return buf.Bytes(), nil
}
