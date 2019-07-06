package protocol

// OpCode OpCode
type OpCode uint8

// enum OpCode
const (
	OpCodeQuery  OpCode = 0 // 0: QUERY
	OpCodeIQuery        = 1 // 1: IQUERY
	OpCodeStatus        = 2 // 2: STATUS
)

// MessageType QR
type MessageType uint8

// enum MessageType
const (
	Query    MessageType = 0
	Response             = 1
)

// ResponseCode RCode
type ResponseCode uint8

// enum ResponseCode
const (
	NoError  ResponseCode = 0
	FormErr               = 1
	ServFail              = 2
	NXDomain              = 3
	NotImp                = 4
	Refused               = 5
)
