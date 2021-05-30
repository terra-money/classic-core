package simulation

import (
	"bytes"
	"fmt"

	tmkv "github.com/tendermint/tendermint/libs/kv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-money/core/x/market/internal/types"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding distribution type
func DecodeStore(cdc *codec.Codec, kvA, kvB tmkv.Pair) string {
	switch {
	case bytes.Equal(kvA.Key[:1], types.TerraPoolDeltaKey):
		var deltaA, deltaB sdk.Dec
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &deltaA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &deltaB)
		return fmt.Sprintf("%v\n%v", deltaA, deltaB)
	default:
		panic(fmt.Sprintf("invalid market key prefix %X", kvA.Key[:1]))
	}
}
