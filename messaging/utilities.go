package messaging

import (
	"encoding/json"
	"log"
)

// SerialiseString - uses json.marshal to serialise strings
func SerialiseString(s string) (b []byte, isOK bool) {
	msg, err := json.Marshal(s)

	if err != nil {
		log.Println(err)
		return nil, false
	}

	return msg, true
}

// DeserialiseString - uses json.Unmarshal to deserialise json data to a string
func DeserialiseString(raw json.RawMessage) (s string, isOK bool) {
	bytes, err := json.Marshal(raw)
	if err != nil {
		log.Println(err)
		return "", false
	}

	err = json.Unmarshal(bytes, &s)
	if err != nil {
		log.Println(err)
		return "", false
	}

	return s, true
}
