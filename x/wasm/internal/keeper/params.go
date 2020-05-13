package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/terra-project/core/x/wasm/internal/types"
)

// ParamKeyTable returns ParamTable for wasm module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&types.Params{})
}

// MaxContractSize defines maximum bytes size of a contract
func (k Keeper) MaxContractSize(ctx sdk.Context) (res int64) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyMaxContractSize, &res)
	return
}

// MaxContractGas defines allowed maximum gas usage per each contract execution
func (k Keeper) MaxContractGas(ctx sdk.Context) (res uint64) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyMaxContractGas, &res)
	return
}

// GasMultiplier defines how many cosmwasm gas points = 1 sdk gas point
func (k Keeper) GasMultiplier(ctx sdk.Context) (res uint64) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyGasMultiplier, &res)
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
