package messaging

import (
	"bytes"
	"testing"
)

type SerialiseStringResult struct {
	s        string
	b        []byte
	expected bool
}

var serialiseStringResults = []SerialiseStringResult{
	{"foobar", []byte{34, 102, 111, 111, 98, 97, 114, 34}, true},
	{"", []byte{34, 34}, true},
}

func TestSerialiseString(t *testing.T) {
	for _, test := range serialiseStringResults {
		b, isOK := SerialiseString(test.s)
		result := bytes.Compare(b, test.b)

		if result != 0 {
			t.Fatalf("Expected result: %v \n Actual result: %v\n", test.b, b)
		}
		if isOK != test.expected {
			t.Fatalf("Expected result: %v \n Actual result: %v\n", test.expected, isOK)
		}
	}
}

type DeserialiseStringResult struct {
	raw      []byte
	s        string
	expected bool
}

var deserialiseStringResults = []DeserialiseStringResult{
	{[]byte{34, 102, 111, 111, 98, 97, 114, 34}, "foobar", true},
	{[]byte{34, 34}, "", true},
}

func TestDeserialiseString(t *testing.T) {
	for _, test := range deserialiseStringResults {
		s, isOK := DeserialiseString(test.raw)

		if s != test.s {
			t.Fatalf("Expected result: %v \n Actual result: %v\n", test.s, &s)
		}
		if isOK != test.expected {
			t.Fatalf("Expected result: %v \n Actual result: %v\n", test.expected, isOK)
		}
	}
}
