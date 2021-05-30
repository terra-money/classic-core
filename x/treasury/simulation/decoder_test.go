package simulation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	tmkv "github.com/tendermint/tendermint/libs/kv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/treasury/internal/types"
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

	taxRate := sdk.NewDecWithPrec(123, 2)
	rewardWeight := sdk.NewDecWithPrec(532, 2)
	taxCap := sdk.NewInt(1600)
	taxProceeds := sdk.NewCoins(sdk.NewInt64Coin(core.MicroKRWDenom, 123124), sdk.NewInt64Coin(core.MicroSDRDenom, 123124))
	epochInitialIssuance := sdk.NewCoins(sdk.NewInt64Coin(core.MicroKRWDenom, 645352342))

	TR := sdk.NewDecWithPrec(123, 2)
	SR := sdk.NewDecWithPrec(43523, 4)
	TSL := sdk.NewInt(1245213)

	kvPairs := tmkv.Pairs{
		tmkv.Pair{Key: types.TaxRateKey, Value: cdc.MustMarshalBinaryLengthPrefixed(taxRate)},
		tmkv.Pair{Key: types.RewardWeightKey, Value: cdc.MustMarshalBinaryLengthPrefixed(rewardWeight)},
		tmkv.Pair{Key: types.TaxCapKey, Value: cdc.MustMarshalBinaryLengthPrefixed(taxCap)},
		tmkv.Pair{Key: types.TaxProceedsKey, Value: cdc.MustMarshalBinaryLengthPrefixed(taxProceeds)},
		tmkv.Pair{Key: types.EpochInitialIssuanceKey, Value: cdc.MustMarshalBinaryLengthPrefixed(epochInitialIssuance)},
		tmkv.Pair{Key: types.TRKey, Value: cdc.MustMarshalBinaryLengthPrefixed(TR)},
		tmkv.Pair{Key: types.SRKey, Value: cdc.MustMarshalBinaryLengthPrefixed(SR)},
		tmkv.Pair{Key: types.TSLKey, Value: cdc.MustMarshalBinaryLengthPrefixed(TSL)},
		tmkv.Pair{Key: []byte{0x99}, Value: []byte{0x99}},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{"TaxRate", fmt.Sprintf("%v\n%v", taxRate, taxRate)},
		{"RewardWeight", fmt.Sprintf("%v\n%v", rewardWeight, rewardWeight)},
		{"TaxCap", fmt.Sprintf("%v\n%v", taxCap, taxCap)},
		{"TaxProceeds", fmt.Sprintf("%v\n%v", taxProceeds, taxProceeds)},
		{"EpochInitialIssuance", fmt.Sprintf("%v\n%v", epochInitialIssuance, epochInitialIssuance)},
		{"TR", fmt.Sprintf("%v\n%v", TR, TR)},
		{"SR", fmt.Sprintf("%v\n%v", SR, SR)},
		{"TSL", fmt.Sprintf("%v\n%v", TSL, TSL)},
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
