package messaging

import (
	"bytes"
	"testing"
)

type SerialiseStringResult struct {
	s    string
	b    []byte
	isOK bool
}

var serialiseStringResults = []SerialiseStringResult{
	{"foobar", []byte{34, 102, 111, 111, 98, 97, 114, 34}, true},
	//{"", nil, false},
}

func TestSerialiseString(t *testing.T) {
	for _, test := range serialiseStringResults {
		b, isOK := SerialiseString(test.s)
		result := bytes.Compare(b, test.b)

		if result != 0 {
			t.Fatalf("Expected result: %v \n Actual result: %v\n", test.b, b)
		}
		if isOK != test.isOK {
			t.Fatalf("Expected result: %v \n Actual result: %v\n", test.isOK, isOK)
		}
	}
}

type DeserialiseStringResult struct {
	raw  []byte
	s    string
	isOK bool
}

var deserialiseStringResults = []DeserialiseStringResult{
	{[]byte{34, 102, 111, 111, 98, 97, 114, 34}, "foobar", true},
	//{nil, "", false},
}

func TestDeserialiseString(t *testing.T) {
	for _, test := range deserialiseStringResults {
		s, isOK := DeserialiseString(test.raw)

		if s != test.s {
			t.Fatalf("Expected result: %v \n Actual result: %v\n", test.s, &s)
		}
		if isOK != test.isOK {
			t.Fatalf("Expected result: %v \n Actual result: %v\n", test.isOK, isOK)
		}
	}
}
