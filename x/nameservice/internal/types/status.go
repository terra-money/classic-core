package types

import (
	"encoding/json"
	"fmt"
)

type (
	// AuctionStatus is a type alias that represents a registry status as a byte
	AuctionStatus byte
)

//nolint
const (
	AuctionStatusNil    AuctionStatus = 0x00
	AuctionStatusBid    AuctionStatus = 0x01
	AuctionStatusReveal AuctionStatus = 0x02
)

// AuctionStatusToString turns a string into a AuctionStatus
func AuctionStatusFromString(str string) (AuctionStatus, error) {
	switch str {
	case "Nil":
		return AuctionStatusNil, nil

	case "Bid":
		return AuctionStatusBid, nil

	case "Reveal":
		return AuctionStatusReveal, nil

	default:
		return AuctionStatus(0xff), fmt.Errorf("'%s' is not a valid auction status", str)
	}
}

// Marshal needed for protobuf compatibility
func (status AuctionStatus) Marshal() ([]byte, error) {
	return []byte{byte(status)}, nil
}

// Unmarshal needed for protobuf compatibility
func (status *AuctionStatus) Unmarshal(data []byte) error {
	*status = AuctionStatus(data[0])
	return nil
}

// Marshals to JSON using string
func (status AuctionStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(status.String())
}

// Unmarshals from JSON assuming Bech32 encoding
func (status *AuctionStatus) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	bz2, err := AuctionStatusFromString(s)
	if err != nil {
		return err
	}

	*status = bz2
	return nil
}

// String implements the Stringer interface.
func (status AuctionStatus) String() string {
	switch status {
	case AuctionStatusBid:
		return "Bid"

	case AuctionStatusReveal:
		return "Reveal"

	case AuctionStatusNil:
		return "Nil"

	default:
		return ""
	}
}

// Format implements the fmt.Formatter interface.
func (status AuctionStatus) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		_, _ = s.Write([]byte(status.String()))
	default:
		_, _ = s.Write([]byte{byte(status)})
	}
}
