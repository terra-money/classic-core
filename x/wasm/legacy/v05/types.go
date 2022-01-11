package v05

import (
	"gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName is the name of the wasm module
	ModuleName = "wasm"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	defaultMaxContractSize    = uint64(600 * 1024) // 600 KB
	defaultMaxContractGas     = uint64(20_000_000) // 20,000,000
	defaultMaxContractMsgSize = uint64(4 * 1024)   // 4KB
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
)

// String implements fmt.Stringer interface
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// GetContractInfoKey returns the key of the WASM contract info for the contract address
func GetContractInfoKey(addr sdk.AccAddress) []byte {
	return append(ContractInfoKey, address.MustLengthPrefix(addr)...)
}
