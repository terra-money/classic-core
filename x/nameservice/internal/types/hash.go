package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"gopkg.in/yaml.v2"
)

var _ yaml.Marshaler = NameHash{}
var _ yaml.Marshaler = BidHash{}

// NameHash holds hash byte of name
type NameHash []byte

// NameHashFromHexString returns NameHash of the given hex string
func NameHashFromHexString(hexStr string) (NameHash, error) {
	hash, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}

	return NameHash(hash), err
}

// GetNameHash returns first 20 bytes of sha256(name)
func GetNameHash(name string) NameHash {
	if len(name) == 0 {
		return NameHash{}
	}

	hash := tmhash.NewTruncated()
	_, err := hash.Write([]byte(name))
	if err != nil {
		panic(err)
	}

	bz := hash.Sum(nil)
	return bz
}

// String implements fmt.Stringer interface
func (h NameHash) String() string {
	return hex.EncodeToString(h)
}

// Equal does bytes equal check
func (h NameHash) Equal(h2 NameHash) bool {
	return bytes.Equal(h, h2)
}

// Empty check the name hash has zero length
func (h NameHash) Empty() bool {
	return len(h) == 0
}

// Format implements the fmt.Formatter interface.
func (h NameHash) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		_, _ = s.Write([]byte(h.String()))
	default:
		_, _ = s.Write(h)
	}
}

// Marshal returns the raw address bytes. It is needed for protobuf
// compatibility.
func (h NameHash) Marshal() ([]byte, error) {
	return h, nil
}

// Unmarshal sets the address to the given data. It is needed for protobuf
// compatibility.
func (h *NameHash) Unmarshal(data []byte) error {
	*h = data
	return nil
}

// MarshalJSON marshals to JSON using Bech32.
func (h NameHash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

// MarshalYAML marshals to YAML using Bech32.
func (h NameHash) MarshalYAML() (interface{}, error) {
	return h.String(), nil
}

// UnmarshalJSON unmarshals from JSON assuming Bech32 encoding.
func (h *NameHash) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	h2, err := NameHashFromHexString(s)
	if err != nil {
		return err
	}

	*h = h2
	return nil
}

// BidHash holds hash of "salt:name:amount:bidder"
type BidHash []byte

// BidHashFromHexString returns BidHash of the given hex string
func BidHashFromHexString(hexStr string) (BidHash, error) {
	hash, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}

	return BidHash(hash), err
}

// GetBidHash returns first 20 bytes of sha256("salt:name:amount:bidder")
func GetBidHash(salt string, name Name, amount sdk.Coin, bidder sdk.AccAddress) BidHash {
	hash := tmhash.NewTruncated()
	_, err := hash.Write([]byte(fmt.Sprintf("%s:%s:%s:%s", salt, name, amount, bidder)))
	if err != nil {
		panic(err)
	}

	bz := hash.Sum(nil)
	return bz
}

// String implements fmt.Stringer interface
func (h BidHash) String() string {
	return hex.EncodeToString(h)
}

// Equal does bytes equal check
func (h BidHash) Equal(h2 BidHash) bool {
	return bytes.Equal(h, h2)
}

// Empty check the bid hash has zero length
func (h BidHash) Empty() bool {
	return len(h) == 0
}

// Format implements the fmt.Formatter interface.
func (h BidHash) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		_, _ = s.Write([]byte(h.String()))
	default:
		_, _ = s.Write(h)
	}
}

// Marshal returns the raw address bytes. It is needed for protobuf
// compatibility.
func (h BidHash) Marshal() ([]byte, error) {
	return h, nil
}

// Unmarshal sets the address to the given data. It is needed for protobuf
// compatibility.
func (h *BidHash) Unmarshal(data []byte) error {
	*h = data
	return nil
}

// MarshalJSON marshals to JSON using Bech32.
func (h BidHash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

// MarshalYAML marshals to YAML using Bech32.
func (h BidHash) MarshalYAML() (interface{}, error) {
	return h.String(), nil
}

// UnmarshalJSON unmarshals from JSON assuming Bech32 encoding.
func (h *BidHash) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	h2, err := BidHashFromHexString(s)
	if err != nil {
		return err
	}

	*h = h2
	return nil
}
