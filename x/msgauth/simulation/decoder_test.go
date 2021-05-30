package simulation

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	tmkv "github.com/tendermint/tendermint/libs/kv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-money/core/x/msgauth/internal/types"
)

func makeTestCodec() (cdc *codec.Codec) {
	cdc = codec.New()
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	types.RegisterCodec(cdc)
	return
}

func TestDecodeMsgAuthStore(t *testing.T) {
	cdc := makeTestCodec()

	grant := types.NewAuthorizationGrant(types.NewSendAuthorization(sdk.NewCoins(sdk.NewInt64Coin("foo", 123))), time.Now().UTC())
	pairs := []types.GGMPair{
		{
			GranteeAddress: sdk.AccAddress{1, 2, 3},
			GranterAddress: sdk.AccAddress{1, 2, 3},
			MsgType:        "send",
		},
	}

	kvPairs := tmkv.Pairs{
		tmkv.Pair{Key: types.GrantKey, Value: cdc.MustMarshalBinaryLengthPrefixed(grant)},
		tmkv.Pair{Key: types.GrantQueueKey, Value: cdc.MustMarshalBinaryLengthPrefixed(pairs)},
		tmkv.Pair{Key: []byte{0x99}, Value: []byte{0x99}},
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
				require.Panics(t, func() { DecodeStore(cdc, kvPairs[i], kvPairs[i]) }, tt.name)
			default:
				require.Equal(t, tt.expectedLog, DecodeStore(cdc, kvPairs[i], kvPairs[i]), tt.name)
			}
		})
	}
}
