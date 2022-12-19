package neti

import (
	"bytes"
	"fmt"
)

// MessageWrap is a wrapper for messages that are sent over the network.
type MessageWrap struct {
	Id   string
	Msg  Message
	code uint16
	buff *bytes.Buffer
}

// String returns a string representation of the message.
func (m MessageWrap) String() string {
	return fmt.Sprintf("%v{cid: %v mcode: %v}", m.Name(), m.Id, m.code)
}

// Name returns the name of the message.
func (m MessageWrap) Name() string {
	return "MessageWrap"
}

// Code returns the code of the message.
func (m MessageWrap) Code() uint16 {
	return 0
}

// Buff returns the buffer of the message.
func (m MessageWrap) Buff() *bytes.Buffer {
	return m.buff
}

// MessageCode is a code for a message.
func (m MessageWrap) MessageCode() uint16 {
	return m.code
}

// Serialize serializes the message.
func (m MessageWrap) Serialize(buff *bytes.Buffer) error {
	_ = EncodeStringToBuffer(m.Id, buff)
	_ = EncodeNumberToBuffer(m.Msg.Code(), buff)
	if err := m.Msg.Serialize(buff); err != nil {
		panic(err)
		return err
	}

	return nil
}

// Deserialize deserializes the message.
func (m MessageWrap) Deserialize(buff *bytes.Buffer) (Message, error) {
	m.Id, _ = DecodeStringFromBuffer(buff)
	_ = DecodeNumberFromBuffer(&m.code, buff)
	//msg, err := m.Msg.Deserialize(buff)
	//if err != nil {
	//	panic(err)
	//	return nil, err
	//}
	//m.Msg = msg
	m.buff = buff
	return m, nil
}
