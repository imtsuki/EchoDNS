package protocol

type Message struct {
	Header      Header
	Questions   []Question
	Answers     []Resource
	NameServers []Resource
	Additionals []Resource
	//rawPacket   NoDisplayPacket
}

type NoDisplayPacket []byte

func (p NoDisplayPacket) String() string {
	return ""
}

func (m *Message) Decode(data []byte, off int) (Message, int) {
	//m.rawPacket = NoDisplayPacket(data)
	_, off = m.Header.Decode(data, off)

	for i := uint16(0); i < m.Header.QuestionCount; i++ {
		question := Question{}
		_, off = question.Decode(data, off)
		m.Questions = append(m.Questions, question)
	}
	return *m, off
}
