package simulation

import (
	"bytes"
	"encoding/binary"
	"fmt"

	tmkv "github.com/tendermint/tendermint/libs/kv"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/terra-money/core/x/wasm/internal/types"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding distribution type
func DecodeStore(cdc *codec.Codec, kvA, kvB tmkv.Pair) string {
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
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &codeInfoA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &codeInfoB)
		return fmt.Sprintf("%v\n%v", codeInfoA, codeInfoB)
	case bytes.Equal(kvA.Key[:1], types.ContractInfoKey):
		var contractInfoA, contractInfoB types.ContractInfo
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &contractInfoA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &contractInfoB)
		return fmt.Sprintf("%v\n%v", contractInfoA, contractInfoB)
	case bytes.Equal(kvA.Key[:1], types.ContractStoreKey):
		var rawDataA, rawDataB []byte
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &rawDataA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &rawDataB)
		return fmt.Sprintf("%v\n%v", rawDataA, rawDataB)
	default:
		panic(fmt.Sprintf("invalid wasm key prefix %X", kvA.Key[:1]))
	}
}
