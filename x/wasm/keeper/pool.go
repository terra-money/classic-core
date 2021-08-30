package keeper

import (
	"context"

	"github.com/terra-money/core/x/wasm/types"
)

func (k Keeper) getWasmVM(ctx context.Context) (types.WasmerEngine, error) {
	err := k.wasmReadVMSemaphore.Acquire(ctx, 1)
	if err != nil {
		return nil, err
	}

	wasmVM := k.wasmReadVMPool[0]
	k.wasmReadVMPool = k.wasmReadVMPool[1:]

	return wasmVM, nil
}

func (k Keeper) putWasmVM(wasmVM types.WasmerEngine) {
	k.wasmReadVMPool = append(k.wasmReadVMPool, wasmVM)
	k.wasmReadVMSemaphore.Release(1)
}
