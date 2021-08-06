package neti

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
)

func EncodeString(s string) ([]byte, error) {
	b := new(bytes.Buffer)
	err := EncodeStringToBuffer(s, b)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func DecodeString(buf []byte) (string, error) {
	b := bytes.NewBuffer(buf)
	return DecodeStringFromBuffer(b)
}


func EncodeStringToBuffer(s string, buffer *bytes.Buffer) error {
	sb := []byte(s)
	return EncodeBytesToBuffer(sb, buffer)
}

func DecodeStringFromBuffer(buffer *bytes.Buffer) (string, error) {
	sb, err := DecodeBytesFromBuffer(buffer)
	if err != nil {
		return "", err
	}
	return string(sb), nil
}


func EncodeBytesToBuffer(b []byte, buffer *bytes.Buffer) error {
	err := binary.Write(buffer, binary.BigEndian, uint16(len(b)))
	if err != nil {
		return err
	}
	n, err := buffer.Write(b)
	if n != len(b) {
		return errors.New(fmt.Sprint("Expected to write ", len(b), " wrote ", n))
	}
	if err != nil {
		return err
	}
	return nil
}

func DecodeBytesFromBuffer(buffer *bytes.Buffer) ([]byte, error) {
	var bLen uint16
	err := binary.Read(buffer, binary.BigEndian, &bLen)
	if err != nil {
		return nil, err
	}

	b := make([]byte, bLen)
	n, err := buffer.Read(b)
	if n != int(bLen) {
		log.Panic("Expected to read ", bLen, " read ", n)
		return nil, errors.New(fmt.Sprint("Expected to read ", bLen, " read ", n))
	}
	if err != nil {
		return nil, err
	}

	return b, nil
}


func EncodeNumberToBuffer(n interface{}, buffer *bytes.Buffer) error {
	return binary.Write(buffer, binary.BigEndian, n)
}

func DecodeNumberFromBuffer(nPointer interface{}, buffer *bytes.Buffer) error {
	return binary.Read(buffer, binary.BigEndian, nPointer)
}