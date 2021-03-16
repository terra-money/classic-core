package simulation

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/types/kv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-project/core/x/msgauth/keeper"
	"github.com/terra-project/core/x/msgauth/types"
)

func TestDecodeMsgAuthStore(t *testing.T) {
	cdc := keeper.MakeTestCodec(t)
	dec := NewDecodeStore(cdc)

	grant, err := types.NewAuthorizationGrant(types.NewSendAuthorization(sdk.NewCoins(sdk.NewInt64Coin("foo", 123))), time.Now().UTC())
	require.NoError(t, err)

	pairs := types.GGMPairs{
		Pairs: []types.GGMPair{
			{
				GranteeAddress: "abc",
				GranterAddress: "cba",
				MsgType:        "send",
			},
		},
	}

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{Key: types.GrantKey, Value: cdc.MustMarshalBinaryBare(&grant)},
			{Key: types.GrantQueueKey, Value: cdc.MustMarshalBinaryBare(&pairs)},
			{Key: []byte{0x99}, Value: []byte{0x99}},
		},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{"Grant", fmt.Sprintf("%v\n%v", grant, grant)},
		{"GGMPair", fmt.Sprintf("%v\n%v", pairs, pairs)},
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
