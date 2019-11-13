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

// RecordTaxProceeds adds tax proceeds that have been added this epoch
func (k Keeper) RecordTaxProceeds(ctx sdk.Context, delta sdk.Coins) {
	if delta.Empty() || delta.IsZero() {
		return
	}

	proceeds := k.PeekTaxProceeds(ctx)
	proceeds = proceeds.Add(delta)

	k.SetTaxProceeds(ctx, proceeds)
}

// SetTaxProceeds stores tax proceeds for the given epoch
func (k Keeper) SetTaxProceeds(ctx sdk.Context, taxProceeds sdk.Coins) {
	store := ctx.KVStore(k.storeKey)

	if taxProceeds.Empty() || taxProceeds.IsZero() {
		store.Delete(types.TaxProceedsKey)
	} else {
		bz := k.cdc.MustMarshalBinaryLengthPrefixed(taxProceeds)
		store.Set(types.TaxProceedsKey, bz)
	}
}

// PeekTaxProceeds peeks the total amount of taxes that have been collected in the given epoch.
func (k Keeper) PeekTaxProceeds(ctx sdk.Context) (res sdk.Coins) {
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

	if issuance.Empty() || issuance.IsZero() {
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
	if epochIssuance.IsZero() {
		epochIssuance = k.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(core.MicroLunaDenom)
	}

	preEpochIssuance := k.GetEpochInitialIssuance(ctx).AmountOf(core.MicroLunaDenom)
	epochSeigniorage := preEpochIssuance.Sub(epochIssuance)

	if epochSeigniorage.LT(sdk.ZeroInt()) {
		return sdk.ZeroInt()
	}

	return epochSeigniorage
}

// GetMR returns MR of the epoch
func (k Keeper) GetMR(ctx sdk.Context, epoch int64) (res sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetMRKey(epoch))

	if bz == nil {
		res = sdk.ZeroDec()
	} else {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	}

	return
}

// SetMR stores MR of the epoch
func (k Keeper) SetMR(ctx sdk.Context, epoch int64, MR sdk.Dec) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(MR)
	store.Set(types.GetMRKey(epoch), bz)
}

// ClearMRs delete all MRs from the store
func (k Keeper) ClearMRs(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)

	iter := sdk.KVStorePrefixIterator(store, types.MRKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}

// GetSR returns SR of the epoch
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

// SetSR stores SR of the epoch
func (k Keeper) SetSR(ctx sdk.Context, epoch int64, SR sdk.Dec) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(SR)
	store.Set(types.GetSRKey(epoch), bz)
}

// ClearSRs delete all SRs from the store
func (k Keeper) ClearSRs(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)

	iter := sdk.KVStorePrefixIterator(store, types.SRKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}

// GetTRL returns TRL of the epoch
func (k Keeper) GetTRL(ctx sdk.Context, epoch int64) (res sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetTRLKey(epoch))

	if bz == nil {
		res = sdk.ZeroDec()
	} else {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	}

	return
}

// SetTRL stores TRL of the epoch
func (k Keeper) SetTRL(ctx sdk.Context, epoch int64, TRL sdk.Dec) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(TRL)
	store.Set(types.GetTRLKey(epoch), bz)
}

// ClearTRLs delete all TRLs from the store
func (k Keeper) ClearTRLs(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)

	iter := sdk.KVStorePrefixIterator(store, types.TRLKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}
