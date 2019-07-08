package protocol

import "encoding/binary"
import "net"

// Resource record format [RFC 1035 4.1.2]
//
// The answer, authority, and additional sections all share the same
// format: a variable number of resource records, where the number of
// records is specified in the corresponding count field in the header.
// Each resource record has the following format:
//
//                                     1  1  1  1  1  1
//       0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                                               |
//     /                                               /
//     /                      NAME                     /
//     |                                               |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                      TYPE                     |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                     CLASS                     |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                      TTL                      |
//     |                                               |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                   RDLENGTH                    |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--|
//     /                     RDATA                     /
//     /                                               /
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
type Resource struct {
	Name   Name
	Type   Type
	Class  Class
	TTL    uint32
	Length uint16
	Data   ResourceData
}

type ResourceData interface {
	Encode() []byte
	ResourceType() Type
}

func (r Resource) Encode() []byte {
	name := r.Name.Encode()
	data := r.Data.Encode()
	r.Length = uint16(len(data))
	fields := make([]byte, 10)
	binary.BigEndian.PutUint16(fields[0:2], uint16(r.Type))
	binary.BigEndian.PutUint16(fields[2:4], uint16(r.Class))
	binary.BigEndian.PutUint32(fields[4:8], uint32(r.TTL))
	binary.BigEndian.PutUint16(fields[8:10], uint16(r.Length))
	return append(append(name, fields...), data...)
}

type AResource struct {
	IP net.IP
}

func (r *AResource) ResourceType() Type {
	return TypeA
}

func (r *AResource) Encode() []byte {
	return r.IP.To4()
}

type AAAAResource struct {
	IP net.IP
}

func (r *AAAAResource) ResourceType() Type {
	return TypeAAAA
}

func (r *AAAAResource) Encode() []byte {
	return r.IP.To16()
}

type ALLResource struct {
	Data []byte
}

func (r *ALLResource) ResourceType() Type {
	return TypeALL

}
func (r *ALLResource) Encode() []byte {
	return r.Data
}
