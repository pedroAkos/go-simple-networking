package neti

import (
	"bytes"
	"fmt"
)

type MessageWrap struct {
	Id   string
	Msg  Message
	code uint16
	buff *bytes.Buffer
}

func (m MessageWrap) String() string {
	return fmt.Sprintf("%v{cid: %v mcode: %v}", m.Name(), m.Id, m.code)
}

func (m MessageWrap) Name() string {
	return "MessageWrap"
}

func (m MessageWrap) Code() uint16 {
	return 0
}

func (m MessageWrap) Buff() *bytes.Buffer {
	return m.buff
}

func (m MessageWrap) MessageCode() uint16 {
	return m.code
}

func (m MessageWrap) Serialize(buff *bytes.Buffer) error {
	_ = EncodeStringToBuffer(m.Id, buff)
	_ = EncodeNumberToBuffer(m.Msg.Code(), buff)
	if err := m.Msg.Serialize(buff); err != nil {
		panic(err)
		return err
	}

	return nil
}

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
