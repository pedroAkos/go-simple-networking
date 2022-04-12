package neti

import "bytes"

type MessageWrap struct {
	Id  string
	Msg Message
}

func (m MessageWrap) String() string {
	return m.Msg.String()
}

func (m MessageWrap) Name() string {
	return m.Msg.Name()
}

func (m MessageWrap) Code() uint16 {
	return m.Msg.Code()
}

func (m MessageWrap) Serialize(buff *bytes.Buffer) error {
	_ = EncodeStringToBuffer(m.Id, buff)
	if err := m.Msg.Serialize(buff); err != nil {
		panic(err)
		return err
	}

	return nil
}

func (m MessageWrap) Deserialize(buff *bytes.Buffer) (Message, error) {
	m.Id, _ = DecodeStringFromBuffer(buff)
	msg, err := m.Msg.Deserialize(buff)
	if err != nil {
		panic(err)
		return nil, err
	}
	m.Msg = msg
	return m, nil
}
