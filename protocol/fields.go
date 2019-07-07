package protocol

import "encoding/binary"

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

type Name struct {
	name string
}

func DecodeName(data []byte, off int) string {
	var name = ""
Loop:
	for {
		c := int(data[off])

		switch c & 0xc0 {
		case 0x00:
			if c == 0x00 {
				break Loop
			}
			off++
			name += string(data[off:off+c]) + "."
			off += c
		case 0xc0:
			off = DecodePtr(data[off : off+2])
		}
	}

	return name
}

func DecodePtr(data []byte) int {
	return int(binary.BigEndian.Uint16(data[0:2]) ^ 0xc000)
}
