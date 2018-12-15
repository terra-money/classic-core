package oracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// nolint
var (
	PrefixVote    = []byte{0x00}
	PrefixElect   = []byte{0x01}
	KeyThreshold  = []byte{0x02}
	KeyVotePeriod = []byte{0x03}
	KeyWhitelist  = []byte{0x04}
)

// GetVotePrefix is in format of prefix||denom
func GetVotePrefix(denom string) []byte {
	return append(PrefixVote, []byte(denom)...)
}

// GetVoteKey Key is in format of PrefixVote||denom||voter.AccAddress
func GetVoteKey(denom string, voter sdk.AccAddress) []byte {

	key := make([]byte, 1+len(denom)+sdk.AddrLen)

	copy(key[0:1], PrefixVote)
	copy(key[1:len(denom)+1], []byte(denom))
	copy(key[len(denom)+1:], voter.Bytes())

	return key
}

// GetElectKey is in format of PrefixElect||denom
func GetElectKey(denom string) []byte {
	return append(PrefixElect, []byte(denom)...)
}
