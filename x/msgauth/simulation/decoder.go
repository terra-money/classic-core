package simulation

import (
	"bytes"
	"fmt"

	tmkv "github.com/tendermint/tendermint/libs/kv"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/terra-money/core/x/msgauth/internal/types"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding distribution type
func DecodeStore(cdc *codec.Codec, kvA, kvB tmkv.Pair) string {
	switch {
	case bytes.Equal(kvA.Key[:1], types.GrantKey):
		var grantA, grantB types.AuthorizationGrant
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &grantA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &grantB)
		return fmt.Sprintf("%v\n%v", grantA, grantB)
	case bytes.Equal(kvA.Key[:1], types.GrantQueueKey):
		var pairsA, pairsB []types.GGMPair
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &pairsA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &pairsB)
		return fmt.Sprintf("%v\n%v", pairsA, pairsB)
	default:
		panic(fmt.Sprintf("invalid market key prefix %X", kvA.Key[:1]))
	}
}
