package protocol

type Encodable interface {
	Encode() []byte
}
