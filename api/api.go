package api

import (
	"bytes"
	"encoding/binary"
	"github.com/pkg/errors"
	"github.com/ulikunitz/xz"
	"io/ioutil"
)

type Query struct {
	DeviceName string
}

type Response struct {
	SocksPort  uint16
	DnsPort    uint16
	DeviceName string
	Debug      bool
}

func MakeQuery(query Query) ([]byte, error) {
	buf := &bytes.Buffer{}
	writer, err := xz.NewWriter(buf)
	strArr := []byte(query.DeviceName)
	if len(strArr) > 128 {
		strArr = strArr[:128]
	}
	err = binary.Write(writer, binary.LittleEndian, uint8(len(strArr)))
	if err == nil {
		_, err = writer.Write(strArr)
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
	if err != nil {
		return nil, errors.WithMessage(err, "read message error")
	}
	var strLen uint8
	err = binary.Read(reader, binary.LittleEndian, &strLen)
	if err == nil {
		strBytes := make([]byte, strLen)
		_, err := reader.Read(strBytes)
		if err == nil {
			query.DeviceName = string(strBytes)
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
		err = binary.Write(writer, binary.LittleEndian, response.DnsPort)
	}
	strArr := []byte(response.DeviceName)
	if len(strArr) > 128 {
		strArr = strArr[:128]
	}
	err = binary.Write(writer, binary.LittleEndian, uint8(len(strArr)))
	if err == nil {
		_, err = writer.Write(strArr)
	}
	if err == nil {
		err = binary.Write(writer, binary.LittleEndian, response.Debug)
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
	if err != nil {
		return nil, errors.WithMessage(err, "read message error")
	}
	err = binary.Read(reader, binary.LittleEndian, &response.SocksPort)
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
	if err != nil {
		return nil, errors.WithMessage(err, "parse binary error")
	}
	return &response, nil
}
