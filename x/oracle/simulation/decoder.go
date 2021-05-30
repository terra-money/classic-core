package simulation

import (
	"bytes"
	"fmt"

	tmkv "github.com/tendermint/tendermint/libs/kv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-money/core/x/oracle/internal/types"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding distribution type
func DecodeStore(cdc *codec.Codec, kvA, kvB tmkv.Pair) string {
	switch {
	case bytes.Equal(kvA.Key[:1], types.PrevoteKey):
		var prevoteA, prevoteB types.ExchangeRatePrevote
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &prevoteA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &prevoteB)
		return fmt.Sprintf("%v\n%v", prevoteA, prevoteB)
	case bytes.Equal(kvA.Key[:1], types.VoteKey):
		var voteA, voteB types.ExchangeRateVote
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &voteA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &voteB)
		return fmt.Sprintf("%v\n%v", voteA, voteB)
	case bytes.Equal(kvA.Key[:1], types.ExchangeRateKey):
		var exchangeRateA, exchangeRateB sdk.Dec
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &exchangeRateA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &exchangeRateB)
		return fmt.Sprintf("%v\n%v", exchangeRateA, exchangeRateB)
	case bytes.Equal(kvA.Key[:1], types.FeederDelegationKey):
		var addressA, addressB sdk.AccAddress
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &addressA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &addressB)
		return fmt.Sprintf("%v\n%v", addressA, addressB)
	case bytes.Equal(kvA.Key[:1], types.MissCounterKey):
		var counterA, counterB int64
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &counterA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &counterB)
		return fmt.Sprintf("%v\n%v", counterA, counterB)
	case bytes.Equal(kvA.Key[:1], types.AggregateExchangeRatePrevoteKey):
		var prevoteA, prevoteB types.AggregateExchangeRatePrevote
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &prevoteA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &prevoteB)
		return fmt.Sprintf("%v\n%v", prevoteA, prevoteB)
	case bytes.Equal(kvA.Key[:1], types.AggregateExchangeRateVoteKey):
		var voteA, voteB types.AggregateExchangeRateVote
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &voteA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &voteB)
		return fmt.Sprintf("%v\n%v", voteA, voteB)
	case bytes.Equal(kvA.Key[:1], types.TobinTaxKey):
		var tobinTaxA, tobinTaxB sdk.Dec
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &tobinTaxA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &tobinTaxB)
		return fmt.Sprintf("%v\n%v", tobinTaxA, tobinTaxB)
	default:
		panic(fmt.Sprintf("invalid oracle key prefix %X", kvA.Key[:1]))
	}
}
