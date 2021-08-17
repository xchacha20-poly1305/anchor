package api

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/pkg/errors"
	"github.com/ulikunitz/xz"
	"io/ioutil"
)

const Version = 0

type Query struct {
	Version    uint8
	DeviceName string
}

type Response struct {
	Version    uint8
	SocksPort  uint16
	DnsPort    uint16
	DeviceName string
	Debug      bool
	BypassLan  bool
}

func MakeQuery(query Query) ([]byte, error) {
	buf := &bytes.Buffer{}
	writer, err := xz.NewWriter(buf)
	if err == nil {
		err = binary.Write(writer, binary.LittleEndian, query.Version)
	}
	if err == nil {
		strArr := []byte(query.DeviceName)
		if len(strArr) > 128 {
			strArr = strArr[:128]
		}
		err = binary.Write(writer, binary.LittleEndian, uint8(len(strArr)))
		if err == nil {
			_, err = writer.Write(strArr)
		}
	}
	if err != nil {
		return nil, errors.WithMessage(err, "write binary error")
	}
	err = writer.Close()
	if err != nil {
		return nil, errors.WithMessage(err, "close writer")
	}

	message, err := ioutil.ReadAll(buf)
	if err != nil {
		return nil, errors.WithMessage(err, "read buf error")
	}

	return message, nil
}

func ParseQuery(message []byte) (*Query, error) {
	query := Query{}
	reader, err := xz.NewReader(bytes.NewReader(message))
	if err == nil {
		err = binary.Read(reader, binary.LittleEndian, &query.Version)
	}
	if query.Version < Version {
		return nil, errors.New(fmt.Sprintf("remote version %d < current version %d, please upgrade your SagerConnect.", query.Version, Version))
	} else if query.Version > Version {
		return nil, errors.New(fmt.Sprintf("remote version %d > current version %d, please upgrade your SagerNet.", query.Version, Version))
	}
	if err == nil {
		var strLen uint8
		err = binary.Read(reader, binary.LittleEndian, &strLen)
		if err == nil {
			strBytes := make([]byte, strLen)
			_, err := reader.Read(strBytes)
			if err == nil {
				query.DeviceName = string(strBytes)
			}
		}
	}
	if err != nil {
		return nil, errors.WithMessage(err, "parse binary error")
	}
	return &query, nil
}

func MakeResponse(response Response) ([]byte, error) {
	buf := &bytes.Buffer{}
	writer, err := xz.NewWriter(buf)
	err = binary.Write(writer, binary.LittleEndian, response.SocksPort)
	if err == nil {
		err = binary.Write(writer, binary.LittleEndian, response.Version)
	}
	if err == nil {
		err = binary.Write(writer, binary.LittleEndian, response.DnsPort)
	}
	if err == nil {
		strArr := []byte(response.DeviceName)
		if len(strArr) > 128 {
			strArr = strArr[:128]
		}
		err = binary.Write(writer, binary.LittleEndian, uint8(len(strArr)))
		if err == nil {
			_, err = writer.Write(strArr)
		}
	}
	if err == nil {
		err = binary.Write(writer, binary.LittleEndian, response.Debug)
	}
	if err == nil {
		err = binary.Write(writer, binary.LittleEndian, response.BypassLan)
	}
	if err != nil {
		return nil, errors.WithMessage(err, "write binary error")
	}
	err = writer.Close()
	if err != nil {
		return nil, errors.WithMessage(err, "close writer")
	}
	message, err := ioutil.ReadAll(buf)
	if err != nil {
		return nil, errors.WithMessage(err, "read buf error")
	}

	return message, nil
}

func ParseResponse(message []byte) (*Response, error) {
	response := Response{}
	reader, err := xz.NewReader(bytes.NewReader(message))
	if err == nil {
		err = binary.Read(reader, binary.LittleEndian, &response.Version)
	}
	if response.Version < Version {
		return nil, errors.New(fmt.Sprintf("remote version %d < current version %d, please upgrade your SagerNet.", response.Version, Version))
	} else if response.Version > Version {
		return nil, errors.New(fmt.Sprintf("remote version %d > current version %d, please upgrade your SagerConnect.", response.Version, Version))
	}
	if err == nil {
		err = binary.Read(reader, binary.LittleEndian, &response.SocksPort)
	}
	if err == nil {
		err = binary.Read(reader, binary.LittleEndian, &response.DnsPort)
	}
	var strLen uint8
	err = binary.Read(reader, binary.LittleEndian, &strLen)
	if err == nil {
		strBytes := make([]byte, strLen)
		_, err := reader.Read(strBytes)
		if err == nil {
			response.DeviceName = string(strBytes)
		}
	}
	if err == nil {
		err = binary.Read(reader, binary.LittleEndian, &response.Debug)
	}
	if err == nil {
		err = binary.Read(reader, binary.LittleEndian, &response.BypassLan)
	}
	if err != nil {
		return nil, errors.WithMessage(err, "parse binary error")
	}
	return &response, nil
}
