// Package anchor provides an API
// that allow switching the proxy info.
package anchor

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"math"

	"github.com/sagernet/sing/common"
	"github.com/sagernet/sing/common/auth"
	E "github.com/sagernet/sing/common/exceptions"
)

const (
	Version    = 0x01
	ListenPort = 45947

	MaxDeviceName = 128
	MaxAuth       = math.MaxUint8 // RFC 1929

	MinQuerySize    = 1 + 1 + 0                                       // 3
	MaxQuerySize    = 1 + 1 + 128                                     // 130
	MinResponseSize = 1 + 2 + 1 + 0 + 1 + 1 + 2 + 1 + 0 + 1 + 0       // 10
	MaxResponseSize = 1 + 2 + 1 + 128 + 1 + 1 + 2 + 1 + 255 + 1 + 255 // 648
)

var _ encoding.BinaryMarshaler = (*Query)(nil)

type Query struct {
	Version    uint8
	DeviceName string
}

func ParseQuery(message []byte) (*Query, error) {
	query := &Query{}
	reader := bytes.NewReader(message)
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

func (q *Query) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, MinQuerySize))
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
	User       auth.User
}

func ParseResponse(message []byte) (*Response, error) {
	response := &Response{}
	reader := bytes.NewReader(message)
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

	var userLen uint8
	err = binary.Read(reader, binary.BigEndian, &userLen)
	if err != nil {
		return nil, E.Cause(err, "read socks user length")
	}
	if userLen > 0 {
		user := make([]byte, userLen)
		_, err = reader.Read(user)
		if err != nil {
			return nil, E.Cause(err, "read socks user")
		}
		response.User.Username = string(user)
	}

	var passwordLen uint8
	err = binary.Read(reader, binary.BigEndian, &passwordLen)
	if err != nil {
		return nil, E.Cause(err, "read password length")
	}
	if passwordLen > 0 {
		password := make([]byte, passwordLen)
		_, err = reader.Read(password)
		if err != nil {
			return nil, E.Cause(err, "read password")
		}
		response.User.Password = string(password)
	}

	return response, nil
}

func (r *Response) MarshalBinary() ([]byte, error) {
	if len(r.User.Username) > MaxAuth || len(r.User.Password) > MaxAuth {
		return nil, E.New("invalid auth length")
	}

	buf := bytes.NewBuffer(make([]byte, 0, MinResponseSize))
	common.Must(binary.Write(buf, binary.BigEndian, r.Version))
	common.Must(binary.Write(buf, binary.BigEndian, r.DnsPort))
	deviceName := []byte(r.DeviceName)
	if len(deviceName) > 128 {
		deviceName = deviceName[:128]
	}
	common.Must(binary.Write(buf, binary.BigEndian, uint8(len(deviceName))))
	common.Must1(buf.Write(deviceName))
	common.Must(binary.Write(buf, binary.BigEndian, r.SocksPort))
	common.Must(binary.Write(buf, binary.BigEndian, uint8(len(r.User.Username))))
	common.Must1(buf.Write([]byte(r.User.Username)))
	common.Must(binary.Write(buf, binary.BigEndian, uint8(len(r.User.Password))))
	common.Must1(buf.Write([]byte(r.User.Password)))

	return buf.Bytes(), nil
}
