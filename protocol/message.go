package protocol

type Message struct {
	Header      Header
	Questions   []Question
	Answers     []Resource
	NameServers []Resource
	Additionals []Resource
	RawPacket   []byte
}
