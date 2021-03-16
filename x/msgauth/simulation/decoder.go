package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/terra-project/core/x/msgauth/types"
)

// NewDecodeStore unmarshals the KVPair's Value to the corresponding distribution type
func NewDecodeStore(cdc codec.Marshaler) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.GrantKey):
			var grantA, grantB types.AuthorizationGrant
			cdc.MustUnmarshalBinaryBare(kvA.Value, &grantA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &grantB)
			return fmt.Sprintf("%v\n%v", grantA, grantB)
		case bytes.Equal(kvA.Key[:1], types.GrantQueueKey):
			var pairsA, pairsB types.GGMPairs
			cdc.MustUnmarshalBinaryBare(kvA.Value, &pairsA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &pairsB)
			return fmt.Sprintf("%v\n%v", pairsA, pairsB)
		default:
			panic(fmt.Sprintf("invalid msgauth key prefix %X", kvA.Key[:1]))
		}
	}
}
