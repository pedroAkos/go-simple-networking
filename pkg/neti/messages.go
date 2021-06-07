package neti


type messageWrap struct {
	id uint16
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
	return m.msg.Serialize()
}

func (m messageWrap) Deserialize(bytes []byte) (Message, error) {
	msg, err := m.msg.Deserialize(bytes)
	if err != nil {
		return nil, err
	}
	m.msg = msg
	return m, nil
}
