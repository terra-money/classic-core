package oracle

import (
	"sort"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper of the oracle store
type Keeper struct {
	key sdk.StoreKey
	cdc *codec.Codec

	valset        sdk.ValidatorSet
	timeout       int64
	supermajority sdk.Dec
}

// NewKeeper constructs a new keeper
func NewKeeper(key sdk.StoreKey, cdc *codec.Codec, valset sdk.ValidatorSet, supermajority sdk.Dec, timeout int64) Keeper {
	if timeout < 10000 {
		panic("Timeout should be a reasonably high number")
	}

	if supermajority.LT(sdk.NewDecWithPrec(5, 1)) {
		panic("Supermajority needs to be at least 50%")
	}

	return Keeper{
		key: key,
		cdc: cdc,

		valset:        valset,
		timeout:       timeout,
		supermajority: supermajority,
	}
}

func (keeper Keeper) AddVote(ctx sdk.Context, vote PriceVote) {
	store := ctx.KVStore(keeper.key)

	key := GetVoteKey(vote.Denom, vote.Feeder)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(vote)
	store.Set(key, bz)
}

func (keeper Keeper) Elect(ctx sdk.Context, denom string) (PriceVote, sdk.Error) {
	store := ctx.KVStore(keeper.key)

	votes := PriceVotes{}
	totalPowerVoted := sdk.ZeroDec()

	iter := sdk.KVStorePrefixIterator(store, GetVotePrefix(denom))
	for ; iter.Valid(); iter.Next() {
		var pv PriceVote
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &pv)

		votes = append(votes, pv)
		store.Delete(iter.Key())
	}
	iter.Close()

	// Not enough people have voted, skip
	if totalPowerVoted.LT(keeper.valset.TotalPower(ctx).Mul(keeper.supermajority)) {
		return PriceVote{}, ErrNotEnoughVotes(DefaultCodespace)
	}

	// Sort votes by price
	sort.Sort(votes)

	medPower := sdk.ZeroDec()
	median := PriceVote{}
	for i := 0; i < len(votes); i++ {
		medPower.Add(votes[i].Power)

		// Get the weighted median of the votes
		if medPower.GTE(totalPowerVoted.Mul(sdk.NewDecWithPrec(5, 1))) {
			median = votes[i]
		}
	}

	// Save new price elect
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(median)
	store.Set(GetElectKey(denom), bz)

	return median, nil
}

//nolint
func (keeper Keeper) GetElect(ctx sdk.Context, denom string) (res PriceVote) {
	store := ctx.KVStore(keeper.key)

	key := GetElectKey(denom)
	bz := store.Get(key)
	if bz == nil {

		return
	}
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)

	return
}
