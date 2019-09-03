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

// return the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

// Load the tax-rate
func (k Keeper) GetTaxRate(ctx sdk.Context, epoch int64) (taxRate sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetTaxRateKey(epoch))
	if b == nil {
		return types.DefaultTaxRate
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &taxRate)
	return
}

// Set the tax-rate
func (k Keeper) SetTaxRate(ctx sdk.Context, taxRate sdk.Dec) {
	epoch := core.GetEpoch(ctx)
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(taxRate)
	store.Set(types.GetTaxRateKey(epoch), b)
}

// Load the reward weight
func (k Keeper) GetRewardWeight(ctx sdk.Context, epoch int64) (rewardWeight sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetRewardWeightKey(epoch))
	if b == nil {
		return types.DefaultRewardWeight
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &rewardWeight)
	return
}

// Set the reward weight
func (k Keeper) SetRewardWeight(ctx sdk.Context, rewardWeight sdk.Dec) {
	epoch := core.GetEpoch(ctx)
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(rewardWeight)
	store.Set(types.GetRewardWeightKey(epoch), b)
}

// setTaxCap sets the Tax Cap. Denominated in integer units of the reference {denom}
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

// RecordTaxProceeds add tax proceeds that have been added this epoch
func (k Keeper) RecordTaxProceeds(ctx sdk.Context, delta sdk.Coins) {
	if delta.Empty() {
		return
	}

	epoch := core.GetEpoch(ctx)
	proceeds := k.PeekTaxProceeds(ctx, epoch)
	proceeds = proceeds.Add(delta)

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(proceeds)
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

// UpdateIssuance update epoch issuance from supply keeper (historical)
func (k Keeper) UpdateIssuance(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)

	epoch := core.GetEpoch(ctx)
	totalCoins := k.supplyKeeper.GetSupply(ctx).GetTotal()
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(totalCoins)
	store.Set(types.GetHistoricalIssuanceKey(epoch), bz)
}

// GetHistoricalIssuance returns epoch issuance of a denom
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
