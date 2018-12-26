package oracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// nolint
var (
	PrefixVote          = []byte{0x01}
	PrefixTargetPrice   = []byte{0x02}
	PrefixObservedPrice = []byte{0x03}
	KeyWhitelist        = []byte{0x04}

	ParamStoreKeyParams = []byte("params")
)

// ParamTable for oracle module
func ParamTypeTable() params.TypeTable {
	return params.NewTypeTable(
		ParamStoreKeyParams, Params{},
	)
}

const (
	// default paramspace for params keeper
	DefaultParamspace = "oracle"
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

// GetTargetPriceKey is in format of PrefixTargetPrice||denom
func GetTargetPriceKey(denom string) []byte {
	return append(PrefixTargetPrice, []byte(denom)...)
}

// GetObservedPriceKey is in format of PrefixObservedPrice||denom
func GetObservedPriceKey(denom string) []byte {
	return append(PrefixObservedPrice, []byte(denom)...)
}
