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
	bytes := make([]byte, 12)
	binary.BigEndian.PutUint16(bytes[0:2], header.ID)
	bytes[2] |= uint8(header.MessageType) << 7
	if header.Authoritative {
		bytes[2] |= 1 << 2
	}
	if header.Truncation {
		bytes[2] |= 1 << 1
	}
	if header.RecursionDesired {
		bytes[2] |= 1 << 0
	}
	if header.RecursionAvailable {
		bytes[3] |= 1 << 7
	}
	bytes[2] |= (uint8(header.OpCode) & 0xF) << 3
	bytes[3] |= uint8(header.ResponseCode) & 0xF
	binary.BigEndian.PutUint16(bytes[4:6], header.QuestionCount)
	binary.BigEndian.PutUint16(bytes[6:8], header.AnswerCount)
	binary.BigEndian.PutUint16(bytes[8:10], header.NameServerCount)
	binary.BigEndian.PutUint16(bytes[10:12], header.AdditionalCount)
	return bytes
}
