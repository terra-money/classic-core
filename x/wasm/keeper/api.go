package keeper

import (
	"math"

	cosmwasm "github.com/CosmWasm/wasmvm"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/core/x/wasm/types"
)

func (k Keeper) getCosmWasmAPI() cosmwasm.GoAPI {
	return cosmwasm.GoAPI{
		HumanAddress: func(canon []byte) (humanAddr string, usedGas uint64, err error) {
			err = sdk.VerifyAddressFormat(canon)
			if err != nil {
				return "", types.HumanizeWasmGasCost, err
			}

			return sdk.AccAddress(canon).String(), types.HumanizeWasmGasCost, nil
		},
		CanonicalAddress: func(human string) (canonicalAddr []byte, usedGas uint64, err error) {
			addr, err := sdk.AccAddressFromBech32(human)
			if err != nil {
				return nil, types.CanonicalizeWasmGasCost, err
			}

			return addr, types.CanonicalizeWasmGasCost, nil
		},
	}
}

// wasmGasMeter wraps the GasMeter from context and multiplies all reads by out defined multiplier
type wasmGasMeter struct {
	originalMeter sdk.GasMeter
	gasMultiplier uint64
}

var _ cosmwasm.GasMeter = wasmGasMeter{}

func (m wasmGasMeter) GasConsumed() sdk.Gas {
	return types.ToWasmVMGas(m.originalMeter.GasConsumed())
}

// return gas meter interface for wasm gas meter
func (k Keeper) getWasmVMGasMeter(ctx sdk.Context) wasmGasMeter {
	return wasmGasMeter{
		originalMeter: ctx.GasMeter(),
		gasMultiplier: types.GasMultiplier,
	}
}

// return remaining gas in wasm gas unit
func (k Keeper) getWasmVMGasRemaining(ctx sdk.Context) uint64 {
	meter := ctx.GasMeter()

	// avoid integer overflow
	if meter.IsOutOfGas() {
		return 0
	}

	// infinite gas meter with limit=0 and not out of gas
	if meter.Limit() == 0 {
		return math.MaxUint64
	}

	remaining := (meter.Limit() - meter.GasConsumed())
	return types.ToWasmVMGas(remaining)
}

// converts contract gas usage to sdk gas and consumes it
func (k Keeper) consumeWasmVMGas(ctx sdk.Context, wasmVMGas uint64, descriptor string) {
	consumed := types.FromWasmVMGas(wasmVMGas)
	ctx.GasMeter().ConsumeGas(consumed, descriptor)
}
