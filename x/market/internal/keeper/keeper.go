package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/terra-project/core/x/market/internal/types"
)

// Keeper of the oracle store
type Keeper struct {
	cdc        *codec.Codec
	storeKey   sdk.StoreKey
	paramSpace params.Subspace

	oracleKeeper types.OracleKeeper
	SupplyKeeper types.SupplyKeeper

	// codespace
	codespace sdk.CodespaceType
}

// NewKeeper constructs a new keeper for oracle
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey,
	paramspace params.Subspace, oracleKeeper types.OracleKeeper,
	supplyKeeper types.SupplyKeeper, codespace sdk.CodespaceType) Keeper {
	return Keeper{
		cdc:          cdc,
		storeKey:     storeKey,
		paramSpace:   paramspace.WithKeyTable(ParamKeyTable()),
		oracleKeeper: oracleKeeper,
		SupplyKeeper: supplyKeeper,
		codespace:    codespace,
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

// GetTerraPoolDelta returns the gap between TerraPool and BasePool
func (k Keeper) GetTerraPoolDelta(ctx sdk.Context) (delta sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.TerraPoolDeltaKey)
	if bz == nil {
		return sdk.ZeroDec()
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &delta)
	return
}

// SetTerraPoolDelta updates TerraPoolDelta which is gap between TerraPool and BasePool
func (k Keeper) SetTerraPoolDelta(ctx sdk.Context, delta sdk.Dec) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(delta)
	store.Set(types.TerraPoolDeltaKey, bz)
}

// ReplenishPools replenishes each pool(Terra,Luna) to BasePool
func (k Keeper) ReplenishPools(ctx sdk.Context) {
	delta := k.GetTerraPoolDelta(ctx)
	regressionAmt := delta.QuoInt64(k.PoolRecoveryPeriod(ctx))

	// Replenish terra pool towards base pool
	if delta.IsPositive() {
		delta = delta.Sub(regressionAmt)
		if delta.IsNegative() {
			delta = sdk.ZeroDec()
		}
	} else if delta.IsNegative() {
		delta = delta.Add(regressionAmt)
		if delta.IsPositive() {
			delta = sdk.ZeroDec()
		}
	}

	k.SetTerraPoolDelta(ctx, delta)
}
