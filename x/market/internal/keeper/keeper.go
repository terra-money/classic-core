package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"

	core "github.com/terra-project/core/types"
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

// GetBasePool returns BasePool
func (k Keeper) GetBasePool(ctx sdk.Context) (pool sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.BasePoolKey)
	if bz == nil {
		return sdk.ZeroDec()
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &pool)
	return
}

// SetBasePool updates BasePool
func (k Keeper) SetBasePool(ctx sdk.Context, pool sdk.Dec) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(pool)
	store.Set(types.BasePoolKey, bz)
}

// GetLunaPool returns LunaPool
func (k Keeper) GetLunaPool(ctx sdk.Context) (pool sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.LunaPoolKey)
	if bz == nil {
		return sdk.ZeroDec()
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &pool)
	return
}

// SetLunaPool updates LunaPool
func (k Keeper) SetLunaPool(ctx sdk.Context, pool sdk.Dec) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(pool)
	store.Set(types.LunaPoolKey, bz)
}

// GetTerraPool returns TerraPool
func (k Keeper) GetTerraPool(ctx sdk.Context) (pool sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.TerraPoolKey)
	if bz == nil {
		return sdk.ZeroDec()
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &pool)
	return
}

// SetTerraPool updates TerraPool
func (k Keeper) SetTerraPool(ctx sdk.Context, pool sdk.Dec) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(pool)
	store.Set(types.TerraPoolKey, bz)
}

// GetLastUpdateHeight returns LastUpdateHeight
func (k Keeper) GetLastUpdateHeight(ctx sdk.Context) (height int64) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.LastUpdateHeightKey)
	if bz == nil {
		return 0
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &height)
	return
}

// SetLastUpdateHeight updates LastUpdateHeight
func (k Keeper) SetLastUpdateHeight(ctx sdk.Context, height int64) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(height)
	store.Set(types.LastUpdateHeightKey, bz)
}

// ReplenishPools replenishes each swap pool to base pool
func (k Keeper) ReplenishPools(ctx sdk.Context) {
	basePool := k.GetBasePool(ctx)
	terraPool := k.GetTerraPool(ctx)
	lunaPool := k.GetLunaPool(ctx)

	regressionAmt := basePool.QuoInt64(core.BlocksPerDay)

	// Replenish terra pool towards base pool
	if terraPool.GT(basePool) {
		terraPool = terraPool.Sub(regressionAmt)
		if terraPool.LT(basePool) {
			terraPool = basePool
		}
	} else if terraPool.LT(basePool) {
		terraPool = terraPool.Add(regressionAmt)
		if terraPool.GT(basePool) {
			terraPool = basePool
		}
	}

	// Replenish luna pool towards base pool
	if lunaPool.GT(basePool) {
		lunaPool = lunaPool.Sub(regressionAmt)
		if lunaPool.LT(basePool) {
			lunaPool = basePool
		}
	} else if lunaPool.LT(basePool) {
		lunaPool = lunaPool.Add(regressionAmt)
		if lunaPool.GT(basePool) {
			lunaPool = basePool
		}
	}

	k.SetTerraPool(ctx, terraPool)
	k.SetLunaPool(ctx, lunaPool)
}

// UpdatePools updates base & terra & luna pool along with sdr swap rate & luna supply
func (k Keeper) UpdatePools(ctx sdk.Context) (sdk.Dec, sdk.Error) {
	lunaSupplyAmt := k.SupplyKeeper.GetSupply(ctx).GetTotal().AmountOf(core.MicroLunaDenom)
	oldBasePool := k.GetBasePool(ctx)

	// swap luna supply to terra supply
	baseSupply, err := k.ComputeInternalSwap(ctx, sdk.NewDecCoin(core.MicroLunaDenom, lunaSupplyAmt), core.MicroSDRDenom)
	if err != nil {
		return sdk.ZeroDec(), err
	}

	basePool := k.DailyTerraLiquidityRatio(ctx).Mul(baseSupply.Amount)
	k.SetBasePool(ctx, basePool)
	k.SetLastUpdateHeight(ctx, ctx.BlockHeight())

	// Initial pool update
	if oldBasePool.IsZero() {
		k.SetLunaPool(ctx, basePool)
		k.SetTerraPool(ctx, basePool)
	} else {
		// Keep pool delta when updating
		oldLunaPool := k.GetLunaPool(ctx)
		oldTerraPool := k.GetTerraPool(ctx)

		// Reset each pools by multifying change ratio
		changeRatio := basePool.Quo(oldBasePool)
		lunaPool := oldLunaPool.Mul(changeRatio)
		terraPool := oldTerraPool.Mul(changeRatio)

		k.SetLunaPool(ctx, lunaPool)
		k.SetTerraPool(ctx, terraPool)
	}

	return basePool, nil
}

// IsMarketActive return current market activeness (check pool update was conducted or not in this interval period)
func (k Keeper) IsMarketActive(ctx sdk.Context) bool {
	height := ctx.BlockHeight()
	lastUpdateHeight := k.GetLastUpdateHeight(ctx)
	interval := k.PoolUpdateInterval(ctx)

	// Ative when base pool is positive and UpdateHeight is same or bigger than previous last block
	previousLastBlock := (height/interval)*interval - 1
	return k.GetBasePool(ctx).IsPositive() && lastUpdateHeight >= previousLastBlock
}
