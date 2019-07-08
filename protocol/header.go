package protocol

import "encoding/binary"

// Header section format [RFC 1035 4.1.1]
//
// The header contains the following fields:
//
//                                     1  1  1  1  1  1
//       0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                      ID                       |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |QR|   Opcode  |AA|TC|RD|RA|   Z    |   RCODE   |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                    QDCOUNT                    |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                    ANCOUNT                    |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                    NSCOUNT                    |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                    ARCOUNT                    |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
type Header struct {
	ID                 uint16       // ID       16
	MessageType        MessageType  // QR        1
	OpCode             OpCode       // OPCODE    4
	Authoritative      bool         // AA        1
	Truncation         bool         // TC        1
	RecursionDesired   bool         // RD        1
	RecursionAvailable bool         // RA        1
	Reserved           uint8        // Z         3
	ResponseCode       ResponseCode // RCODE     4
	QuestionCount      uint16       // QDCOUNT  16
	AnswerCount        uint16       // ANCOUNT  16
	NameServerCount    uint16       // NSCOUNT  16
	AdditionalCount    uint16       // ARCOUNT  16
}

func (header Header) Encode() []byte {
	data := make([]byte, 12)
	binary.BigEndian.PutUint16(data[0:2], header.ID)
	data[2] |= uint8(header.MessageType) << 7
	if header.Authoritative {
		data[2] |= 1 << 2
	}
	if header.Truncation {
		data[2] |= 1 << 1
	}
	if header.RecursionDesired {
		data[2] |= 1 << 0
	}
	if header.RecursionAvailable {
		data[3] |= 1 << 7
	}
	data[2] |= (uint8(header.OpCode) & 0xF) << 3
	data[3] |= uint8(header.ResponseCode) & 0xF
	binary.BigEndian.PutUint16(data[4:6], header.QuestionCount)
	binary.BigEndian.PutUint16(data[6:8], header.AnswerCount)
	binary.BigEndian.PutUint16(data[8:10], header.NameServerCount)
	binary.BigEndian.PutUint16(data[10:12], header.AdditionalCount)
	return data
}

func (header *Header) Decode(data []byte, off int) (Header, int) {
	*header = Header{
		ID:                 binary.BigEndian.Uint16(data[off : off+2]),
		MessageType:        MessageType((data[off+2] >> 7) & 0x01),
		OpCode:             OpCode((data[off+2] >> 3) & 0x0F),
		Authoritative:      data[off+2]&0x04 != 0,
		Truncation:         data[off+2]&0x02 != 0,
		RecursionDesired:   data[off+2]&0x01 != 0,
		RecursionAvailable: data[off+3]&0x80 != 0,
		Reserved:           0,
		ResponseCode:       ResponseCode(data[off+3] & 0x0F),
		QuestionCount:      binary.BigEndian.Uint16(data[off+4 : off+6]),
		AnswerCount:        binary.BigEndian.Uint16(data[off+6 : off+8]),
		NameServerCount:    binary.BigEndian.Uint16(data[off+8 : off+10]),
		AdditionalCount:    binary.BigEndian.Uint16(data[off+10 : off+12]),
	}
	return *header, off + 12
}
