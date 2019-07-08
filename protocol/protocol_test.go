package protocol

import "testing"
import "reflect"

func TestEncodeQueryHeader(t *testing.T) {
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

func TestEncodeResponseHeader(t *testing.T) {
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

func TestDecodeHeader(t *testing.T) {
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
	header.Decode(data, 0)

	if !reflect.DeepEqual(header, expected) {
		t.Fail()
	}
}

func TestDecodePtr(t *testing.T) {
	data := []byte{
		0xc0, 0x0c,
	}
	if DecodePtr(data, 0) != 0x0c {
		t.Log(DecodePtr(data, 0))
		t.Fail()
	}
}

func TestDecodeName(t *testing.T) {
	data := []byte{
		0x04, 0x6f, 0x63, 0x73,
		0x70, 0x06, 0x64, 0x63,
		0x6f, 0x63, 0x73, 0x70,
		0x02, 0x63, 0x6e, 0x01,
		0x77, 0x08, 0x6b, 0x75,
		0x6e, 0x6c, 0x75, 0x6e,
		0x61, 0x72, 0x03, 0x63,
		0x6f, 0x6d, 0x00,
	}
	expected := "ocsp.dcocsp.cn.w.kunlunar.com."
	name := Name{}
	name.Decode(data, 0)
	if name.Domain != expected {
		t.Fail()
	}
}

func TestCompressedName(t *testing.T) {
	data := []byte{
		0xe3, 0xf4, 0x81, 0x80, 0x00, 0x01, 0x00, 0x02,
		0x00, 0x00, 0x00, 0x00, 0x04, 0x6f, 0x63, 0x73,
		0x70, 0x08, 0x64, 0x69, 0x67, 0x69, 0x63, 0x65,
		0x72, 0x74, 0x03, 0x63, 0x6f, 0x6d, 0x00, 0x00,
		0x01, 0x00, 0x01, 0xc0, 0x0c, 0x00, 0x05, 0x00,
		0x01, 0x00, 0x01, 0x19, 0x27, 0x00, 0x14, 0x03,
		0x63, 0x73, 0x39, 0x03, 0x77, 0x61, 0x63, 0x06,
		0x70, 0x68, 0x69, 0x63, 0x64, 0x6e, 0x03, 0x6e,
		0x65, 0x74, 0x00, 0xc0, 0x2f, 0x00, 0x01, 0x00,
		0x01, 0x00, 0x00, 0x0a, 0x66, 0x00, 0x04, 0x75,
		0x12, 0xed, 0x1d,
	}
	name := Name{}

	_, off := name.Decode(data, 35)
	expected := "ocsp.digicert.com."

	if name.Domain != expected || off != 37 {
		t.Fail()
	}

	_, off = name.Decode(data, 67)
	expected = "cs9.wac.phicdn.net."
	if name.Domain != expected || off != 69 {
		t.Fail()
	}
}

func TestEncodeName(t *testing.T) {
	name := Name{
		Domain: "deepzz.com.",
	}

	encoded := name.Encode()

	expected := []byte("\x06\x64\x65\x65\x70\x7a\x7a\x03\x63\x6f\x6d\x00")

	if !reflect.DeepEqual(encoded, expected) {
		t.Log(encoded)
		t.Log(expected)
		t.Fail()
	}
}

func TestDecodeQuestion(t *testing.T) {
	data := []byte("\x06\x64\x65\x65\x70\x7a\x7a\x03\x63\x6f\x6d\x00\x00\x01\x00\x01")

	question := Question{}

	expected := Question{
		Name: Name{
			Domain: "deepzz.com.",
		},
		Type:  TypeA,
		Class: ClassINET,
	}

	_, off := question.Decode(data, 0)
	if !reflect.DeepEqual(question, expected) {
		t.Log(question)
		t.Fail()
	}
	if off != 16 {
		t.Log(off)
		t.Fail()
	}
}

func TestEncodeQuestion(t *testing.T) {
	expected := []byte("\x06\x64\x65\x65\x70\x7a\x7a\x03\x63\x6f\x6d\x00\x00\x01\x00\x01")

	question := Question{
		Name: Name{
			Domain: "deepzz.com",
		},
		Type:  TypeA,
		Class: ClassINET,
	}

	encoded := question.Encode()
	if !reflect.DeepEqual(encoded, expected) {
		t.Log(encoded)
		t.Log(expected)
		t.Fail()
	}
}
