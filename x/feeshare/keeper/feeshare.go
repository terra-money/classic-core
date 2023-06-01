package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/classic-terra/core/v2/x/feeshare/types"
)

// GetFeeShares returns all registered FeeShares.
func (k Keeper) GetFeeShares(ctx sdk.Context) []types.FeeShare {
	feeshares := []types.FeeShare{}

	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixFeeShare)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var feeshare types.FeeShare
		k.cdc.MustUnmarshal(iterator.Value(), &feeshare)

		feeshares = append(feeshares, feeshare)
	}

	return feeshares
}

// IterateFeeShares iterates over all registered contracts and performs a
// callback with the corresponding FeeShare.
func (k Keeper) IterateFeeShares(
	ctx sdk.Context,
	handlerFn func(fee types.FeeShare) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixFeeShare)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var feeshare types.FeeShare
		k.cdc.MustUnmarshal(iterator.Value(), &feeshare)

		if handlerFn(feeshare) {
			break
		}
	}
}

// GetFeeShare returns the FeeShare for a registered contract
func (k Keeper) GetFeeShare(
	ctx sdk.Context,
	contract sdk.Address,
) (types.FeeShare, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixFeeShare)
	bz := store.Get(contract.Bytes())
	if len(bz) == 0 {
		return types.FeeShare{}, false
	}

	var feeshare types.FeeShare
	k.cdc.MustUnmarshal(bz, &feeshare)
	return feeshare, true
}

// SetFeeShare stores the FeeShare for a registered contract.
func (k Keeper) SetFeeShare(ctx sdk.Context, feeshare types.FeeShare) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixFeeShare)
	key := feeshare.GetContractAddr()
	bz := k.cdc.MustMarshal(&feeshare)
	store.Set(key.Bytes(), bz)
}

// DeleteFeeShare deletes a FeeShare of a registered contract.
func (k Keeper) DeleteFeeShare(ctx sdk.Context, fee types.FeeShare) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixFeeShare)
	key := fee.GetContractAddr()
	store.Delete(key.Bytes())
}

// SetDeployerMap stores a contract-by-deployer mapping
func (k Keeper) SetDeployerMap(
	ctx sdk.Context,
	deployer sdk.AccAddress,
	contract sdk.Address,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixDeployer)
	key := append(deployer.Bytes(), contract.Bytes()...)
	store.Set(key, []byte{1})
}

// DeleteDeployerMap deletes a contract-by-deployer mapping
func (k Keeper) DeleteDeployerMap(
	ctx sdk.Context,
	deployer sdk.AccAddress,
	contract sdk.Address,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixDeployer)
	key := append(deployer.Bytes(), contract.Bytes()...)
	store.Delete(key)
}

// SetWithdrawerMap stores a contract-by-withdrawer mapping
func (k Keeper) SetWithdrawerMap(
	ctx sdk.Context,
	withdrawer sdk.AccAddress,
	contract sdk.Address,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixWithdrawer)
	key := append(withdrawer.Bytes(), contract.Bytes()...)
	store.Set(key, []byte{1})
}

// DeleteWithdrawMap deletes a contract-by-withdrawer mapping
func (k Keeper) DeleteWithdrawerMap(
	ctx sdk.Context,
	withdrawer sdk.AccAddress,
	contract sdk.Address,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixWithdrawer)
	key := append(withdrawer.Bytes(), contract.Bytes()...)
	store.Delete(key)
}

// IsFeeShareRegistered checks if a contract was registered for receiving
// transaction fees
func (k Keeper) IsFeeShareRegistered(
	ctx sdk.Context,
	contract sdk.Address,
) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixFeeShare)
	return store.Has(contract.Bytes())
}

// IsDeployerMapSet checks if a given contract-by-withdrawer mapping is set in
// store
func (k Keeper) IsDeployerMapSet(
	ctx sdk.Context,
	deployer sdk.AccAddress,
	contract sdk.Address,
) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixDeployer)
	key := append(deployer.Bytes(), contract.Bytes()...)
	return store.Has(key)
}

// IsWithdrawerMapSet checks if a give contract-by-withdrawer mapping is set in
// store
func (k Keeper) IsWithdrawerMapSet(
	ctx sdk.Context,
	withdrawer sdk.AccAddress,
	contract sdk.Address,
) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixWithdrawer)
	key := append(withdrawer.Bytes(), contract.Bytes()...)
	return store.Has(key)
}
