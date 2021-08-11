package simulation

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/terra-money/core/x/wasm/types"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding market type.
func NewDecodeStore(cdc codec.Codec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.LastCodeIDKey):
			lastCodeIDA := binary.LittleEndian.Uint64(kvA.Value)
			lastCodeIDB := binary.LittleEndian.Uint64(kvB.Value)
			return fmt.Sprintf("lastCodeIDA: %d\nlastCodeIDB: %d", lastCodeIDA, lastCodeIDB)
		case bytes.Equal(kvA.Key[:1], types.LastInstanceIDKey):
			lastInstanceIDKeyA := binary.LittleEndian.Uint64(kvA.Value)
			lastInstanceIDKeyB := binary.LittleEndian.Uint64(kvB.Value)
			return fmt.Sprintf("lastInstanceIDA: %d\nlastInstanceIDB: %d", lastInstanceIDKeyA, lastInstanceIDKeyB)
		case bytes.Equal(kvA.Key[:1], types.CodeKey):
			var codeInfoA, codeInfoB types.CodeInfo
			cdc.MustUnmarshal(kvA.Value, &codeInfoA)
			cdc.MustUnmarshal(kvB.Value, &codeInfoB)
			return fmt.Sprintf("%v\n%v", codeInfoA, codeInfoB)
		case bytes.Equal(kvA.Key[:1], types.ContractInfoKey):
			var contractInfoA, contractInfoB types.ContractInfo
			cdc.MustUnmarshal(kvA.Value, &contractInfoA)
			cdc.MustUnmarshal(kvB.Value, &contractInfoB)
			return fmt.Sprintf("%v\n%v", contractInfoA, contractInfoB)
		case bytes.Equal(kvA.Key[:1], types.ContractStoreKey):
			return fmt.Sprintf("%v\n%v", kvA.Value, kvB.Value)
		default:
			panic(fmt.Sprintf("invalid wasm key prefix %X", kvA.Key[:1]))
		}
	}
}
