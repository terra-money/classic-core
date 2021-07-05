package keeper

import (
	"github.com/terra-money/core/x/wasm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MaxContractSize defines maximum bytes size of a contract
func (k Keeper) MaxContractSize(ctx sdk.Context) (res uint64) {
	k.paramSpace.Get(ctx, types.KeyMaxContractSize, &res)
	return
}

// MaxContractGas defines allowed maximum gas usage per each contract execution
func (k Keeper) MaxContractGas(ctx sdk.Context) (res uint64) {
	k.paramSpace.Get(ctx, types.KeyMaxContractGas, &res)
	return
}

// MaxContractMsgSize defines maximum bytes size of a contract
func (k Keeper) MaxContractMsgSize(ctx sdk.Context) (res uint64) {
	k.paramSpace.Get(ctx, types.KeyMaxContractMsgSize, &res)
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
