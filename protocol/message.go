package protocol

type Encodable interface {
	Encode() []byte
}

type Message struct {
	Header    Header
	RawPacket []byte
}
