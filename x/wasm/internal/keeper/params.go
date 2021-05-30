package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-money/core/x/wasm/internal/types"
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

// GetParams returns the total set of oracle parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of oracle parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
