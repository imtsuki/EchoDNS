package protocol

import "testing"
import "reflect"

func TestQueryHeaderEncode(t *testing.T) {
	header := Header{
		ID:               0x0b65,
		MessageType:      Query,
		OpCode:           OpCodeQuery,
		Truncation:       false,
		RecursionDesired: true,
		QuestionCount:    1,
	}
	expected := []byte{
		0x0b, 0x65, 0x01, 0x00,
		0x00, 0x01, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	if !reflect.DeepEqual(header.Encode(), expected) {
		t.Fail()
	}
}

func TestResponseHeaderEncode(t *testing.T) {
	header := Header{
		ID:                 0x4a5a,
		MessageType:        Response,
		OpCode:             OpCodeQuery,
		Authoritative:      false,
		Truncation:         false,
		RecursionDesired:   true,
		RecursionAvailable: true,
		QuestionCount:      1,
		AnswerCount:        1,
	}
	expected := []byte{
		0x4a, 0x5a, 0x81, 0x80,
		0x00, 0x01, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00,
	}

	if !reflect.DeepEqual(header.Encode(), expected) {
		t.Fail()
	}
}

func TestHeaderDecode(t *testing.T) {
	expected := Header{
		ID:                 0x4a5a,
		MessageType:        Response,
		OpCode:             OpCodeQuery,
		Authoritative:      false,
		Truncation:         false,
		RecursionDesired:   true,
		RecursionAvailable: true,
		QuestionCount:      1,
		AnswerCount:        1,
	}
	data := []byte{
		0x4a, 0x5a, 0x81, 0x80,
		0x00, 0x01, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00,
	}
	var header Header
	if !reflect.DeepEqual(header.Decode(data), expected) {
		t.Fail()
	}
	if !reflect.DeepEqual(header, expected) {
		t.Fail()
	}
}
