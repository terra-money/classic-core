package ante

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/core/x/wasm/types"
)

// CountTXDecorator ante handler to count the tx position in a block.
type CountTXDecorator struct {
	storeKey sdk.StoreKey
}

// NewCountTXDecorator constructor
func NewCountTXDecorator(storeKey sdk.StoreKey) *CountTXDecorator {
	return &CountTXDecorator{storeKey: storeKey}
}

// AnteHandle handler stores a tx counter with current height encoded in the store to let the app handle
// global rollback behavior instead of keeping state in the handler itself.
// The ante handler passes the counter value via sdk.Context upstream. See `types.TXCounter(ctx)` to read the value.
// Simulations don't get a tx counter value assigned.
func (a CountTXDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	if simulate {
		return next(ctx, tx, simulate)
	}
	store := ctx.KVStore(a.storeKey)
	currentHeight := ctx.BlockHeight()

	var txCounter uint32 // start with 0
	// load counter when exists
	if bz := store.Get(types.TXCounterKey); bz != nil {
		lastHeight, val := decodeHeightCounter(bz)
		if currentHeight == lastHeight {
			// then use stored counter
			txCounter = val
		} // else use `0` from above to start with
	}
	// store next counter value for current height
	store.Set(types.TXCounterKey, encodeHeightCounter(currentHeight, txCounter+1))

	return next(types.WithTXCounter(ctx, txCounter), tx, simulate)
}

func encodeHeightCounter(height int64, counter uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, counter)
	return append(sdk.Uint64ToBigEndian(uint64(height)), b...)
}

func decodeHeightCounter(bz []byte) (int64, uint32) {
	return int64(sdk.BigEndianToUint64(bz[0:8])), binary.BigEndian.Uint32(bz[8:])
}
