package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/terra-money/core/x/market/internal/types"
)

// Keeper of the oracle store
type Keeper struct {
	cdc        *codec.Codec
	storeKey   sdk.StoreKey
	paramSpace params.Subspace

	oracleKeeper types.OracleKeeper
	SupplyKeeper types.SupplyKeeper
}

// NewKeeper constructs a new keeper for oracle
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey,
	paramspace params.Subspace, oracleKeeper types.OracleKeeper,
	supplyKeeper types.SupplyKeeper) Keeper {

	// ensure oracle module account is set
	if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	// set KeyTable if it has not already been set
	if !paramspace.HasKeyTable() {
		paramspace = paramspace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:          cdc,
		storeKey:     storeKey,
		paramSpace:   paramspace,
		oracleKeeper: oracleKeeper,
		SupplyKeeper: supplyKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetTerraPoolDelta returns the gap between the TerraPool and the BasePool
func (k Keeper) GetTerraPoolDelta(ctx sdk.Context) (delta sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.TerraPoolDeltaKey)
	if bz == nil {
		return sdk.ZeroDec()
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &delta)
	return
}

// SetTerraPoolDelta updates TerraPoolDelta which is gap between the TerraPool and the BasePool
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
	// regressionAmt cannot make delta zero
	delta = delta.Sub(regressionAmt)

	k.SetTerraPoolDelta(ctx, delta)
}
