package protocol

type Message struct {
	Header    Header
	RawPacket []byte
}
