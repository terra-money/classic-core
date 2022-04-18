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
)

// Keys for wasm store
// Items are stored with the following key: values
//
// - 0x04<accAddress_Bytes>: ContractInfo
var (
	ContractInfoKey = []byte{0x04}
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
