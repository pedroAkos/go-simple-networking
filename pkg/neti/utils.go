package neti

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
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
	err := binary.Write(buffer, binary.BigEndian, uint16(len(sb)))
	if err != nil {
		return err
	}
	n, err := buffer.Write(sb)
	if n != len(sb) {
		return errors.New(fmt.Sprint("Expected to write", len(sb), "wrote", n))
	}
	if err != nil {
		return err
	}
	return nil
}

func DecodeStringFromBuffer(buffer *bytes.Buffer) (string, error) {
	var sLen uint16
	err := binary.Read(buffer, binary.BigEndian, &sLen)
	if err != nil {
		return "", err
	}

	sb := make([]byte, sLen)
	n, err := buffer.Read(sb)
	if n != int(sLen) {
		return "", errors.New(fmt.Sprint("Expected to read", sLen, "read", n))
	}

	if err != nil {
		return "", err
	}

	return string(sb), nil
}


func EncodeNumberToBuffer(n interface{}, buffer *bytes.Buffer) error {
	return binary.Write(buffer, binary.BigEndian, n)
}

func DecodeNumberFromBuffer(nPointer interface{}, buffer *bytes.Buffer) error {
	return binary.Read(buffer, binary.BigEndian, nPointer)
}