package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-project/core/x/wasm/internal/types"
)

// MaxContractSize defines maximum bytes size of a contract
func (k Keeper) MaxContractSize(ctx sdk.Context) (res uint64) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyMaxContractSize, &res)
	return
}

// MaxContractGas defines allowed maximum gas usage per each contract execution
func (k Keeper) MaxContractGas(ctx sdk.Context) (res uint64) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyMaxContractGas, &res)
	return
}

// MaxContractMsgSize defines maximum bytes size of a contract
func (k Keeper) MaxContractMsgSize(ctx sdk.Context) (res uint64) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyMaxContractMsgSize, &res)
	return
}

// GasMultiplier defines how many cosmwasm gas points = 1 sdk gas point
func (k Keeper) GasMultiplier(ctx sdk.Context) (res uint64) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyGasMultiplier, &res)
	return
}

// CompileCostPerByte defines how much SDK gas we charge *per byte* for compiling WASM code.
func (k Keeper) CompileCostPerByte(ctx sdk.Context) (res uint64) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyCompileCostPerByte, &res)
	return
}

// InstanceCost defines how much SDK gas we charge each time we load a WASM instance.
func (k Keeper) InstanceCost(ctx sdk.Context) (res uint64) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyInstanceCost, &res)
	return
}

// HumanizeCost defines how much SDK gas we charge each time we humanize adress.
func (k Keeper) HumanizeCost(ctx sdk.Context) (res uint64) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyInstanceCost, &res)
	return
}

// CanonicalizeCost defines how much SDK gas we charge each time we canonicalize adress.
func (k Keeper) CanonicalizeCost(ctx sdk.Context) (res uint64) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyInstanceCost, &res)
	return
}

// GetParams returns the total set of oracle parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of oracle parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
