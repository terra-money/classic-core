package budget

import (
	"strconv"
	"strings"

	"github.com/terra-project/core/x/mint"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// nolint
type Keeper struct {
	cdc    *codec.Codec     // Codec to encore/decode structs
	key    sdk.StoreKey     // Key to our module's store
	valset sdk.ValidatorSet // Needed to compute voting power.

	mk         mint.Keeper // Needed to handle deposits. This module only requires read/writes to Terra balance
	paramSpace params.Subspace
}

// NewKeeper crates a new keeper
func NewKeeper(cdc *codec.Codec,
	key sdk.StoreKey,
	mk mint.Keeper,
	valset sdk.ValidatorSet,
	paramspace params.Subspace) Keeper {
	return Keeper{
		cdc:        cdc,
		key:        key,
		mk:         mk,
		valset:     valset,
		paramSpace: paramspace.WithKeyTable(paramKeyTable()),
	}
}

//-----------------------------------
// Vote logic

// GetVote returns the given option of a Program stored in the keeper
func (k Keeper) GetVote(ctx sdk.Context, programID uint64, voter sdk.AccAddress) (res bool, err sdk.Error) {
	store := ctx.KVStore(k.key)
	if bz := store.Get(keyVote(programID, voter)); bz != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	} else {
		err = ErrVoteNotFound()
	}
	return
}

// AddVote adds the vote option to the store
func (k Keeper) AddVote(ctx sdk.Context, programID uint64, voter sdk.AccAddress, option bool) {
	store := ctx.KVStore(k.key)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(option)
	store.Set(keyVote(programID, voter), bz)
}

// DeleteVote deletes the vote from the store
func (k Keeper) DeleteVote(ctx sdk.Context, programID uint64, voter sdk.AccAddress) {
	store := ctx.KVStore(k.key)
	store.Delete(keyVote(programID, voter))
}

// DeleteVotesForProgram deletes the vote from the store
func (k Keeper) DeleteVotesForProgram(ctx sdk.Context, programID uint64) {
	store := ctx.KVStore(k.key)
	iter := sdk.KVStorePrefixIterator(store, keyVote(programID, sdk.AccAddress{}))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}

// IterateVotes iterates votes in the store
func (k Keeper) IterateVotes(ctx sdk.Context, handler func(uint64, sdk.AccAddress, bool) (stop bool)) {
	k.IterateVotesWithPrefix(ctx, prefixVote, handler)
}

// IterateVotesWithPrefix iterates votes with given {prefix} in the store
func (k Keeper) IterateVotesWithPrefix(ctx sdk.Context, prefix []byte, handler func(uint64, sdk.AccAddress, bool) (stop bool)) {
	store := ctx.KVStore(k.key)
	iter := sdk.KVStorePrefixIterator(store, prefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		voteKey := string(iter.Key())
		var option bool
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &option)

		elems := strings.Split(voteKey, ":")
		programID, err := strconv.ParseUint(elems[1], 10, 0)
		if err != nil {
			continue
		}

		voterAddrStr := elems[2]
		voterAddr, err := sdk.AccAddressFromBech32(voterAddrStr)
		if err != nil {
			continue
		}

		if handler(programID, voterAddr, option) {
			break
		}

	}
}

//-----------------------------------
// Deposit logic

// PayDeposit pays the deposit by withdrawing from the submitter's balance.
func (k Keeper) PayDeposit(ctx sdk.Context, submitter sdk.AccAddress) (err sdk.Error) {
	deposit := k.GetParams(ctx).Deposit
	err = k.mk.Burn(ctx, submitter, deposit)
	return
}

// RefundDeposit refunds the deposit, by crediting the submitter's balance.
func (k Keeper) RefundDeposit(ctx sdk.Context, submitter sdk.AccAddress) (err sdk.Error) {
	deposit := k.GetParams(ctx).Deposit
	err = k.mk.Mint(ctx, submitter, deposit)
	return
}

//-----------------------------------
// Params logic

// GetParams get budget params from the global param store
func (k Keeper) GetParams(ctx sdk.Context) Params {
	var resultParams Params
	k.paramSpace.Get(ctx, paramStoreKeyParams, &resultParams)
	return resultParams
}

// SetParams set budget params from the global param store
func (k Keeper) SetParams(ctx sdk.Context, params Params) {
	k.paramSpace.Set(ctx, paramStoreKeyParams, &params)
}

//-----------------------------------
// Program logic

// NewProgramID generates a new program id; advances sequentially from 1; 0 conflits with vote querier
func (k Keeper) NewProgramID(ctx sdk.Context) (programID uint64) {
	store := ctx.KVStore(k.key)
	if bz := store.Get(keyNextProgramID); bz != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &programID)
		programID++
	} else {
		programID = 1
	}

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(programID)
	store.Set(keyNextProgramID, bz)
	return
}

// GetProgram gets the Program with the given id from the store.
func (k Keeper) GetProgram(ctx sdk.Context, programID uint64) (res Program, err sdk.Error) {
	store := ctx.KVStore(k.key)

	if bz := store.Get(keyProgram(programID)); bz != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	} else {
		err = ErrProgramNotFound(programID)
	}
	return
}

// StoreProgram sets a Program to the store
func (k Keeper) StoreProgram(ctx sdk.Context, program Program) {
	store := ctx.KVStore(k.key)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(program)
	store.Set(keyProgram(program.ProgramID), bz)
}

// DeleteProgram deletes a program from the store
func (k Keeper) DeleteProgram(ctx sdk.Context, programID uint64) {
	store := ctx.KVStore(k.key)
	store.Delete(keyProgram(programID))
}

// IteratePrograms iterates programs in the store
func (k Keeper) IteratePrograms(ctx sdk.Context, filterInactive bool, handler func(Program) (stop bool)) {
	store := ctx.KVStore(k.key)
	iter := sdk.KVStorePrefixIterator(store, prefixProgram)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {

		var program Program
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &program)

		// Filter out candidate programs if filterInactive is true
		if filterInactive && k.CandQueueHas(ctx, program.getVotingEndBlock(ctx, k), program.ProgramID) {
			continue
		}

		if handler(program) {
			break
		}

	}
}

//-----------------------------------
// Candidate Queue logic

// CandQueueIterate iterate all the Programs in the candidate queue
func (k Keeper) CandQueueIterate(ctx sdk.Context, handler func(uint64) (stop bool)) {
	store := ctx.KVStore(k.key)
	iter := sdk.KVStorePrefixIterator(store, prefixCandQueue)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var programID uint64
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &programID)

		if handler(programID) {
			break
		}
	}
}

// CandQueueIterateExpired iterate all the Programs in the candidate queue that have outspent their voteperiod
func (k Keeper) CandQueueIterateExpired(ctx sdk.Context, endBlock int64, handler func(uint64) (stop bool)) {
	store := ctx.KVStore(k.key)
	iter := store.Iterator(prefixCandQueue, sdk.PrefixEndBytes(prefixCandQueueEndBlock(endBlock)))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var programID uint64
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &programID)

		if handler(programID) {
			break
		}
	}
}

// CandQueueInsert Inserts a ProgramID into the Candidate Program queue at endTime
func (k Keeper) CandQueueInsert(ctx sdk.Context, endBlock int64, programID uint64) {
	store := ctx.KVStore(k.key)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(programID)
	store.Set(keyCandidate(endBlock, programID), bz)
}

// CandQueueHas Checks if a progrma exists in accordance with the given parameters
func (k Keeper) CandQueueHas(ctx sdk.Context, endBlock int64, programID uint64) (res bool) {
	store := ctx.KVStore(k.key)
	bz := store.Get(keyCandidate(endBlock, programID))
	return bz != nil
}

// CandQueueRemove removes a ProgramID from the Candidate Program Queue
func (k Keeper) CandQueueRemove(ctx sdk.Context, endBlock int64, programID uint64) {
	store := ctx.KVStore(k.key)
	store.Delete(keyCandidate(endBlock, programID))
}
