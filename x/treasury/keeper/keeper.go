package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	core "github.com/terra-money/core/types"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/terra-money/core/x/treasury/types"
)

// TaxPowerUpgradeHeight is when taxes are allowed to go into effect
// This will still need a parameter change proposal, but can be activated
// anytime after this height
const TaxPowerUpgradeHeight = 9346889

// Keeper of the treasury store
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        codec.BinaryCodec
	paramSpace paramstypes.Subspace

	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	marketKeeper  types.MarketKeeper
	stakingKeeper types.StakingKeeper
	distrKeeper   types.DistributionKeeper
	oracleKeeper  types.OracleKeeper

	distributionModuleName string
}

// NewKeeper creates a new treasury Keeper instance
func NewKeeper(cdc codec.BinaryCodec, storeKey sdk.StoreKey,
	paramSpace paramstypes.Subspace,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	marketKeeper types.MarketKeeper,
	oracleKeeper types.OracleKeeper,
	stakingKeeper types.StakingKeeper,
	distrKeeper types.DistributionKeeper,
	distributionModuleName string) Keeper {

	// ensure treasury module account is set
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	// ensure burn module account is set
	if addr := accountKeeper.GetModuleAddress(types.BurnModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.BurnModuleName))
	}

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:                    cdc,
		storeKey:               storeKey,
		paramSpace:             paramSpace,
		accountKeeper:          accountKeeper,
		bankKeeper:             bankKeeper,
		marketKeeper:           marketKeeper,
		oracleKeeper:           oracleKeeper,
		stakingKeeper:          stakingKeeper,
		distrKeeper:            distrKeeper,
		distributionModuleName: distributionModuleName,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetTaxRate loads the tax rate
func (k Keeper) GetTaxRate(ctx sdk.Context) sdk.Dec {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.TaxRateKey)
	if b == nil {
		return types.DefaultTaxRate
	}

	dp := sdk.DecProto{}
	k.cdc.MustUnmarshal(b, &dp)
	return dp.Dec
}

// SetTaxRate sets the tax rate
func (k Keeper) SetTaxRate(ctx sdk.Context, taxRate sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&sdk.DecProto{Dec: taxRate})
	store.Set(types.TaxRateKey, b)
}

// GetRewardWeight loads the reward weight
func (k Keeper) GetRewardWeight(ctx sdk.Context) sdk.Dec {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.RewardWeightKey)
	if b == nil {
		return types.DefaultRewardWeight
	}

	dp := sdk.DecProto{}
	k.cdc.MustUnmarshal(b, &dp)
	return dp.Dec
}

// SetRewardWeight sets the reward weight
func (k Keeper) SetRewardWeight(ctx sdk.Context, rewardWeight sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&sdk.DecProto{Dec: rewardWeight})
	store.Set(types.RewardWeightKey, b)
}

// SetTaxCap sets the tax cap denominated in integer units of the reference {denom}
func (k Keeper) SetTaxCap(ctx sdk.Context, denom string, cap sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&sdk.IntProto{Int: cap})
	store.Set(types.GetTaxCapKey(denom), bz)
}

// GetTaxCap gets the tax cap denominated in integer units of the reference {denom}
func (k Keeper) GetTaxCap(ctx sdk.Context, denom string) sdk.Int {
	currHeight := ctx.BlockHeight()
	// Allow tax cap for uluna
	if denom == core.MicroLunaDenom && currHeight < TaxPowerUpgradeHeight {
		return sdk.ZeroInt()
	}

	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetTaxCapKey(denom))
	if bz == nil {
		// if no tax-cap registered, return SDR tax-cap
		return k.TaxPolicy(ctx).Cap.Amount
	}

	ip := sdk.IntProto{}
	k.cdc.MustUnmarshal(bz, &ip)
	return ip.Int
}

// IterateTaxCap iterates all tax cap
func (k Keeper) IterateTaxCap(ctx sdk.Context, handler func(denom string, taxCap sdk.Int) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.TaxCapKey)

	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		denom := string(iter.Key()[len(types.TaxCapKey):])
		var ip sdk.IntProto
		k.cdc.MustUnmarshal(iter.Value(), &ip)

		if handler(denom, ip.Int) {
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
	proceeds = proceeds.Add(delta...)

	k.SetEpochTaxProceeds(ctx, proceeds)
}

// SetEpochTaxProceeds stores tax proceeds for the given epoch
func (k Keeper) SetEpochTaxProceeds(ctx sdk.Context, taxProceeds sdk.Coins) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshal(&types.EpochTaxProceeds{TaxProceeds: taxProceeds})
	store.Set(types.TaxProceedsKey, bz)
}

// PeekEpochTaxProceeds peeks the total amount of taxes that have been collected in the given epoch.
func (k Keeper) PeekEpochTaxProceeds(ctx sdk.Context) sdk.Coins {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.TaxProceedsKey)
	taxProceeds := types.EpochTaxProceeds{}
	if bz == nil {
		taxProceeds.TaxProceeds = sdk.Coins{}
	} else {
		k.cdc.MustUnmarshal(bz, &taxProceeds)
	}

	return taxProceeds.TaxProceeds
}

// RecordEpochInitialIssuance updates epoch initial issuance from supply keeper
func (k Keeper) RecordEpochInitialIssuance(ctx sdk.Context) {
	whitelist := k.oracleKeeper.Whitelist(ctx)

	totalSupply := make(sdk.Coins, len(whitelist)+1)
	totalSupply[0] = k.bankKeeper.GetSupply(ctx, core.MicroLunaDenom)

	for i, denom := range whitelist {
		totalSupply[i+1] = k.bankKeeper.GetSupply(ctx, denom.Name)
	}

	k.SetEpochInitialIssuance(ctx, totalSupply.Sort())
}

// SetEpochInitialIssuance stores epoch initial issuance
func (k Keeper) SetEpochInitialIssuance(ctx sdk.Context, issuance sdk.Coins) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshal(&types.EpochInitialIssuance{Issuance: issuance})
	store.Set(types.EpochInitialIssuanceKey, bz)
}

// GetEpochInitialIssuance returns epoch initial issuance
func (k Keeper) GetEpochInitialIssuance(ctx sdk.Context) sdk.Coins {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.EpochInitialIssuanceKey)

	initialIssuance := types.EpochInitialIssuance{}
	if bz == nil {
		initialIssuance.Issuance = sdk.Coins{}
	} else {
		k.cdc.MustUnmarshal(bz, &initialIssuance)
	}

	return initialIssuance.Issuance
}

// PeekEpochSeigniorage returns epoch seigniorage
func (k Keeper) PeekEpochSeigniorage(ctx sdk.Context) sdk.Int {
	epochIssuance := k.bankKeeper.GetSupply(ctx, core.MicroLunaDenom).Amount
	preEpochIssuance := k.GetEpochInitialIssuance(ctx).AmountOf(core.MicroLunaDenom)
	epochSeigniorage := preEpochIssuance.Sub(epochIssuance)

	if epochSeigniorage.IsNegative() {
		return sdk.ZeroInt()
	}

	return epochSeigniorage
}

// GetTR returns the tax rewards for the epoch
func (k Keeper) GetTR(ctx sdk.Context, epoch int64) sdk.Dec {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetTRKey(epoch))

	dp := sdk.DecProto{}
	if bz == nil {
		dp.Dec = sdk.ZeroDec()
	} else {
		k.cdc.MustUnmarshal(bz, &dp)
	}

	return dp.Dec
}

// SetTR stores the tax rewards for the epoch
func (k Keeper) SetTR(ctx sdk.Context, epoch int64, TR sdk.Dec) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshal(&sdk.DecProto{Dec: TR})
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
func (k Keeper) GetSR(ctx sdk.Context, epoch int64) sdk.Dec {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetSRKey(epoch))

	dp := sdk.DecProto{}
	if bz == nil {
		dp.Dec = sdk.ZeroDec()
	} else {
		k.cdc.MustUnmarshal(bz, &dp)
	}

	return dp.Dec
}

// SetSR stores the seigniorage rewards for the epoch
func (k Keeper) SetSR(ctx sdk.Context, epoch int64, SR sdk.Dec) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshal(&sdk.DecProto{Dec: SR})
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

// GetTSL returns the total staked luna for the epoch
func (k Keeper) GetTSL(ctx sdk.Context, epoch int64) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetTSLKey(epoch))

	ip := sdk.IntProto{}
	if bz == nil {
		ip.Int = sdk.ZeroInt()
	} else {
		k.cdc.MustUnmarshal(bz, &ip)
	}

	return ip.Int
}

// SetTSL stores the total staked luna for the epoch
func (k Keeper) SetTSL(ctx sdk.Context, epoch int64, TSL sdk.Int) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshal(&sdk.IntProto{Int: TSL})
	store.Set(types.GetTSLKey(epoch), bz)
}

// ClearTSLs delete all the total staked luna from the store
func (k Keeper) ClearTSLs(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)

	iter := sdk.KVStorePrefixIterator(store, types.TSLKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}
