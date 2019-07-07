package protocol

import "encoding/binary"
import "strings"

type OpCode uint8
type MessageType uint8
type ResponseCode uint8
type Type uint16
type Class uint16

const (
	OpCodeQuery  OpCode = 0 // 0: QUERY
	OpCodeIQuery OpCode = 1 // 1: IQUERY
	OpCodeStatus OpCode = 2 // 2: STATUS

	Query    MessageType = 0
	Response MessageType = 1

	NoError  ResponseCode = 0
	FormErr  ResponseCode = 1
	ServFail ResponseCode = 2
	NXDomain ResponseCode = 3
	NotImp   ResponseCode = 4
	Refused  ResponseCode = 5

	// Types for Question and Resource
	TypeA     Type = 1
	TypeNS    Type = 2
	TypeCNAME Type = 5
	TypeSOA   Type = 6
	TypePTR   Type = 12
	TypeMX    Type = 15
	TypeTXT   Type = 16
	TypeAAAA  Type = 28
	TypeSRV   Type = 33

	// Types for Question only
	TypeWKS   Type = 11
	TypeHINFO Type = 13
	TypeMINFO Type = 14
	TypeAXFR  Type = 252
	TypeALL   Type = 255

	// Classes for Question and Resource
	ClassINET   Class = 1
	ClassCSNET  Class = 2
	ClassCHAOS  Class = 3
	ClassHESIOD Class = 4

	// Classes for Question only
	ClassANY Class = 255
)

func (responseCode ResponseCode) String() string {
	switch responseCode {
	case NoError:
		return "NOERROR"
	case FormErr:
		return "FORMERR"
	case ServFail:
		return "SERVFAIL"
	case NXDomain:
		return "NXDOMAIN"
	case NotImp:
		return "NOTIMP"
	case Refused:
		return "REFUSED"
	default:
		return "Unknown RCode"
	}
}

func (opcode OpCode) String() string {
	switch opcode {
	case OpCodeQuery:
		return "Standard Query"
	case OpCodeIQuery:
		return "Inverse Query"
	case OpCodeStatus:
		return "Status"
	default:
		return "Unknown OpCode"
	}
}

func (messageType MessageType) String() string {
	switch messageType {
	case Query:
		return "Query"
	case Response:
		return "Response"
	default:
		return "Unknown QR"
	}
}

func (typ Type) String() string {
	switch typ {
	case TypeA:
		return "A"
	case TypeNS:
		return "NS"
	case TypeCNAME:
		return "CNAME"
	case TypeSOA:
		return "SOA"
	case TypeWKS:
		return "WKS"
	case TypePTR:
		return "PTR"
	case TypeMX:
		return "MX"
	case TypeTXT:
		return "TXT"
	case TypeAAAA:
		return "AAAA"
	case TypeSRV:
		return "SRV"
	case TypeHINFO:
		return "HINFO"
	case TypeMINFO:
		return "MINFO"
	case TypeAXFR:
		return "AXFR"
	case TypeALL:
		return "ALL"
	default:
		return "Unknown Type"
	}
}

func (class Class) String() string {
	switch class {
	case ClassINET:
		return "IN"
	case ClassCSNET:
		return "CS"
	case ClassCHAOS:
		return "CH"
	case ClassHESIOD:
		return "HS"
	case ClassANY:
		return "Any"
	default:
		return "Unknown Class"
	}
}

type Name struct {
	Domain string
	Compressed bool
}

func (name *Name) Decode(data []byte, off int) (Name, int) {
	name.Domain = ""
Loop:
	for {
		c := int(data[off])

		switch c & 0xc0 {
		case 0x00:
			if c == 0x00 {
				off++
				break Loop
			}
			off++
			name.Domain += string(data[off:off+c]) + "."
			off += c
		case 0xc0:
			off = DecodePtr(data, off)
		}
	}
	return *name, off
}


func DecodePtr(data []byte, off int) int {
	return int(binary.BigEndian.Uint16(data[off:off+2]) ^ 0xc000)
}

func (name *Name) Encode() []byte {
	if name.Compressed {
		return []byte{0xc0, 0x0c}
	}

	data := make([]byte, 0)
	parts := strings.Split(name.Domain, ".")
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}
		data = append(data, byte(len(part)))
		data = append(data, []byte(part)...)
	}
	data = append(data, 0)
	return data
}