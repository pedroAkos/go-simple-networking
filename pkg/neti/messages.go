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

func (m MessageWrap) Serialize() ([]byte, error) {
	buff := new(bytes.Buffer)
	if err := EncodeStringToBuffer(m.Id, buff); err != nil {
		return nil, err
	}
	if b, err := m.Msg.Serialize(); err != nil {
		return nil, err
	} else if n, err :=  buff.Write(b); n != len(b) || err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

func (m MessageWrap) Deserialize(b []byte) (Message, error) {
	buff := bytes.NewBuffer(b)
	var err error
	if m.Id, err = DecodeStringFromBuffer(buff); err != nil {
		return nil, err
	}
	msg, err := m.Msg.Deserialize(buff.Bytes())
	if err != nil {
		return nil, err
	}
	m.Msg = msg
	return m, nil
}
