package neti

import "bytes"

type messageWrap struct {
	id string
	msg Message
}

func (m messageWrap) String() string {
	return m.msg.String()
}

func (m messageWrap) Name() string {
	return m.msg.Name()
}

func (m messageWrap) Code() uint16 {
	return m.msg.Code()
}

func (m messageWrap) Serialize() ([]byte, error) {
	buff := new(bytes.Buffer)
	if err := EncodeStringToBuffer(m.id, buff); err != nil {
		return nil, err
	}
	if b, err := m.msg.Serialize(); err != nil {
		return nil, err
	} else if n, err :=  buff.Write(b); n != len(b) || err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

func (m messageWrap) Deserialize(b []byte) (Message, error) {
	buff := bytes.NewBuffer(b)
	var err error
	if m.id, err = DecodeStringFromBuffer(buff); err != nil {
		return nil, err
	}
	msg, err := m.msg.Deserialize(buff.Bytes())
	if err != nil {
		return nil, err
	}
	m.msg = msg
	return m, nil
}
