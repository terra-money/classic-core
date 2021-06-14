package keeper

import (
	cosmwasm "github.com/CosmWasm/wasmvm"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/core/x/wasm/types"
)

func (k Keeper) getCosmWasmAPI(ctx sdk.Context) cosmwasm.GoAPI {
	return cosmwasm.GoAPI{
		HumanAddress: func(canon []byte) (humanAddr string, usedGas uint64, err error) {
			err = sdk.VerifyAddressFormat(canon)
			if err != nil {
				return "", 0, nil
			}

			return sdk.AccAddress(canon).String(), types.HumanizeCost * types.GasMultiplier, nil
		},
		CanonicalAddress: func(human string) (canonicalAddr []byte, usedGas uint64, err error) {
			addr, err := sdk.AccAddressFromBech32(human)
			if err != nil {
				return nil, 0, err
			}

			return addr, types.CanonicalizeCost * types.GasMultiplier, nil
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
	return m.originalMeter.GasConsumed() * m.gasMultiplier
}

// return gas meter interface for wasm gas meter
func (k Keeper) getGasMeter(ctx sdk.Context) wasmGasMeter {
	return wasmGasMeter{
		originalMeter: ctx.GasMeter(),
		gasMultiplier: types.GasMultiplier,
	}
}

// return remaining gas in wasm gas unit
func (k Keeper) getGasRemaining(ctx sdk.Context) uint64 {
	meter := ctx.GasMeter()

	// avoid integer overflow
	if meter.IsOutOfGas() {
		return 0
	}

	remaining := (meter.Limit() - meter.GasConsumed())
	return remaining * types.GasMultiplier
}

// converts contract gas usage to sdk gas and consumes it
func (k Keeper) consumeGas(ctx sdk.Context, gas uint64, descriptor string) {
	consumed := gas / types.GasMultiplier
	ctx.GasMeter().ConsumeGas(consumed, descriptor)
}
