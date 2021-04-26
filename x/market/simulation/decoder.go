package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/terra-project/core/x/market/types"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding market type.
func NewDecodeStore(cdc codec.Marshaler) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.MintPoolDeltaKey):
			var deltaA, deltaB sdk.DecProto
			cdc.MustUnmarshalBinaryBare(kvA.Value, &deltaA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &deltaB)
			return fmt.Sprintf("%v\n%v", deltaA, deltaB)
		case bytes.Equal(kvA.Key[:1], types.BurnPoolDeltaKey):
			var deltaA, deltaB sdk.DecProto
			cdc.MustUnmarshalBinaryBare(kvA.Value, &deltaA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &deltaB)
			return fmt.Sprintf("%v\n%v", deltaA, deltaB)
		default:
			panic(fmt.Sprintf("invalid market key prefix %X", kvA.Key[:1]))
		}
	}
}
