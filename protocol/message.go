package protocol

type Message struct {
	Header    Header
	Questions []Question
	RawPacket []byte
}
