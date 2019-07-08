package protocol

type Message struct {
	Header      Header
	Questions   []Question
	Answers     []Resource
	NameServers []Resource
	Additionals []Resource
	RawPacket   []byte
}

func (m *Message) Decode(data []byte, off int) (Message, int) {
	m.RawPacket = data
	_, off = m.Header.Decode(data, off)

	for i := uint16(0); i < m.Header.QuestionCount; i++ {
		question := Question{}
		_, off = question.Decode(data, off)
		m.Questions = append(m.Questions, question)
	}
	return *m, off
}

func (m Message) Encode() []byte {
	data := make([]byte, 0)

	m.Header.QuestionCount = uint16(len(m.Questions))
	m.Header.AnswerCount = uint16(len(m.Answers))
	m.Header.NameServerCount = uint16(len(m.NameServers))
	m.Header.AdditionalCount = uint16(len(m.Additionals))

	data = append(data, m.Header.Encode()...)

	for _, q := range m.Questions {
		data = append(data, q.Encode()...)
	}
	for _, a := range m.Answers {
		data = append(data, a.Encode()...)
	}
	for _, n := range m.NameServers {
		data = append(data, n.Encode()...)
	}
	for _, a := range m.Additionals {
		data = append(data, a.Encode()...)
	}

	return data
}
