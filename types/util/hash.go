package util

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v2"
)

var _ yaml.Marshaler = HexBytes{}

// HexBytes is hash value of wasm bytes code
type HexBytes []byte

// HexBytesFromHexString convert hex string to HexBytes
func HexBytesFromHexString(s string) (HexBytes, error) {
	h, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	return h, nil
}

// String implements fmt.Stringer interface
func (h HexBytes) String() string {
	return hex.EncodeToString(h)
}

// Equal does bytes equal check
func (h HexBytes) Equal(h2 HexBytes) bool {
	return bytes.Equal(h, h2)
}

// Empty check the name hash has zero length
func (h HexBytes) Empty() bool {
	return len(h) == 0
}

// Bytes returns the raw address bytes.
func (h HexBytes) Bytes() []byte {
	return h
}

// Format implements the fmt.Formatter interface.
func (h HexBytes) Format(s fmt.State, verb rune) {
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
func (h HexBytes) Marshal() ([]byte, error) {
	return h, nil
}

// Unmarshal sets the address to the given data. It is needed for protobuf
// compatibility.
func (h *HexBytes) Unmarshal(data []byte) error {
	*h = data
	return nil
}

// MarshalJSON marshals to JSON using Bech32.
func (h HexBytes) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

// MarshalYAML marshals to YAML using Bech32.
func (h HexBytes) MarshalYAML() (interface{}, error) {
	return h.String(), nil
}

// UnmarshalJSON unmarshals from JSON assuming Bech32 encoding.
func (h *HexBytes) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	h2, err := HexBytesFromHexString(s)
	if err != nil {
		return err
	}

	*h = h2
	return nil
}
