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
	EnforcedMaxContractSize    = uint64(500 * 1024)     // 500KB
	EnforcedMaxContractGas     = uint64(10_000_000_000) // 10,000,000,000
	EnforcedMaxContractMsgSize = uint64(10 * 1024)      // 10KB
)

// Parameter keys
var (
	ParamStoreKeyMaxContractSize    = []byte("maxcontractsize")
	ParamStoreKeyMaxContractGas     = []byte("maxcontractgas")
	ParamStoreKeyMaxContractMsgSize = []byte("maxcontractmsgsize")
	ParamStoreKeyGasMultiplier      = []byte("gasmultiplier")
	ParamStoreKeyCompileCostPerByte = []byte("compilecostperbyte")
	ParamStoreKeyInstanceCost       = []byte("instancecost")
	ParamStoreKeyHumanizeCost       = []byte("humanizecost")
	ParamStoreKeyCanonicalizeCost   = []byte("canonicalizecost")
)

// Default parameter values
const (
	DefaultMaxContractSize    = EnforcedMaxContractSize // 500 KB
	DefaultMaxContractGas     = EnforcedMaxContractGas  // 10,000,000,000
	DefaultMaxContractMsgSize = uint64(1 * 1024)        // 1KB
	DefaultGasMultiplier      = uint64(100)             // Please note that all gas prices returned to the wasmer engine should have this multiplied
	DefaultCompileCostPerByte = uint64(2)               // sdk gas cost per bytes
	DefaultInstanceCost       = uint64(40_000)          // sdk gas cost for executing wasmer engine
	DefaultHumanizeCost       = uint64(5)               // sdk gas cost to convert canonical address to human address
	DefaultCanonicalizeCost   = uint64(4)               // sdk gas cost to convert human address to canonical address
)

var _ params.ParamSet = &Params{}

// Params wasm parameters
type Params struct {
	MaxContractSize    uint64 `json:"max_contract_size" yaml:"max_contract_size"`         // allowed max contract bytes size
	MaxContractGas     uint64 `json:"max_contract_gas" yaml:"max_contract_gas"`           // allowed max gas usages per each contract execution
	MaxContractMsgSize uint64 `json:"max_contract_msg_size" yaml:"max_contract_msg_size"` // allowed max contract exe msg bytes size
	GasMultiplier      uint64 `json:"gas_multiplier" yaml:"gas_multiplier"`               // defines how many cosmwasm gas points = 1 sdk gas point
	CompileCostPerByte uint64 `json:"compile_cost_per_byte" yaml:"compile_cost_per_byte"` // defines how much SDK gas we charge *per byte* for compiling WASM code.
	InstanceCost       uint64 `json:"instance_cost" yaml:"instance_cost"`                 // defines how much SDK gas we charge each time we load a WASM instance.
	HumanizeCost       uint64 `json:"humanize_cost" yaml:"humanize_cost"`                 // defines how much SDK gas we charge each time we humanize address
	CanonicalizeCost   uint64 `json:"canonicalize_cost" yaml:"canonicalize_cost"`         // defines how much SDK gas we charge each time we canonicalize address
}

// DefaultParams creates default treasury module parameters
func DefaultParams() Params {
	return Params{
		MaxContractSize:    DefaultMaxContractSize,
		MaxContractGas:     DefaultMaxContractGas,
		MaxContractMsgSize: DefaultMaxContractMsgSize,
		GasMultiplier:      DefaultGasMultiplier,
		CompileCostPerByte: DefaultCompileCostPerByte,
		InstanceCost:       DefaultInstanceCost,
		HumanizeCost:       DefaultHumanizeCost,
		CanonicalizeCost:   DefaultCanonicalizeCost,
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
		params.NewParamSetPair(ParamStoreKeyCompileCostPerByte, &p.CompileCostPerByte, validateCompileCostPerByte),
		params.NewParamSetPair(ParamStoreKeyInstanceCost, &p.InstanceCost, validateInstanceCost),
		params.NewParamSetPair(ParamStoreKeyHumanizeCost, &p.HumanizeCost, validateHumanizeCost),
		params.NewParamSetPair(ParamStoreKeyCanonicalizeCost, &p.CanonicalizeCost, validateCanonicalizeCost),
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

func validateCompileCostPerByte(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateInstanceCost(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateHumanizeCost(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateCanonicalizeCost(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}
