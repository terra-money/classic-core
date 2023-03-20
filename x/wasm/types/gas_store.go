package types

import (
	"io"
	"time"

	"github.com/cosmos/cosmos-sdk/store/gaskv"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ storetypes.KVStore = &Store{}

// KVStore return new gas KVStore which fixed
// https://github.com/cosmos/cosmos-sdk/issues/10243
func KVStore(ctx sdk.Context, key sdk.StoreKey) storetypes.KVStore {
	if (ctx.ChainID() == "bombay-12" && ctx.BlockHeight() < 7_800_000) ||
		(ctx.ChainID() == "columbus-5" && ctx.BlockHeight() < 6_470_000) {
		return gaskv.NewStore(ctx.MultiStore().GetKVStore(key), ctx.GasMeter(), storetypes.KVGasConfig())
	}

	return NewStore(ctx.MultiStore().GetKVStore(key), ctx.GasMeter(), storetypes.KVGasConfig())
}

// Store applies gas tracking to an underlying KVStore. It implements the
// KVStore interface.
type Store struct {
	gasMeter  storetypes.GasMeter
	gasConfig storetypes.GasConfig
	parent    storetypes.KVStore
}

// NewStore returns a reference to a new GasKVStore.
func NewStore(parent storetypes.KVStore, gasMeter storetypes.GasMeter, gasConfig storetypes.GasConfig) *Store {
	kvs := &Store{
		gasMeter:  gasMeter,
		gasConfig: gasConfig,
		parent:    parent,
	}
	return kvs
}

// GetStoreType implements Store.
func (gs *Store) GetStoreType() storetypes.StoreType {
	return gs.parent.GetStoreType()
}

// Get implements KVStore.
func (gs *Store) Get(key []byte) (value []byte) {
	gs.gasMeter.ConsumeGas(gs.gasConfig.ReadCostFlat, storetypes.GasReadCostFlatDesc)
	value = gs.parent.Get(key)

	// TODO overflow-safe math?
	gs.gasMeter.ConsumeGas(gs.gasConfig.ReadCostPerByte*storetypes.Gas(len(key)), storetypes.GasReadPerByteDesc)
	gs.gasMeter.ConsumeGas(gs.gasConfig.ReadCostPerByte*storetypes.Gas(len(value)), storetypes.GasReadPerByteDesc)

	return value
}

// Set implements KVStore.
func (gs *Store) Set(key []byte, value []byte) {
	storetypes.AssertValidKey(key)
	storetypes.AssertValidValue(value)
	gs.gasMeter.ConsumeGas(gs.gasConfig.WriteCostFlat, storetypes.GasWriteCostFlatDesc)
	// TODO overflow-safe math?
	gs.gasMeter.ConsumeGas(gs.gasConfig.WriteCostPerByte*storetypes.Gas(len(key)), storetypes.GasWritePerByteDesc)
	gs.gasMeter.ConsumeGas(gs.gasConfig.WriteCostPerByte*storetypes.Gas(len(value)), storetypes.GasWritePerByteDesc)
	gs.parent.Set(key, value)
}

// Has implements KVStore.
func (gs *Store) Has(key []byte) bool {
	defer telemetry.MeasureSince(time.Now(), "store", "gaskv", "has")
	gs.gasMeter.ConsumeGas(gs.gasConfig.HasCost, storetypes.GasHasDesc)
	return gs.parent.Has(key)
}

// Delete implements KVStore.
func (gs *Store) Delete(key []byte) {
	defer telemetry.MeasureSince(time.Now(), "store", "gaskv", "delete")
	// charge gas to prevent certain attack vectors even though space is being freed
	gs.gasMeter.ConsumeGas(gs.gasConfig.DeleteCost, storetypes.GasDeleteDesc)
	gs.parent.Delete(key)
}

// Iterator implements the KVStore interface. It returns an iterator which
// incurs a flat gas cost for seeking to the first key/value pair and a variable
// gas cost based on the current value's length if the iterator is valid.
func (gs *Store) Iterator(start, end []byte) storetypes.Iterator {
	return gs.iterator(start, end, true)
}

// ReverseIterator implements the KVStore interface. It returns a reverse
// iterator which incurs a flat gas cost for seeking to the first key/value pair
// and a variable gas cost based on the current value's length if the iterator
// is valid.
func (gs *Store) ReverseIterator(start, end []byte) storetypes.Iterator {
	return gs.iterator(start, end, false)
}

// CacheWrap implements KVStore.
func (gs *Store) CacheWrap() storetypes.CacheWrap {
	panic("cannot CacheWrap a GasKVStore")
}

// CacheWrapWithTrace implements the KVStore interface.
func (gs *Store) CacheWrapWithTrace(_ io.Writer, _ storetypes.TraceContext) storetypes.CacheWrap {
	panic("cannot CacheWrapWithTrace a GasKVStore")
}

// CacheWrapWithListeners implements the CacheWrapper interface.
func (gs *Store) CacheWrapWithListeners(_ storetypes.StoreKey, _ []storetypes.WriteListener) storetypes.CacheWrap {
	panic("cannot CacheWrapWithListeners a GasKVStore")
}

func (gs *Store) iterator(start, end []byte, ascending bool) storetypes.Iterator {
	var parent storetypes.Iterator
	if ascending {
		parent = gs.parent.Iterator(start, end)
	} else {
		parent = gs.parent.ReverseIterator(start, end)
	}

	gi := newGasIterator(gs.gasMeter, gs.gasConfig, parent)
	gi.(*gasIterator).consumeSeekGas()

	return gi
}

type gasIterator struct {
	gasMeter  storetypes.GasMeter
	gasConfig storetypes.GasConfig
	parent    storetypes.Iterator
}

func newGasIterator(gasMeter storetypes.GasMeter, gasConfig storetypes.GasConfig, parent storetypes.Iterator) storetypes.Iterator {
	return &gasIterator{
		gasMeter:  gasMeter,
		gasConfig: gasConfig,
		parent:    parent,
	}
}

// Implements Iterator.
func (gi *gasIterator) Domain() (start []byte, end []byte) {
	return gi.parent.Domain()
}

// Implements Iterator.
func (gi *gasIterator) Valid() bool {
	return gi.parent.Valid()
}

// Next implements the Iterator interface. It seeks to the next key/value pair
// in the iterator. It incurs a flat gas cost for seeking and a variable gas
// cost based on the current value's length if the iterator is valid.
func (gi *gasIterator) Next() {
	gi.consumeSeekGas()
	gi.parent.Next()
}

// Key implements the Iterator interface. It returns the current key and it does
// not incur any gas cost.
func (gi *gasIterator) Key() (key []byte) {
	key = gi.parent.Key()
	return key
}

// Value implements the Iterator interface. It returns the current value and it
// does not incur any gas cost.
func (gi *gasIterator) Value() (value []byte) {
	value = gi.parent.Value()
	return value
}

// Implements Iterator.
func (gi *gasIterator) Close() error {
	return gi.parent.Close()
}

// Error delegates the Error call to the parent iterator.
func (gi *gasIterator) Error() error {
	return gi.parent.Error()
}

// consumeSeekGas consumes on each iteration step a flat gas cost and a variable gas cost
// based on the current value's length.
func (gi *gasIterator) consumeSeekGas() {
	if gi.Valid() {
		key := gi.Key()
		value := gi.Value()

		gi.gasMeter.ConsumeGas(gi.gasConfig.ReadCostPerByte*storetypes.Gas(len(key)), storetypes.GasValuePerByteDesc)
		gi.gasMeter.ConsumeGas(gi.gasConfig.ReadCostPerByte*storetypes.Gas(len(value)), storetypes.GasValuePerByteDesc)
	}
	gi.gasMeter.ConsumeGas(gi.gasConfig.IterNextCostFlat, storetypes.GasIterNextCostFlatDesc)
}
