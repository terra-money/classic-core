package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/classic-terra/core/v2/x/dyncomm/types"
)

// Keeper of the market store
type Keeper struct {
	storeKey      storetypes.StoreKey
	cdc           codec.BinaryCodec
	paramSpace    paramstypes.Subspace
	StakingKeeper types.StakingKeeper
}

// NewKeeper constructs a new keeper for oracle
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	paramstore paramstypes.Subspace,
	stakingKeeper types.StakingKeeper,
) Keeper {
	// set KeyTable if it has not already been set
	if !paramstore.HasKeyTable() {
		paramstore = paramstore.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		paramSpace:    paramstore,
		StakingKeeper: stakingKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
