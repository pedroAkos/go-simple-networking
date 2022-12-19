package neti

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

// EncodeString encodes a string to a byte array
func EncodeString(s string) ([]byte, error) {
	b := new(bytes.Buffer)
	err := EncodeStringToBuffer(s, b)
	if err != nil {
		panic(err)
		return nil, err
	}
	return b.Bytes(), nil
}

// DecodeString decodes a string from a byte array
func DecodeString(buf []byte) (string, error) {
	b := bytes.NewBuffer(buf)
	return DecodeStringFromBuffer(b)
}

// EncodeStringToBuffer encodes a string to a buffer
func EncodeStringToBuffer(s string, buffer *bytes.Buffer) error {
	sb := []byte(s)
	return EncodeBytesToBuffer(sb, buffer)
}

// DecodeStringFromBuffer decodes a string from a buffer
func DecodeStringFromBuffer(buffer *bytes.Buffer) (string, error) {
	sb, err := DecodeBytesFromBuffer(buffer)
	if err != nil {
		panic(err)
		return "", err
	}
	return string(sb), nil
}

// EncodeBytesToBuffer encodes a byte array to a buffer
func EncodeBytesToBuffer(b []byte, buffer *bytes.Buffer) error {
	err := binary.Write(buffer, binary.BigEndian, uint16(len(b)))
	if err != nil {
		panic(err)
		return err
	}
	n, err := buffer.Write(b)
	if err != nil {
		panic(err)
		return err
	}
	if n != len(b) {
		err = errors.New(fmt.Sprint("Expected to write ", len(b), " wrote ", n))
		panic(err)
		return err
	}
	return nil
}

// DecodeBytesFromBuffer decodes a byte array from a buffer
func DecodeBytesFromBuffer(buffer *bytes.Buffer) ([]byte, error) {
	var bLen uint16
	err := binary.Read(buffer, binary.BigEndian, &bLen)
	if err != nil {
		panic(err)
		return nil, err
	}

	b := make([]byte, bLen)
	n, err := buffer.Read(b)
	if err != nil {
		panic(err)
		return nil, err
	}
	if n != int(bLen) {
		err = errors.New(fmt.Sprint("Expected to read ", bLen, " read ", n))
		panic(err)
		return nil, err
	}

	return b, nil
}

// EncodeNumberToBuffer encodes a number to a buffer
func EncodeNumberToBuffer(n interface{}, buffer *bytes.Buffer) error {
	if err := binary.Write(buffer, binary.BigEndian, n); err != nil {
		panic(err)
	}
	return nil
}

// DecodeNumberFromBuffer decodes a number from a buffer
// Note that nPointer must be a pointer to the type of the number
func DecodeNumberFromBuffer(nPointer interface{}, buffer *bytes.Buffer) error {
	if err := binary.Read(buffer, binary.BigEndian, nPointer); err != nil {
		panic(err)
	}
	return nil
}
