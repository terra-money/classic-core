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

// GetTaxRate loads the tax-rate
func (k Keeper) GetTaxRate(ctx sdk.Context, epoch int64) (taxRate sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetTaxRateKey(epoch))
	if b == nil {
		return types.DefaultTaxRate
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &taxRate)
	return
}

// SetTaxRate sets the tax-rate
func (k Keeper) SetTaxRate(ctx sdk.Context, taxRate sdk.Dec) {
	epoch := core.GetEpoch(ctx)
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(taxRate)
	store.Set(types.GetTaxRateKey(epoch), b)
}

// GetRewardWeight loads the reward weight
func (k Keeper) GetRewardWeight(ctx sdk.Context, epoch int64) (rewardWeight sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetRewardWeightKey(epoch))
	if b == nil {
		return types.DefaultRewardWeight
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &rewardWeight)
	return
}

// SetRewardWeight sets the reward weight
func (k Keeper) SetRewardWeight(ctx sdk.Context, rewardWeight sdk.Dec) {
	epoch := core.GetEpoch(ctx)
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(rewardWeight)
	store.Set(types.GetRewardWeightKey(epoch), b)
}

// SetTaxCap sets the Tax Cap. Denominated in integer units of the reference {denom}
func (k Keeper) SetTaxCap(ctx sdk.Context, denom string, cap sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(cap)
	store.Set(types.GetTaxCapKey(denom), bz)
}

// GetTaxCap gets the Tax Cap. Denominated in integer units of the reference {denom}
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

// IterateTaxCap iterates tax caps
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

// RecordTaxProceeds adds tax proceeds that have been added this epoch
func (k Keeper) RecordTaxProceeds(ctx sdk.Context, delta sdk.Coins) {
	if delta.Empty() || delta.IsZero() {
		return
	}

	epoch := core.GetEpoch(ctx)
	proceeds := k.PeekTaxProceeds(ctx, epoch)
	proceeds = proceeds.Add(delta)

	k.SetTaxProceeds(ctx, epoch, proceeds)
}

// SetTaxProceeds stores tax proceeds for the given epoch
func (k Keeper) SetTaxProceeds(ctx sdk.Context, epoch int64, taxProceeds sdk.Coins) {
	if taxProceeds.Empty() || taxProceeds.IsZero() {
		return
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(taxProceeds)
	store.Set(types.GetTaxProceedsKey(epoch), bz)
}

// PeekTaxProceeds peeks the total amount of taxes that have been collected in the given epoch.
func (k Keeper) PeekTaxProceeds(ctx sdk.Context, epoch int64) (res sdk.Coins) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetTaxProceedsKey(epoch))
	if bz == nil {
		res = sdk.Coins{}
	} else {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	}
	return
}

// ClearTaxProceeds clear all taxProceeds
func (k Keeper) ClearTaxProceeds(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.TaxProceedsKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}

// RecordHistoricalIssuance update epoch issuance from supply keeper (historical)
func (k Keeper) RecordHistoricalIssuance(ctx sdk.Context) {
	epoch := core.GetEpoch(ctx)
	totalCoins := k.supplyKeeper.GetSupply(ctx).GetTotal()
	k.SetHistoricalIssuance(ctx, epoch, totalCoins)
}

// SetHistoricalIssuance stores epoch issuance
func (k Keeper) SetHistoricalIssuance(ctx sdk.Context, epoch int64, issuance sdk.Coins) {
	if issuance.Empty() || issuance.IsZero() {
		return
	}

	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(issuance)
	store.Set(types.GetHistoricalIssuanceKey(epoch), bz)
	return
}

// GetHistoricalIssuance returns epoch issuance
func (k Keeper) GetHistoricalIssuance(ctx sdk.Context, epoch int64) (res sdk.Coins) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetHistoricalIssuanceKey(epoch))

	if bz == nil {
		res = sdk.Coins{}
	} else {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	}
	return
}

// ClearHistoricalIssuance clear all taxProceeds
func (k Keeper) ClearHistoricalIssuance(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.HistoricalIssuanceKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}

// PeekEpochSeigniorage retursn epoch seigniorage
func (k Keeper) PeekEpochSeigniorage(ctx sdk.Context, epoch int64) sdk.Int {
	if epoch == 0 {
		return sdk.ZeroInt()
	}

	epochIssuance := k.GetHistoricalIssuance(ctx, epoch).AmountOf(core.MicroLunaDenom)
	if epochIssuance.IsZero() {
		epochIssuance = k.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(core.MicroLunaDenom)
	}

	preEpochIssuance := k.GetHistoricalIssuance(ctx, epoch-1).AmountOf(core.MicroLunaDenom)
	epochSeigniorage := preEpochIssuance.Sub(epochIssuance)

	if epochSeigniorage.LT(sdk.ZeroInt()) {
		return sdk.ZeroInt()
	}

	return epochSeigniorage
}
