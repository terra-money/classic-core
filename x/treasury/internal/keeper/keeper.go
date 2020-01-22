package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"

	core "github.com/terra-project/core/types"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/terra-project/core/x/treasury/internal/types"
)

// Keeper of the treasury store
type Keeper struct {
	cdc      *codec.Codec
	storeKey sdk.StoreKey

	paramSpace params.Subspace
	codespace  sdk.CodespaceType

	supplyKeeper  types.SupplyKeeper
	marketKeeper  types.MarketKeeper
	stakingKeeper types.StakingKeeper
	distrKeeper   types.DistributionKeeper

	oracleModuleName       string
	distributionModuleName string
}

// NewKeeper creates a new treasury Keeper instance
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, paramSpace params.Subspace,
	supplyKeeper types.SupplyKeeper, marketKeeper types.MarketKeeper,
	stakingKeeper types.StakingKeeper, distrKeeper types.DistributionKeeper,
	oracleModuleName string, distributionModuleName string, codespace sdk.CodespaceType) Keeper {

	return Keeper{
		cdc:                    cdc,
		storeKey:               storeKey,
		paramSpace:             paramSpace.WithKeyTable(ParamKeyTable()),
		codespace:              codespace,
		supplyKeeper:           supplyKeeper,
		marketKeeper:           marketKeeper,
		stakingKeeper:          stakingKeeper,
		distrKeeper:            distrKeeper,
		oracleModuleName:       oracleModuleName,
		distributionModuleName: distributionModuleName,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Codespace returns the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

// GetTaxRate loads the tax rate
func (k Keeper) GetTaxRate(ctx sdk.Context) (taxRate sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.TaxRateKey)
	if b == nil {
		return types.DefaultTaxRate
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &taxRate)
	return
}

// SetTaxRate sets the tax rate
func (k Keeper) SetTaxRate(ctx sdk.Context, taxRate sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(taxRate)
	store.Set(types.TaxRateKey, b)
}

// GetRewardWeight loads the reward weight
func (k Keeper) GetRewardWeight(ctx sdk.Context) (rewardWeight sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.RewardWeightKey)
	if b == nil {
		return types.DefaultRewardWeight
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &rewardWeight)
	return
}

// SetRewardWeight sets the reward weight
func (k Keeper) SetRewardWeight(ctx sdk.Context, rewardWeight sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(rewardWeight)
	store.Set(types.RewardWeightKey, b)
}

// SetTaxCap sets the tax cap denominated in integer units of the reference {denom}
func (k Keeper) SetTaxCap(ctx sdk.Context, denom string, cap sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(cap)
	store.Set(types.GetTaxCapKey(denom), bz)
}

// GetTaxCap gets the tax cap denominated in integer units of the reference {denom}
func (k Keeper) GetTaxCap(ctx sdk.Context, denom string) (taxCap sdk.Int) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetTaxCapKey(denom))
	if bz == nil {
		// if no tax-cap registered, return SDR tax-cap
		return k.TaxPolicy(ctx).Cap.Amount
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &taxCap)
	return
}

// IterateTaxCap iterates all tax cap
func (k Keeper) IterateTaxCap(ctx sdk.Context, handler func(denom string, taxCap sdk.Int) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.TaxCapKey)

	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		denom := string(iter.Key()[len(types.TaxCapKey):])
		var taxCap sdk.Int
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &taxCap)

		if handler(denom, taxCap) {
			break
		}
	}

	return
}

// RecordEpochTaxProceeds adds tax proceeds that have been added this epoch
func (k Keeper) RecordEpochTaxProceeds(ctx sdk.Context, delta sdk.Coins) {
	if delta.IsZero() {
		return
	}

	proceeds := k.PeekEpochTaxProceeds(ctx)
	proceeds = proceeds.Add(delta)

	k.SetEpochTaxProceeds(ctx, proceeds)
}

// SetEpochTaxProceeds stores tax proceeds for the given epoch
func (k Keeper) SetEpochTaxProceeds(ctx sdk.Context, taxProceeds sdk.Coins) {
	store := ctx.KVStore(k.storeKey)

	if taxProceeds.IsZero() {
		store.Delete(types.TaxProceedsKey)
	} else {
		bz := k.cdc.MustMarshalBinaryLengthPrefixed(taxProceeds)
		store.Set(types.TaxProceedsKey, bz)
	}
}

// PeekEpochTaxProceeds peeks the total amount of taxes that have been collected in the given epoch.
func (k Keeper) PeekEpochTaxProceeds(ctx sdk.Context) (res sdk.Coins) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.TaxProceedsKey)
	if bz == nil {
		res = sdk.Coins{}
	} else {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	}
	return
}

// RecordEpochInitialIssuance updates epoch initial issuance from supply keeper
func (k Keeper) RecordEpochInitialIssuance(ctx sdk.Context) {
	totalCoins := k.supplyKeeper.GetSupply(ctx).GetTotal()
	k.SetEpochInitialIssuance(ctx, totalCoins)
}

// SetEpochInitialIssuance stores epoch initial issuance
func (k Keeper) SetEpochInitialIssuance(ctx sdk.Context, issuance sdk.Coins) {
	store := ctx.KVStore(k.storeKey)

	if issuance.IsZero() {
		store.Delete(types.EpochInitialIssuanceKey)
	} else {
		bz := k.cdc.MustMarshalBinaryLengthPrefixed(issuance)
		store.Set(types.EpochInitialIssuanceKey, bz)
	}
}

// GetEpochInitialIssuance returns epoch initial issuance
func (k Keeper) GetEpochInitialIssuance(ctx sdk.Context) (res sdk.Coins) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.EpochInitialIssuanceKey)

	if bz == nil {
		res = sdk.Coins{}
	} else {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	}
	return
}

// PeekEpochSeigniorage returns epoch seigniorage
func (k Keeper) PeekEpochSeigniorage(ctx sdk.Context) sdk.Int {
	epochIssuance := k.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(core.MicroLunaDenom)
	preEpochIssuance := k.GetEpochInitialIssuance(ctx).AmountOf(core.MicroLunaDenom)
	epochSeigniorage := preEpochIssuance.Sub(epochIssuance)

	if epochSeigniorage.LT(sdk.ZeroInt()) {
		return sdk.ZeroInt()
	}

	return epochSeigniorage
}

// GetCumulatedHeight returns last block height of past chain
func (k Keeper) GetCumulatedHeight(ctx sdk.Context) (res int64) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.CumulatedHeightKey)

	if bz == nil {
		res = 0
	} else {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	}
	return
}

// SetCumulatedHeight sets cumulated block height of past chains
func (k Keeper) SetCumulatedHeight(ctx sdk.Context, cumulatedHeight int64) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(cumulatedHeight)
	store.Set(types.CumulatedHeightKey, b)
}

// GetTR returns the tax rewards for the epoch
func (k Keeper) GetTR(ctx sdk.Context, epoch int64) (res sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetTRKey(epoch))

	if bz == nil {
		res = sdk.ZeroDec()
	} else {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	}

	return
}

// SetTR stores the tax rewards for the epoch
func (k Keeper) SetTR(ctx sdk.Context, epoch int64, TR sdk.Dec) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(TR)
	store.Set(types.GetTRKey(epoch), bz)
}

// ClearTRs delete all tax rewards from the store
func (k Keeper) ClearTRs(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)

	iter := sdk.KVStorePrefixIterator(store, types.TRKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}

// GetSR returns the seigniorage rewards for the epoch
func (k Keeper) GetSR(ctx sdk.Context, epoch int64) (res sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetSRKey(epoch))

	if bz == nil {
		res = sdk.ZeroDec()
	} else {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	}

	return
}

// SetSR stores the seigniorage rewards for the epoch
func (k Keeper) SetSR(ctx sdk.Context, epoch int64, SR sdk.Dec) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(SR)
	store.Set(types.GetSRKey(epoch), bz)
}

// ClearSRs delete all seigniorage rewards from the store
func (k Keeper) ClearSRs(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)

	iter := sdk.KVStorePrefixIterator(store, types.SRKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}

// GetTSL returns the total saked luna for the epoch
func (k Keeper) GetTSL(ctx sdk.Context, epoch int64) (res sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetTSLKey(epoch))

	if bz == nil {
		res = sdk.ZeroInt()
	} else {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	}

	return
}

// SetTSL stores the total saked luna for the epoch
func (k Keeper) SetTSL(ctx sdk.Context, epoch int64, TSL sdk.Int) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(TSL)
	store.Set(types.GetTSLKey(epoch), bz)
}

// ClearTSLs delete all the total saked luna from the store
func (k Keeper) ClearTSLs(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)

	iter := sdk.KVStorePrefixIterator(store, types.TSLKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}
