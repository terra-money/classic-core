package simulation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/terra-project/core/x/market/keeper"
	"github.com/terra-project/core/x/market/types"
)

func TestDecodeDistributionStore(t *testing.T) {
	cdc := keeper.MakeTestCodec(t)
	dec := NewDecodeStore(cdc)

	mintDelta := sdk.NewDecWithPrec(12, 2)
	burnDelta := sdk.NewDecWithPrec(121, 2)

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{Key: types.MintPoolDeltaKey, Value: cdc.MustMarshalBinaryBare(&sdk.DecProto{Dec: mintDelta})},
			{Key: types.BurnPoolDeltaKey, Value: cdc.MustMarshalBinaryBare(&sdk.DecProto{Dec: burnDelta})},
			{Key: []byte{0x99}, Value: []byte{0x99}},
		},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{"MintPoolDelta", fmt.Sprintf("%v\n%v", mintDelta, mintDelta)},
		{"BurnPoolDelta", fmt.Sprintf("%v\n%v", burnDelta, burnDelta)},
		{"other", ""},
	}

	for i, tt := range tests {
		i, tt := i, tt
		t.Run(tt.name, func(t *testing.T) {
			switch i {
			case len(tests) - 1:
				require.Panics(t, func() { dec(kvPairs.Pairs[i], kvPairs.Pairs[i]) }, tt.name)
			default:
				require.Equal(t, tt.expectedLog, dec(kvPairs.Pairs[i], kvPairs.Pairs[i]), tt.name)
			}
		})
	}
}
