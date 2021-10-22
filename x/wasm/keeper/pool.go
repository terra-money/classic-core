package keeper

import (
	"context"

	"github.com/terra-money/core/x/wasm/types"
)

func (k Keeper) acquireWasmVM(ctx context.Context) (types.WasmerEngine, error) {
	err := k.wasmReadVMSemaphore.Acquire(ctx, 1)
	if err != nil {
		return nil, err
	}

	k.wasmReadVMMutex.Lock()
	wasmVM := (*k.wasmReadVMPool)[0]
	*k.wasmReadVMPool = (*k.wasmReadVMPool)[1:]
	k.wasmReadVMMutex.Unlock()

	return wasmVM, nil
}

func (k Keeper) releaseWasmVM(wasmVM types.WasmerEngine) {
	k.wasmReadVMMutex.Lock()
	*k.wasmReadVMPool = append(*k.wasmReadVMPool, wasmVM)
	k.wasmReadVMMutex.Unlock()

	k.wasmReadVMSemaphore.Release(1)
}
