package budget

import (
	"strconv"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// StoreKey is string representation of the store key for budget
const StoreKey = "budget"

// nolint
type Keeper struct {
	key       sdk.StoreKey      // Key to our module's store
	codespace sdk.CodespaceType // Reserves space for error codes
	cdc       *codec.Codec      // Codec to encore/decode structs
	valset    sdk.ValidatorSet  // Needed to compute voting power.

	bk         bank.Keeper // Needed to handle deposits. This module only requires read/writes to Terra balance
	paramSpace params.Subspace
}

// NewKeeper crates a new keeper with write and read access
func NewKeeper(key sdk.StoreKey,
	cdc *codec.Codec,
	bk bank.Keeper,
	codespace sdk.CodespaceType,
	valset sdk.ValidatorSet,
	paramspace params.Subspace) Keeper {
	return Keeper{
		key:        key,
		cdc:        cdc,
		bk:         bk,
		valset:     valset,
		codespace:  codespace,
		paramSpace: paramspace.WithKeyTable(ParamKeyTable()),
	}
}

// Get the last used proposal ID
func (k Keeper) NewProgramID(ctx sdk.Context) (programID uint64) {
	store := ctx.KVStore(k.key)

	bz := store.Get(KeyNextProgramID)
	if bz == nil {
		programID = 0
	} else {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &programID)
		programID++
	}

	bz = k.cdc.MustMarshalBinaryLengthPrefixed(programID)
	store.Set(KeyNextProgramID, bz)
	return
}

// GetProgram gets the Program with the given id from the context.
func (k Keeper) GetProgram(ctx sdk.Context, programID uint64) (res Program, err sdk.Error) {
	store := ctx.KVStore(k.key)
	bz := store.Get(KeyProgram(programID))
	if bz == nil {
		err = ErrProgramNotFound(programID)
		return
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	return
}

// SetProgram sets a Program to the context
func (k Keeper) SetProgram(ctx sdk.Context, programID uint64, program Program) {
	store := ctx.KVStore(k.key)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(program)
	store.Set(KeyProgram(programID), bz)
}

// DeleteProgram deletes a program from the context
func (k Keeper) DeleteProgram(ctx sdk.Context, programID uint64) {
	store := ctx.KVStore(k.key)
	store.Delete(KeyProgram(programID))
}

// IteratePrograms iterates through programs in the store
func (k Keeper) IterateActivePrograms(ctx sdk.Context, handler func(uint64, Program) (stop bool)) {
	store := ctx.KVStore(k.key)
	iter := sdk.KVStorePrefixIterator(store, PrefixCandidateQueue)
	for ; iter.Valid(); iter.Next() {
		var programStoreKey string
		var program Program
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Key(), &programStoreKey)
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &program)

		if programID, err := strconv.Atoi(strings.Split(programStoreKey, ":")[1]); err == nil {
			if k.CandidateQueueHas(ctx, program, uint64(programID)) {
				continue
			}

			if handler(uint64(programID), program) {
				break
			}
		}

	}
	iter.Close()
}

// GetVote returns the given option of a Program stored in the keeper
// Used to check if an address already voted
func (k Keeper) GetVote(ctx sdk.Context, programID uint64, voter sdk.AccAddress) (res bool, err sdk.Error) {
	store := ctx.KVStore(k.key)
	bz := store.Get(KeyVote(programID, voter))
	if bz == nil {
		err = ErrVoteNotFound()
		return
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	return
}

// SetVote sets the vote option to the Program stored in the context store
func (k Keeper) SetVote(ctx sdk.Context, programID uint64, voter sdk.AccAddress, option bool) {
	store := ctx.KVStore(k.key)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(option)
	store.Set(KeyVote(programID, voter), bz)
}

func (k Keeper) ClearVotesForProgram(ctx sdk.Context, programID uint64) {
	store := ctx.KVStore(k.key)
	iter := sdk.KVStorePrefixIterator(store, PrefixVoteForProgram(programID))
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
	iter.Close()
}

// RefundDeposit refunds the deposit
func (k Keeper) RefundDeposit(ctx sdk.Context, programID uint64) (err sdk.Error) {
	program, err := k.GetProgram(ctx, programID)
	if err != nil {
		return err
	}
	_, _, err = k.bk.AddCoins(ctx, program.Submitter, sdk.Coins{program.Deposit})
	return
}

//______________________________________________________________________

// GetParams get oralce params from the global param store
func (k Keeper) GetParams(ctx sdk.Context) Params {
	var params Params
	k.paramSpace.Get(ctx, ParamStoreKeyParams, &params)
	return params
}

// SetParams set oracle params from the global param store
func (k Keeper) SetParams(ctx sdk.Context, params Params) {
	k.paramSpace.Set(ctx, ParamStoreKeyParams, &params)
}

// =====================================================
// ProgramQueues

// IterateMatureCandidates Returns an iterator for all the Programs in the Candidate Queue that expire by endTime
func (k Keeper) IterateMatureCandidates(ctx sdk.Context, endTime time.Time, handler func(uint64, Program) (stop bool)) {
	store := ctx.KVStore(k.key)
	iter := store.Iterator(PrefixCandidateQueue, sdk.PrefixEndBytes(PrefixCandidateQueueTime(endTime)))
	for ; iter.Valid(); iter.Next() {
		var programStoreKey string
		var program Program
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Key(), &programStoreKey)
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &program)

		if programID, err := strconv.Atoi(strings.Split(programStoreKey, ":")[1]); err == nil {
			if handler(uint64(programID), program) {
				break
			}
		}
	}
	iter.Close()
}

// CandidateQueueInsert Inserts a ProgramID into the Candidate Program queue at endTime
func (k Keeper) CandidateQueueInsert(ctx sdk.Context, program Program, programID uint64) {
	votingEndTime := program.getVotingEndTime(k.GetParams(ctx).VotePeriod)

	store := ctx.KVStore(k.key)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(programID)
	store.Set(KeyCandidate(votingEndTime, programID), bz)
}

// CandidateQueueHas Checks if a progrma exists in accordance with the given parameters
func (k Keeper) CandidateQueueHas(ctx sdk.Context, program Program, programID uint64) (res bool) {
	votingEndTime := program.getVotingEndTime(k.GetParams(ctx).VotePeriod)

	store := ctx.KVStore(k.key)
	bz := store.Get(KeyCandidate(votingEndTime, programID))
	return bz != nil
}

// CandidateQueueRemove removes a ProgramID from the Candidate Program Queue
func (k Keeper) CandidateQueueRemove(ctx sdk.Context, program Program, programID uint64) {
	votingEndTime := program.getVotingEndTime(k.GetParams(ctx).VotePeriod)

	store := ctx.KVStore(k.key)
	store.Delete(KeyCandidate(votingEndTime, programID))
}
