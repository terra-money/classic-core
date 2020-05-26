package types

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/cosmos/cosmos-sdk/x/params"
)

// DefaultParamspace defines default space for treasury params
const DefaultParamspace = ModuleName

// Max params for static check
const (
	EnforcedMaxContractSize    = uint64(500 * 1024)  // 500KB
	EnforcedMaxContractGas     = uint64(900_000_000) // 900,000,000
	EnforcedMaxContractMsgSize = uint64(10 * 1024)   // 10KB
)

// Parameter keys
var (
	ParamStoreKeyMaxContractSize    = []byte("maxcontractsize")
	ParamStoreKeyMaxContractGas     = []byte("maxcontractgas")
	ParamStoreKeyMaxContractMsgSize = []byte("maxcontractmsgsize")
	ParamStoreKeyGasMultiplier      = []byte("gasmultiplier")
)

// Default parameter values
const (
	DefaultMaxContractSize    = EnforcedMaxContractSize        // 500 KB
	DefaultMaxContractGas     = uint64(EnforcedMaxContractGas) // 900,000,000
	DefaultMaxContractMsgSize = uint64(1 * 1024)               // 1KB
	// SDK reference costs can be found here: https://github.com/cosmos/cosmos-sdk/blob/02c6c9fafd58da88550ab4d7d494724a477c8a68/store/types/gas.go#L153-L164
	// A write at ~3000 gas and ~200us = 10 gas per us (microsecond) cpu/io
	// Rough timing have 88k gas at 90us, which is equal to 1k sdk gas... (one read)
	DefaultGasMultiplier = uint64(100)
)

var _ params.ParamSet = &Params{}

// Params wasm parameters
type Params struct {
	MaxContractSize    uint64 `json:"max_contract_size" yaml:"max_contract_size"`         // allowed max contract bytes size
	MaxContractGas     uint64 `json:"max_contract_gas" yaml:"max_contract_gas"`           // allowed max gas usages per each contract execution
	MaxContractMsgSize uint64 `json:"max_contract_msg_size" yaml:"max_contract_msg_size"` // allowed max contract exe msg bytes size
	GasMultiplier      uint64 `json:"gas_multiplier" yaml:"gas_multiplier"`               // defines how many cosmwasm gas points = 1 sdk gas point
}

// DefaultParams creates default treasury module parameters
func DefaultParams() Params {
	return Params{
		MaxContractSize:    DefaultMaxContractSize,
		MaxContractGas:     DefaultMaxContractGas,
		MaxContractMsgSize: DefaultMaxContractMsgSize,
		GasMultiplier:      DefaultGasMultiplier,
	}
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of treasury module's parameters.
// nolint
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(ParamStoreKeyMaxContractSize, &p.MaxContractSize, validateMaxContractSize),
		params.NewParamSetPair(ParamStoreKeyMaxContractGas, &p.MaxContractGas, validateMaxContractGas),
		params.NewParamSetPair(ParamStoreKeyMaxContractMsgSize, &p.MaxContractMsgSize, validateMaxContractMsgSize),
		params.NewParamSetPair(ParamStoreKeyGasMultiplier, &p.GasMultiplier, validateGasMultiplier),
	}
}

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// String implements fmt.Stringer interface
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Validate params
func (p Params) Validate() error {
	if p.MaxContractSize > EnforcedMaxContractSize {
		return fmt.Errorf("max contract byte size %d must be equal or smaller than 500KB", p.MaxContractSize)
	}

	if p.GasMultiplier <= 0 {
		return fmt.Errorf("gas multiplier %d must be positive", p.GasMultiplier)
	}

	if p.MaxContractGas > EnforcedMaxContractGas {
		return fmt.Errorf("max contract gas %d must be equal or smaller than 900,000,000 (enforced in rust)", p.MaxContractGas)
	}

	if p.MaxContractMsgSize > EnforcedMaxContractMsgSize {
		return fmt.Errorf("max contract msg byte size %d must be equal or smaller than 10KB", p.MaxContractMsgSize)
	}

	return nil
}

func validateMaxContractSize(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v > EnforcedMaxContractSize {
		return fmt.Errorf("max contract byte size %d must be equal or smaller than 500KB", v)
	}

	return nil
}

func validateGasMultiplier(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("gas multiplier %d must be positive", v)
	}

	return nil
}

func validateMaxContractGas(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v > EnforcedMaxContractGas {
		return fmt.Errorf("max contract gas %d must be equal or smaller than 900,000,000 (enforced in rust)", v)
	}

	return nil
}

func validateMaxContractMsgSize(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v > EnforcedMaxContractMsgSize {
		return fmt.Errorf("max contract msg byte size %d must be equal or smaller than 10KB", v)
	}

	return nil
}
