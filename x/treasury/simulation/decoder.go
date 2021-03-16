package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/terra-project/core/x/treasury/types"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding treasury type.
func NewDecodeStore(cdc codec.Marshaler) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.TaxRateKey):
			var taxRateA, taxRateB sdk.DecProto
			cdc.MustUnmarshalBinaryBare(kvA.Value, &taxRateA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &taxRateB)
			return fmt.Sprintf("%v\n%v", taxRateA, taxRateB)
		case bytes.Equal(kvA.Key[:1], types.RewardWeightKey):
			var rewardWeightA, rewardWeightB sdk.DecProto
			cdc.MustUnmarshalBinaryBare(kvA.Value, &rewardWeightA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &rewardWeightB)
			return fmt.Sprintf("%v\n%v", rewardWeightA, rewardWeightB)
		case bytes.Equal(kvA.Key[:1], types.TaxCapKey):
			var taxCapA, taxCapB sdk.IntProto
			cdc.MustUnmarshalBinaryBare(kvA.Value, &taxCapA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &taxCapB)
			return fmt.Sprintf("%v\n%v", taxCapA, taxCapB)
		case bytes.Equal(kvA.Key[:1], types.TaxProceedsKey):
			var taxProceedsA, taxProceedsB types.EpochTaxProceeds
			cdc.MustUnmarshalBinaryBare(kvA.Value, &taxProceedsA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &taxProceedsB)
			return fmt.Sprintf("%v\n%v", taxProceedsA.TaxProceeds, taxProceedsB.TaxProceeds)
		case bytes.Equal(kvA.Key[:1], types.EpochInitialIssuanceKey):
			var epochInitialIssuanceA, epochInitialIssuanceB types.EpochInitialIssuance
			cdc.MustUnmarshalBinaryBare(kvA.Value, &epochInitialIssuanceA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &epochInitialIssuanceB)
			return fmt.Sprintf("%v\n%v", epochInitialIssuanceA.Issuance, epochInitialIssuanceB.Issuance)
		case bytes.Equal(kvA.Key[:1], types.TRKey):
			var TaxRateA, TaxRateB sdk.DecProto
			cdc.MustUnmarshalBinaryBare(kvA.Value, &TaxRateA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &TaxRateB)
			return fmt.Sprintf("%v\n%v", TaxRateA, TaxRateB)
		case bytes.Equal(kvA.Key[:1], types.SRKey):
			var SeigniorageRateA, SeigniorageRateB sdk.DecProto
			cdc.MustUnmarshalBinaryBare(kvA.Value, &SeigniorageRateA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &SeigniorageRateB)
			return fmt.Sprintf("%v\n%v", SeigniorageRateA, SeigniorageRateB)
		case bytes.Equal(kvA.Key[:1], types.TSLKey):
			var TotalStakedLunaA, TotalStakedLunaB sdk.IntProto
			cdc.MustUnmarshalBinaryBare(kvA.Value, &TotalStakedLunaA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &TotalStakedLunaB)
			return fmt.Sprintf("%v\n%v", TotalStakedLunaA, TotalStakedLunaB)
		default:
			panic(fmt.Sprintf("invalid oracle key prefix %X", kvA.Key[:1]))
		}
	}
}
