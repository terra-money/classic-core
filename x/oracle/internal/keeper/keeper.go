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
// Prevote logic

// GetPrevote retrieves an oracle prevote from the store
func (k Keeper) GetPrevote(ctx sdk.Context, denom string, voter sdk.ValAddress) (prevote types.Prevote, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetPrevoteKey(denom, voter))
	if b == nil {
		err = types.ErrNoPrevote(k.codespace, voter, denom)
		return
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &prevote)
	return
}

// AddPrevote adds an oracle prevote to the store
func (k Keeper) AddPrevote(ctx sdk.Context, prevote types.Prevote) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(prevote)
	store.Set(types.GetPrevoteKey(prevote.Denom, prevote.Voter), bz)
}

// DeletePrevote deletes an oracle prevote from the store
func (k Keeper) DeletePrevote(ctx sdk.Context, prevote types.Prevote) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetPrevoteKey(prevote.Denom, prevote.Voter))
}

// IteratePrevotes iterates rate over prevotes in the store
func (k Keeper) IteratePrevotes(ctx sdk.Context, handler func(prevote types.Prevote) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.PrevoteKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var prevote types.Prevote
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &prevote)
		if handler(prevote) {
			break
		}
	}
}

// iteratePrevotesWithPrefix iterates over prevotes in the store with given prefix
func (k Keeper) iteratePrevotesWithPrefix(ctx sdk.Context, prefix []byte, handler func(vote types.Prevote) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, prefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var prevote types.Prevote
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &prevote)
		if handler(prevote) {
			break
		}
	}
}

//-----------------------------------
// Votes logic

// IterateVotes iterates over votes in the store
func (k Keeper) IterateVotes(ctx sdk.Context, handler func(vote types.Vote) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.VoteKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var vote types.Vote
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &vote)
		if handler(vote) {
			break
		}
	}
}

// Iterate over oracle votes in the store
func (k Keeper) iterateVotesWithPrefix(ctx sdk.Context, prefix []byte, handler func(vote types.Vote) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, prefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var vote types.Vote
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &vote)
		if handler(vote) {
			break
		}
	}
}

// Retrieves an oracle vote from the store
func (k Keeper) getVote(ctx sdk.Context, denom string, voter sdk.ValAddress) (vote types.Vote, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetVoteKey(denom, voter))
	if b == nil {
		err = types.ErrNoVote(k.codespace, voter, denom)
		return
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &vote)
	return
}

// AddVote adds an oracle vote to the store
func (k Keeper) AddVote(ctx sdk.Context, vote types.Vote) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(vote)
	store.Set(types.GetVoteKey(vote.Denom, vote.Voter), bz)
}

// DeleteVote deletes an oracle vote from the store
func (k Keeper) DeleteVote(ctx sdk.Context, vote types.Vote) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetVoteKey(vote.Denom, vote.Voter))
}

//-----------------------------------
// Price logic

// GetLunaExchangeRate gets the consensus exchange rate of Luna denominated in the denom asset from the store.
func (k Keeper) GetLunaExchangeRate(ctx sdk.Context, denom string) (exchangeRate sdk.Dec, err sdk.Error) {
	if denom == core.MicroLunaDenom {
		return sdk.OneDec(), nil
	}

	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetPriceKey(denom))
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
	store.Set(types.GetPriceKey(denom), bz)
}

// DeleteLunaExchangeRate deletes the consensus exchange rate of Luna denominated in the denom asset from the store.
func (k Keeper) DeleteLunaExchangeRate(ctx sdk.Context, denom string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetPriceKey(denom))
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
	handler func(delegator sdk.ValAddress, delegatee sdk.AccAddress) (stop bool)) {

	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.FeederDelegationKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		delegator := sdk.ValAddress(iter.Key()[len(types.FeederDelegationKey):])

		var delegatee sdk.AccAddress
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &delegatee)

		if handler(delegator, delegatee) {
			break
		}
	}
}

//-----------------------------------
// Reward pool logic

// getRewardPool retrieves the balance of the oracle module account
func (k Keeper) getRewardPool(ctx sdk.Context) sdk.Coins {
	acc := k.supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
	return acc.GetCoins()
}
