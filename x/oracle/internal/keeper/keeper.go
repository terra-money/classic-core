package keeper

import (
	"fmt"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/oracle/internal/types"
)

// Keeper of the oracle store
type Keeper struct {
	cdc        *codec.Codec
	storeKey   sdk.StoreKey
	paramSpace params.Subspace

	distrKeeper   types.DistributionKeeper
	StakingKeeper types.StakingKeeper
	supplyKeeper  types.SupplyKeeper

	distrName string

	// codespace
	codespace sdk.CodespaceType
}

// NewKeeper constructs a new keeper for oracle
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey,
	paramspace params.Subspace, distrKeeper types.DistributionKeeper,
	stakingKeeper types.StakingKeeper, supplyKeeper types.SupplyKeeper,
	distrName string, codespace sdk.CodespaceType) Keeper {
	return Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		paramSpace:    paramspace.WithKeyTable(ParamKeyTable()),
		distrKeeper:   distrKeeper,
		StakingKeeper: stakingKeeper,
		supplyKeeper:  supplyKeeper,
		distrName:     distrName,
		codespace:     codespace,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Codespace returns a codespace of keeper
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

//-----------------------------------
// ExchangeRatePrevote logic

// GetExchangeRatePrevote retrieves an oracle prevote from the store
func (k Keeper) GetExchangeRatePrevote(ctx sdk.Context, denom string, voter sdk.ValAddress) (prevote types.ExchangeRatePrevote, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetExchangeRatePrevoteKey(denom, voter))
	if b == nil {
		err = types.ErrNoPrevote(k.codespace, voter, denom)
		return
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &prevote)
	return
}

// AddExchangeRatePrevote adds an oracle prevote to the store
func (k Keeper) AddExchangeRatePrevote(ctx sdk.Context, prevote types.ExchangeRatePrevote) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(prevote)
	store.Set(types.GetExchangeRatePrevoteKey(prevote.Denom, prevote.Voter), bz)
}

// DeleteExchangeRatePrevote deletes an oracle prevote from the store
func (k Keeper) DeleteExchangeRatePrevote(ctx sdk.Context, prevote types.ExchangeRatePrevote) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetExchangeRatePrevoteKey(prevote.Denom, prevote.Voter))
}

// IterateExchangeRatePrevotes iterates rate over prevotes in the store
func (k Keeper) IterateExchangeRatePrevotes(ctx sdk.Context, handler func(prevote types.ExchangeRatePrevote) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.PrevoteKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var prevote types.ExchangeRatePrevote
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &prevote)
		if handler(prevote) {
			break
		}
	}
}

// iterateExchangeRatePrevotesWithPrefix iterates over prevotes in the store with given prefix
func (k Keeper) iterateExchangeRatePrevotesWithPrefix(ctx sdk.Context, prefix []byte, handler func(vote types.ExchangeRatePrevote) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, prefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var prevote types.ExchangeRatePrevote
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &prevote)
		if handler(prevote) {
			break
		}
	}
}

//-----------------------------------
// ExchangeRateVotes logic

// IterateExchangeRateVotes iterates over votes in the store
func (k Keeper) IterateExchangeRateVotes(ctx sdk.Context, handler func(vote types.ExchangeRateVote) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.VoteKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var vote types.ExchangeRateVote
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &vote)
		if handler(vote) {
			break
		}
	}
}

// Iterate over oracle votes in the store
func (k Keeper) iterateExchangeRateVotesWithPrefix(ctx sdk.Context, prefix []byte, handler func(vote types.ExchangeRateVote) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, prefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var vote types.ExchangeRateVote
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &vote)
		if handler(vote) {
			break
		}
	}
}

// Retrieves an oracle vote from the store
func (k Keeper) getExchangeRateVote(ctx sdk.Context, denom string, voter sdk.ValAddress) (vote types.ExchangeRateVote, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetVoteKey(denom, voter))
	if b == nil {
		err = types.ErrNoVote(k.codespace, voter, denom)
		return
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &vote)
	return
}

// AddExchangeRateVote adds an oracle vote to the store
func (k Keeper) AddExchangeRateVote(ctx sdk.Context, vote types.ExchangeRateVote) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(vote)
	store.Set(types.GetVoteKey(vote.Denom, vote.Voter), bz)
}

// DeleteExchangeRateVote deletes an oracle vote from the store
func (k Keeper) DeleteExchangeRateVote(ctx sdk.Context, vote types.ExchangeRateVote) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetVoteKey(vote.Denom, vote.Voter))
}

//-----------------------------------
// ExchangeRate logic

// GetLunaExchangeRate gets the consensus exchange rate of Luna denominated in the denom asset from the store.
func (k Keeper) GetLunaExchangeRate(ctx sdk.Context, denom string) (exchangeRate sdk.Dec, err sdk.Error) {
	if denom == core.MicroLunaDenom {
		return sdk.OneDec(), nil
	}

	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetExchangeRateKey(denom))
	if b == nil {
		return sdk.ZeroDec(), types.ErrUnknownDenomination(k.codespace, denom)
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &exchangeRate)
	return
}

// SetLunaExchangeRate sets the consensus exchange rate of Luna denominated in the denom asset to the store.
func (k Keeper) SetLunaExchangeRate(ctx sdk.Context, denom string, exchangeRate sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(exchangeRate)
	store.Set(types.GetExchangeRateKey(denom), bz)
}

// DeleteLunaExchangeRate deletes the consensus exchange rate of Luna denominated in the denom asset from the store.
func (k Keeper) DeleteLunaExchangeRate(ctx sdk.Context, denom string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetExchangeRateKey(denom))
}

// IterateLunaExchangeRates iterates over luna rates in the store
func (k Keeper) IterateLunaExchangeRates(ctx sdk.Context, handler func(denom string, exchangeRate sdk.Dec) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ExchangeRateKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		denom := string(iter.Key()[len(types.ExchangeRateKey):])
		var exchangeRate sdk.Dec
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &exchangeRate)
		if handler(denom, exchangeRate) {
			break
		}
	}
}

//-----------------------------------
// Oracle delegation logic

// GetOracleDelegate gets the account address that the validator operator delegated oracle vote rights to
func (k Keeper) GetOracleDelegate(ctx sdk.Context, operator sdk.ValAddress) (delegate sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetFeederDelegationKey(operator))
	if b == nil {
		// By default the right is delegated to the validator itself
		return sdk.AccAddress(operator)
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &delegate)
	return
}

// SetOracleDelegate sets the account address that the validator operator delegated oracle vote rights to
func (k Keeper) SetOracleDelegate(ctx sdk.Context, operator sdk.ValAddress, delegatedFeeder sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(delegatedFeeder)
	store.Set(types.GetFeederDelegationKey(operator), bz)
}

// IterateOracleDelegates iterates over the feed delegates and performs a callback function.
func (k Keeper) IterateOracleDelegates(ctx sdk.Context,
	handler func(delegator sdk.ValAddress, delegate sdk.AccAddress) (stop bool)) {

	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.FeederDelegationKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		delegator := sdk.ValAddress(iter.Key()[len(types.FeederDelegationKey):])

		var delegate sdk.AccAddress
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &delegate)

		if handler(delegator, delegate) {
			break
		}
	}
}

//-----------------------------------
// Reward pool logic

// GetRewardPool retrieves the balance of the oracle module account
func (k Keeper) GetRewardPool(ctx sdk.Context) sdk.Coins {
	acc := k.supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
	return acc.GetCoins()
}

//-----------------------------------
// Miss counter logic

// GetMissCounter retrieves the # of vote periods missed in this oracle slash window
func (k Keeper) GetMissCounter(ctx sdk.Context, operator sdk.ValAddress) (missCounter int64) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetMissCounterKey(operator))
	if b == nil {
		// By default the counter is zero
		return 0
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &missCounter)
	return
}

// SetMissCounter updates the # of vote periods missed in this oracle slash window
func (k Keeper) SetMissCounter(ctx sdk.Context, operator sdk.ValAddress, missCounter int64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(missCounter)
	store.Set(types.GetMissCounterKey(operator), bz)
}

// IterateMissCounters iterates over the miss counters and performs a callback function.
func (k Keeper) IterateMissCounters(ctx sdk.Context,
	handler func(operator sdk.ValAddress, missCounter int64) (stop bool)) {

	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.MissCounterKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		operator := sdk.ValAddress(iter.Key()[len(types.MissCounterKey):])

		var missCounter int64
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &missCounter)

		if handler(operator, missCounter) {
			break
		}
	}
}

//-----------------------------------
// AggregateExchangeRatePrevote logic

// GetAggregateExchangeRatePrevote retrieves an oracle prevote from the store
func (k Keeper) GetAggregateExchangeRatePrevote(ctx sdk.Context, voter sdk.ValAddress) (aggregatePrevote types.AggregateExchangeRatePrevote, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetAggregateExchangeRatePrevoteKey(voter))
	if b == nil {
		err = types.ErrNoAggregatePrevote(k.codespace, voter)
		return
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &aggregatePrevote)
	return
}

// AddAggregateExchangeRatePrevote adds an oracle aggregate prevote to the store
func (k Keeper) AddAggregateExchangeRatePrevote(ctx sdk.Context, aggregatePrevote types.AggregateExchangeRatePrevote) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(aggregatePrevote)
	store.Set(types.GetAggregateExchangeRatePrevoteKey(aggregatePrevote.Voter), bz)
}

// DeleteAggregateExchangeRatePrevote deletes an oracle prevote from the store
func (k Keeper) DeleteAggregateExchangeRatePrevote(ctx sdk.Context, aggregatePrevote types.AggregateExchangeRatePrevote) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetAggregateExchangeRatePrevoteKey(aggregatePrevote.Voter))
}

// IterateAggregateExchangeRatePrevotes iterates rate over prevotes in the store
func (k Keeper) IterateAggregateExchangeRatePrevotes(ctx sdk.Context, handler func(aggregatePrevote types.AggregateExchangeRatePrevote) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.AggregatePrevoteKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var aggregatePrevote types.AggregateExchangeRatePrevote
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &aggregatePrevote)
		if handler(aggregatePrevote) {
			break
		}
	}
}

//-----------------------------------
// AggregateExchangeRateVote logic

// GetAggregateExchangeRateVote retrieves an oracle prevote from the store
func (k Keeper) GetAggregateExchangeRateVote(ctx sdk.Context, voter sdk.ValAddress) (aggregateVote types.AggregateExchangeRateVote, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetAggregateExchangeRateVoteKey(voter))
	if b == nil {
		err = types.ErrNoAggregateVote(k.codespace, voter)
		return
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &aggregateVote)
	return
}

// AddAggregateExchangeRateVote adds an oracle aggregate prevote to the store
func (k Keeper) AddAggregateExchangeRateVote(ctx sdk.Context, aggregateVote types.AggregateExchangeRateVote) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(aggregateVote)
	store.Set(types.GetAggregateExchangeRateVoteKey(aggregateVote.Voter), bz)
}

// DeleteAggregateExchangeRateVote deletes an oracle prevote from the store
func (k Keeper) DeleteAggregateExchangeRateVote(ctx sdk.Context, aggregateVote types.AggregateExchangeRateVote) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetAggregateExchangeRateVoteKey(aggregateVote.Voter))
}

// IterateAggregateExchangeRateVotes iterates rate over prevotes in the store
func (k Keeper) IterateAggregateExchangeRateVotes(ctx sdk.Context, handler func(aggregateVote types.AggregateExchangeRateVote) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.AggregateVoteKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var aggregateVote types.AggregateExchangeRateVote
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &aggregateVote)
		if handler(aggregateVote) {
			break
		}
	}
}

// HashAggregateExchangeRateVote
func (k Keeper) HashAggregateExchangeRateVote(ctx sdk.Context, voter sdk.ValAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetAggregateExchangeRateVoteKey(voter))
}

// GetTobinTax return tobin tax for the denom
func (k Keeper) GetTobinTax(ctx sdk.Context, denom string) (tobinTax sdk.Dec, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetTobinTaxKey(denom))
	if bz == nil {
		err = types.ErrNoTobinTax(k.codespace, denom)
		return
	} else {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &tobinTax)
	}

	return
}

// SetTobinTax updates tobin tax for the denom
func (k Keeper) SetTobinTax(ctx sdk.Context, denom string, tobinTax sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(tobinTax)
	store.Set(types.GetTobinTaxKey(denom), bz)
}

// IterateTobinTaxs iterates rate over tobin taxes in the store
func (k Keeper) IterateTobinTaxes(ctx sdk.Context, handler func(denom string, tobinTax sdk.Dec) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.TobinTaxKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		denom := types.ExtractDenomFromTobinTaxKey(iter.Key())

		var tobinTax sdk.Dec
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &tobinTax)
		if handler(denom, tobinTax) {
			break
		}
	}
}

// ClearTobinTaxes clears tobin taxes
func (k Keeper) ClearTobinTaxes(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.TobinTaxKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}
