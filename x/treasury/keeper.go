package treasury

import (
	"github.com/terra-project/core/types/util"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Keeper of the treasury store
type Keeper struct {
	cdc *codec.Codec
	key sdk.StoreKey

	valset sdk.ValidatorSet

	mtk MintKeeper
	mk  MarketKeeper
	dk  DistributionKeeper
	fck FeeCollectionKeeper

	paramSpace params.Subspace
}

// NewKeeper constructs a new keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, valset sdk.ValidatorSet,
	mtk MintKeeper, mk MarketKeeper, dk DistributionKeeper, fck FeeCollectionKeeper, paramspace params.Subspace) Keeper {
	return Keeper{
		cdc:        cdc,
		key:        key,
		valset:     valset,
		mtk:        mtk,
		mk:         mk,
		dk:         dk,
		fck:        fck,
		paramSpace: paramspace.WithKeyTable(paramKeyTable()),
	}
}

//-----------------------------------
// Reward weight logic

// SetRewardWeight sets the ratio of the treasury that goes to mining rewards, i.e.
// supply of Luna that is burned. You can only set the reward weight of the current epoch.
func (k Keeper) SetRewardWeight(ctx sdk.Context, weight sdk.Dec) {
	store := ctx.KVStore(k.key)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(weight)
	store.Set(keyRewardWeight(util.GetEpoch(ctx)), bz)
}

// GetRewardWeight returns the mining reward weight
func (k Keeper) GetRewardWeight(ctx sdk.Context, epoch sdk.Int) (rewardWeight sdk.Dec) {
	store := ctx.KVStore(k.key)

	if bz := store.Get(keyRewardWeight(epoch)); bz != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &rewardWeight)
		return
	}

	for e := epoch; e.GTE(sdk.ZeroInt()); e = e.Sub(sdk.OneInt()) {
		if bz := store.Get(keyRewardWeight(e)); bz != nil {
			k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &rewardWeight)
			break
		} else if epoch.LTE(sdk.ZeroInt()) {
			// Genesis epoch; nothing exists in store so we set to default state
			rewardWeight = DefaultGenesisState().GenesisRewardWeight
			break
		}
	}

	// Set reward weight to the store
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(rewardWeight)
	store.Set(keyRewardWeight(epoch), bz)

	return
}

//-----------------------------------
// Params logic

// GetParams get treasury params from the global param store
func (k Keeper) GetParams(ctx sdk.Context) Params {
	var params Params
	k.paramSpace.Get(ctx, paramStoreKeyParams, &params)
	return params
}

// SetParams set treasury params from the global param store
func (k Keeper) SetParams(ctx sdk.Context, params Params) {
	k.paramSpace.Set(ctx, paramStoreKeyParams, &params)
}

//-----------------------------------
// Tax logic

// SetTaxRate sets the tax rate; called from the treasury.
func (k Keeper) SetTaxRate(ctx sdk.Context, rate sdk.Dec) {
	store := ctx.KVStore(k.key)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(rate)
	store.Set(keyTaxRate(util.GetEpoch(ctx)), bz)
}

// GetTaxRate gets the tax rate
func (k Keeper) GetTaxRate(ctx sdk.Context, epoch sdk.Int) (rate sdk.Dec) {
	store := ctx.KVStore(k.key)
	if bz := store.Get(keyTaxRate(epoch)); bz != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &rate)
	} else {
		if epoch.LTE(sdk.ZeroInt()) {
			rate = DefaultGenesisState().GenesisTaxRate
		} else {
			// Fetch the tax rate of the previous epoch
			rate = k.GetTaxRate(ctx, epoch.Sub(sdk.OneInt()))
		}

		// Set issuance to the store
		store := ctx.KVStore(k.key)
		bz := k.cdc.MustMarshalBinaryLengthPrefixed(rate)
		store.Set(keyTaxRate(epoch), bz)
	}
	return
}

// setTaxCap sets the Tax Cap. Denominated in integer units of the reference {denom}
func (k Keeper) setTaxCap(ctx sdk.Context, denom string, cap sdk.Int) {
	store := ctx.KVStore(k.key)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(cap)
	store.Set(keyTaxCap(denom), bz)
}

// GetTaxCap gets the Tax Cap. Denominated in integer units of the reference {denom}
func (k Keeper) GetTaxCap(ctx sdk.Context, denom string) (taxCap sdk.Int) {
	store := ctx.KVStore(k.key)

	if bz := store.Get(keyTaxCap(denom)); bz != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &taxCap)
	} else {
		// Tax cap does not exist for the asset; compute it by
		// comparing it with the tax cap for TerraSDR
		referenceCap := k.GetParams(ctx).TaxPolicy.Cap
		reqCap, _, err := k.mk.GetSwapCoins(ctx, referenceCap, denom, true)

		// The coin is more valuable than TaxPolicy asset. just follow the Policy Cap.
		if err != nil {
			reqCap = sdk.NewCoin(denom, referenceCap.Amount)
		}

		taxCap = reqCap.Amount
		k.setTaxCap(ctx, denom, taxCap)
	}

	return
}

// RecordTaxProceeds add tax proceeds that have been added this epoch
func (k Keeper) RecordTaxProceeds(ctx sdk.Context, delta sdk.Coins) {
	epoch := util.GetEpoch(ctx)
	proceeds := k.PeekTaxProceeds(ctx, epoch)
	proceeds = proceeds.Add(delta)

	store := ctx.KVStore(k.key)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(proceeds)
	store.Set(keyTaxProceeds(epoch), bz)
}

// PeekTaxProceeds peeks the total amount of taxes that have been collected in the given epoch.
func (k Keeper) PeekTaxProceeds(ctx sdk.Context, epoch sdk.Int) (res sdk.Coins) {
	store := ctx.KVStore(k.key)
	bz := store.Get(keyTaxProceeds(epoch))
	if bz == nil {
		res = sdk.Coins{}
	} else {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	}
	return
}
