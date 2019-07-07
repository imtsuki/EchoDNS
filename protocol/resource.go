package protocol

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
	Name   string
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

type AResource struct {
	IP [4]byte
}

func (r *AResource) ResourceType() Type {
	return TypeA
}

func (r *AResource) Encode() []byte {
	return r.IP[:]
}

type AAAAResource struct {
}

func (r *AAAAResource) ResourceType() Type {
	return TypeAAAA
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
