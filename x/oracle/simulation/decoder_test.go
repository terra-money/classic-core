package simulation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto/ed25519"
	tmkv "github.com/tendermint/tendermint/libs/kv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/oracle/internal/types"
)

var (
	delPk      = ed25519.GenPrivKey().PubKey()
	feederAddr = sdk.AccAddress(delPk.Address())
	valAddr    = sdk.ValAddress(delPk.Address())
)

func makeTestCodec() (cdc *codec.Codec) {
	cdc = codec.New()
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	types.RegisterCodec(cdc)
	return
}

func TestDecodeDistributionStore(t *testing.T) {
	cdc := makeTestCodec()

	prevote := types.NewExchangeRatePrevote("12345", core.MicroKRWDenom, valAddr, 123)
	vote := types.NewExchangeRateVote(sdk.NewDecWithPrec(1234, 1), core.MicroKRWDenom, valAddr)
	exchangeRate := sdk.NewDecWithPrec(1234, 1)
	missCounter := 123

	kvPairs := tmkv.Pairs{
		tmkv.Pair{Key: types.PrevoteKey, Value: cdc.MustMarshalBinaryLengthPrefixed(prevote)},
		tmkv.Pair{Key: types.VoteKey, Value: cdc.MustMarshalBinaryLengthPrefixed(vote)},
		tmkv.Pair{Key: types.ExchangeRateKey, Value: cdc.MustMarshalBinaryLengthPrefixed(exchangeRate)},
		tmkv.Pair{Key: types.FeederDelegationKey, Value: cdc.MustMarshalBinaryLengthPrefixed(feederAddr)},
		tmkv.Pair{Key: types.MissCounterKey, Value: cdc.MustMarshalBinaryLengthPrefixed(missCounter)},
		tmkv.Pair{Key: []byte{0x99}, Value: []byte{0x99}},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{"Prevote", fmt.Sprintf("%v\n%v", prevote, prevote)},
		{"Vote", fmt.Sprintf("%v\n%v", vote, vote)},
		{"ExchangeRate", fmt.Sprintf("%v\n%v", exchangeRate, exchangeRate)},
		{"FeederDelegation", fmt.Sprintf("%v\n%v", feederAddr, feederAddr)},
		{"MissCounter", fmt.Sprintf("%v\n%v", missCounter, missCounter)},
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
