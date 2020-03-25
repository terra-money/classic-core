package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/params/subspace"
)

// DefaultParamspace defines default space for treasury params
const DefaultParamspace = ModuleName

// Parameter keys
var (
	ParamStoreKeyMaxContractSize = []byte("maxcontractsize")
	ParamStoreKeyMaxContractGas  = []byte("maxcontractgas")
	ParamStoreKeyGasMultiplier   = []byte("gasmultiplier")
)

// Default parameter values
var (
	DefaultMaxContractSize = int64(500 * 1024) // 500 KB
	DefaultMaxContractGas = uint64(900_000_000) // 900,000,000
	// SDK reference costs can be found here: https://github.com/cosmos/cosmos-sdk/blob/02c6c9fafd58da88550ab4d7d494724a477c8a68/store/types/gas.go#L153-L164
	// A write at ~3000 gas and ~200us = 10 gas per us (microsecond) cpu/io
	// Rough timing have 88k gas at 90us, which is equal to 1k sdk gas... (one read)
	DefaultGasMultiplier  = uint64(100)
)

var _ subspace.ParamSet = &Params{}

// Params wasm parameters
type Params struct {
	MaxContractSize int64  `json:"max_contract_size" yaml:"max_contract_size"` // allowed max contract bytes size
	MaxContractGas  uint64 `json:"max_contract_gas" yaml:"max_contract_gas"`   // allowed max gas usages per each contract execution
	GasMultiplier   uint64 `json:"gas_multiplier" yaml:"gas_multiplier"`       // defines how many cosmwasm gas points = 1 sdk gas point
}

// DefaultParams creates default treasury module parameters
func DefaultParams() Params {
	return Params{
		MaxContractSize: DefaultMaxContractSize,
		MaxContractGas:  DefaultMaxContractGas,
		GasMultiplier:   DefaultGasMultiplier,
	}
}

// Validate params
func (params Params) Validate() error {
	if params.MaxContractSize < 1024 || params.MaxContractSize > 400*1024 {
		return fmt.Errorf("max contract byte size %d must be between [1KB, 400KB]", params.MaxContractSize)
	}

	if params.GasMultiplier <= 0 {
		return fmt.Errorf("gas multiplier %d must be positive", params.GasMultiplier)
	}

	if params.MaxContractGas > 900_000_000 {
		return fmt.Errorf("max contract gas %d must be equal or smaller than 900,000,000 (enforced in rust)", params.MaxContractGas)
	}

	return nil
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of treasury module's parameters.
// nolint
func (params *Params) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{
		{Key: ParamStoreKeyMaxContractSize, Value: &params.MaxContractSize},
		{Key: ParamStoreKeyMaxContractGas, Value: &params.MaxContractGas},
		{Key: ParamStoreKeyGasMultiplier, Value: &params.GasMultiplier},
	}
}

// String implements fmt.Stringer interface
func (params Params) String() string {
	return fmt.Sprintf(`Treasury Params:
  Max Contract Size        : %d
  Max Contract Gas         : %d

  Gas Multiplier  : %d
  `, params.MaxContractSize, params.MaxContractGas, params.GasMultiplier)
}
