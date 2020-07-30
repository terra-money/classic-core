package util

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v2"
)

var _ yaml.Marshaler = Base64Bytes{}

// Base64Bytes is hash value of wasm bytes code
type Base64Bytes []byte

// Base64FromHexString convert base64 string to Base64
func Base64FromHexString(s string) (Base64Bytes, error) {
	h, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}

	return h, nil
}

// String implements fmt.Stringer interface
func (h Base64Bytes) String() string {
	return base64.StdEncoding.EncodeToString(h)
}

// Equal does bytes equal check
func (h Base64Bytes) Equal(h2 Base64Bytes) bool {
	return bytes.Equal(h, h2)
}

// Empty check the name hash has zero length
func (h Base64Bytes) Empty() bool {
	return len(h) == 0
}

// Bytes returns the raw address bytes.
func (h Base64Bytes) Bytes() []byte {
	return h
}

// Format implements the fmt.Formatter interface.
func (h Base64Bytes) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		_, _ = s.Write([]byte(h.String()))
	case 'p':
		_, _ = s.Write([]byte(fmt.Sprintf("%p", h)))
	default:
		_, _ = s.Write([]byte(fmt.Sprintf("%X", []byte(h))))
	}
}

// Marshal returns the raw address bytes. It is needed for protobuf
// compatibility.
func (h Base64Bytes) Marshal() ([]byte, error) {
	return h, nil
}

// Unmarshal sets the address to the given data. It is needed for protobuf
// compatibility.
func (h *Base64Bytes) Unmarshal(data []byte) error {
	*h = data
	return nil
}

// MarshalJSON marshals to JSON using base64.
func (h Base64Bytes) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

// MarshalYAML marshals to YAML using base64.
func (h Base64Bytes) MarshalYAML() (interface{}, error) {
	return h.String(), nil
}

// UnmarshalJSON unmarshals from JSON assuming base64 encoding.
func (h *Base64Bytes) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	h2, err := Base64FromHexString(s)
	if err != nil {
		return err
	}

	*h = h2
	return nil
}
