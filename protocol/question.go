package protocol

import "encoding/binary"

// Question section format [RFC 1035 4.1.2]
//
// The question section is used to carry the "question" in most queries,
// i.e., the parameters that define what is being asked.  The section
// contains QDCOUNT (usually 1) entries, each of the following format:
//
//                                     1  1  1  1  1  1
//       0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                                               |
//     /                     QNAME                     /
//     /                                               /
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                     QTYPE                     |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                     QCLASS                    |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
type Question struct {
	Name  Name // QNAME    *
	Type  Type   // QTYPE   16
	Class Class  // OCLASS  16
}

func (question Question) Encode() []byte {
	data := make([]byte, 0)

	data = append(data, question.Name.Encode()...)

	var ending [4]byte
	binary.BigEndian.PutUint16(ending[0:2], uint16(question.Type))
	binary.BigEndian.PutUint16(ending[2:4], uint16(question.Class))

	data = append(data, ending[:]...)
	return data
}

func (question *Question) Decode(data []byte, off int) (Question, int) {
	_, off = question.Name.Decode(data, off)
	question.Type = Type(binary.BigEndian.Uint16(data[off:off+2]))
	question.Class = Class(binary.BigEndian.Uint16(data[off+2:off+4]))
	off += 4
	return *question, off
}
