package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

type key int

const (
	// ModuleName is the name of the wasm module
	ModuleName = "wasm"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// TStoreKey is the string transient store representation
	TStoreKey = "transient_" + ModuleName

	// QuerierRoute is the querier route for the wasm module
	QuerierRoute = ModuleName

	// RouterKey is the msg router key for the wasm module
	RouterKey = ModuleName

	// WasmVMQueryDepthContextKey context key to keep query depth
	WasmVMQueryDepthContextKey key = iota
)

// Keys for wasm store
// Items are stored with the following key: values
//
// - 0x01: uint64
//
// - 0x02: uint64
//
// - 0x03<uint64>: CodeInfo
//
// - 0x04<accAddress_Bytes>: ContractInfo
//
// - 0x05<accAddress_Bytes>: KVStore for contract
var (
	LastCodeIDKey     = []byte{0x01}
	LastInstanceIDKey = []byte{0x02}
	CodeKey           = []byte{0x03}
	ContractInfoKey   = []byte{0x04}
	ContractStoreKey  = []byte{0x05}
	TXCounterKey      = []byte{0x06}
)

// GetCodeInfoKey constructs the key of the WASM code info for the ID
func GetCodeInfoKey(codeID uint64) []byte {
	contractIDBz := sdk.Uint64ToBigEndian(codeID)
	return append(CodeKey, contractIDBz...)
}

// GetContractInfoKey returns the key of the WASM contract info for the contract address
func GetContractInfoKey(addr sdk.AccAddress) []byte {
	return append(ContractInfoKey, address.MustLengthPrefix(addr)...)
}

// GetContractStoreKey returns the store prefix for the WASM contract store
func GetContractStoreKey(addr sdk.AccAddress) []byte {
	return append(ContractStoreKey, address.MustLengthPrefix(addr)...)
}
