package simulation

import (
	"bytes"
	"fmt"

	tmkv "github.com/tendermint/tendermint/libs/kv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/core/x/treasury/internal/types"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding distribution type
func DecodeStore(cdc *codec.Codec, kvA, kvB tmkv.Pair) string {
	switch {
	case bytes.Equal(kvA.Key[:1], types.TaxRateKey):
		var taxRateA, taxRateB sdk.Dec
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &taxRateA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &taxRateB)
		return fmt.Sprintf("%v\n%v", taxRateA, taxRateB)
	case bytes.Equal(kvA.Key[:1], types.RewardWeightKey):
		var rewardWeightA, rewardWeightB sdk.Dec
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &rewardWeightA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &rewardWeightB)
		return fmt.Sprintf("%v\n%v", rewardWeightA, rewardWeightB)
	case bytes.Equal(kvA.Key[:1], types.TaxCapKey):
		var taxCapA, taxCapB sdk.Int
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &taxCapA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &taxCapB)
		return fmt.Sprintf("%v\n%v", taxCapA, taxCapB)
	case bytes.Equal(kvA.Key[:1], types.TaxProceedsKey):
		var taxProceedsA, taxProceedsB sdk.Coins
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &taxProceedsA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &taxProceedsB)
		return fmt.Sprintf("%v\n%v", taxProceedsA, taxProceedsB)
	case bytes.Equal(kvA.Key[:1], types.EpochInitialIssuanceKey):
		var epochInitialIssuanceA, epochInitialIssuanceB sdk.Coins
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &epochInitialIssuanceA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &epochInitialIssuanceB)
		return fmt.Sprintf("%v\n%v", epochInitialIssuanceA, epochInitialIssuanceB)
	case bytes.Equal(kvA.Key[:1], types.TRKey):
		var TaxRateA, TaxRateB sdk.Dec
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &TaxRateA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &TaxRateB)
		return fmt.Sprintf("%v\n%v", TaxRateA, TaxRateB)
	case bytes.Equal(kvA.Key[:1], types.SRKey):
		var SeigniorageRateA, SeigniorageRateB sdk.Dec
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &SeigniorageRateA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &SeigniorageRateB)
		return fmt.Sprintf("%v\n%v", SeigniorageRateA, SeigniorageRateB)
	case bytes.Equal(kvA.Key[:1], types.TSLKey):
		var TotalStakedLunaA, TotalStakedLunaB sdk.Int
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &TotalStakedLunaA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &TotalStakedLunaB)
		return fmt.Sprintf("%v\n%v", TotalStakedLunaA, TotalStakedLunaB)
	default:
		panic(fmt.Sprintf("invalid oracle key prefix %X", kvA.Key[:1]))
	}
}
