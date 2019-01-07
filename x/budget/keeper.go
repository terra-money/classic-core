package budget

import (
	"terra/types/util"
	"terra/x/treasury"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// nolint
type Keeper struct {
	key       sdk.StoreKey      // Key to our module's store
	codespace sdk.CodespaceType // Reserves space for error codes
	cdc       *codec.Codec      // Codec to encore/decode structs
	valset    sdk.ValidatorSet  // Needed to compute voting power.

	bk bank.Keeper // Needed to handle deposits. This module onlyl requires read/writes to Atom balance
	tk treasury.Keeper
}

// NewKeeper crates a new keeper with write and read access
func NewKeeper(key sdk.StoreKey, cdc *codec.Codec, bk bank.Keeper, tk treasury.Keeper, codespace sdk.CodespaceType, valset sdk.ValidatorSet) Keeper {
	return Keeper{
		key:       key,
		cdc:       cdc,
		bk:        bk,
		tk:        tk,
		valset:    valset,
		codespace: codespace,
	}
}

// Get the last used proposal ID
func (keeper Keeper) NewProgramID(ctx sdk.Context) (programID uint64) {
	programID, err = util.Get(keeper.key, keeper.cdc, ctx, KeyNextProgramID).(uint64)
	if err != nil {
		programID = 0
	} else {
		programID++
	}

	util.Set(keeper.key, keeper.cdc, ctx, KeyNextProgramID, programID)
	return
}

// GetProgram gets the Program with the given id from the context.
func (k Keeper) GetProgram(ctx sdk.Context, programID uint64) Program {
	return util.Get(k.key, k.cdc, ctx, GenerateProgramKey(programID)).(Program)
}

// SetProgram sets a Program to the context
func (k Keeper) SetProgram(ctx sdk.Context, programID uint64, program Program) sdk.Error {
	util.Set(k.key, k.cdc, ctx, GenerateProgramKey(programID))
}

// DeleteProgram deletes a program from the context
func (k Keeper) DeleteProgram(ctx sdk.Context, programID uint64) sdk.Error {
	util.Delete(k.key, k.cdc, ctx, GenerateProgramKey(programID))
}

// GetVote returns the given option of a Program stored in the keeper
// Used to check if an address already voted
func (k Keeper) GetVote(ctx sdk.Context, programID uint64, voter sdk.AccAddress) (ProgramVote, sdk.Error) {
	return util.Get(k.key, k.cdc, ctx, GenerateProgramVoteKey(programID, voter))
}

// SetVote sets the vote option to the Program stored in the context store
func (k Keeper) SetVote(ctx sdk.Context, programID uint64, voter sdk.AccAddress, option ProgramVote) {
	util.Set(k.key, k.cdc, ctx, GenerateProgramVoteKey(programID, voter), option)
}

// SetVote sets the vote option to the Program stored in the context store
func (k Keeper) RefundDeposit(ctx sdk.Context, programID uint64) (err sdk.Error) {
	program := k.GetProgram(ctx, programID)
	_, _, err := k.bk.AddCoins(ctx, program.Submitter, program.Deposit)
	return
}

//______________________________________________________________________

// GetParams get oralce params from the global param store
func (keeper Keeper) GetParams(ctx sdk.Context) Params {
	var params Params
	keeper.paramSpace.Get(ctx, ParamStoreKeyParams, &params)
	return params
}

// SetParams set oracle params from the global param store
func (keeper Keeper) SetParams(ctx sdk.Context, params Params) {
	keeper.paramSpace.Set(ctx, ParamStoreKeyParams, &params)
}

// =====================================================
// ProgramQueues

// Returns an iterator for all the Programs in the Inactive Queue that expire by endTime
func (keeper Keeper) InactiveProgramQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	return store.Iterator(PrefixInactiveProgramQueue, sdk.PrefixEndBytes(PrefixInactiveProgramQueueTime(endTime)))
}

// Inserts a ProgramID into the inactive Program queue at endTime
func (keeper Keeper) InsertInactiveProgramQueue(ctx sdk.Context, endTime time.Time, programID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(programID)
	store.Set(KeyInactiveProgramQueueProgram(endTime, programID), bz)
}

// Checks if a progrma exists in accordance with the given parameters
func (keeper Keeper) ProgramExistsInactiveProgramQueue(ctx sdk.Context, endTime time.Time, programID uint64) (res bool) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(programID)
	_, err := store.Get(KeyInactiveProgramQueueProgram(endTime, programID), bz)
	res = (err != nil)
}

// removes a ProgramID from the Inactive Program Queue
func (keeper Keeper) RemoveFromInactiveProgramQueue(ctx sdk.Context, endTime time.Time, programID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(KeyInactiveProgramQueueProgram(endTime, programID))
}
