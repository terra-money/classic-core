package types

import (
	"fmt"

	"gopkg.in/yaml.v2"

	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Max params for static check
const (
	EnforcedMaxContractSize     = uint64(500 * 1024)  // 500KB
	EnforcedMaxContractGas      = uint64(100_000_000) // 100,000,000
	EnforcedMaxContractMsgSize  = uint64(10 * 1024)   // 10KB
	EnforcedMaxContractDataSize = uint64(1024)        // 1KB
)

// Parameter keys
var (
	KeyMaxContractSize     = []byte("MaxContractSize")
	KeyMaxContractGas      = []byte("MaxContractGas")
	KeyMaxContractMsgSize  = []byte("MaxContractMsgSize")
	KeyMaxContractDataSize = []byte("MaxContractDataSize")
	KeyEventParams         = []byte("EventParams")
)

// Default parameter values
const (
	DefaultMaxContractSize     = EnforcedMaxContractSize // 500 KB
	DefaultMaxContractGas      = EnforcedMaxContractGas  // 100,000,000
	DefaultMaxContractMsgSize  = uint64(1 * 1024)        // 1KB
	DefaultMaxContractDataSize = uint64(256)             // 256 bytes
)

// Default event parameter values
var (
	DefaultEventParams = EventParams{
		MaxAttributeNum:         16,
		MaxAttributeKeyLength:   64,
		MaxAttributeValueLength: 256,
	}
)

// Constant gas parameters
const (
	InstanceCost       = uint64(40_000) // sdk gas cost for executing wasmVM engine
	CompileCostPerByte = uint64(2)      // sdk gas cost per bytes

	GasMultiplier    = uint64(100) // Please note that all gas prices returned to the wasmVM engine should have this multiplied
	HumanizeCost     = uint64(5)   // wasm gas cost to convert canonical address to human address
	CanonicalizeCost = uint64(4)   // wasm gas cost to convert human address to canonical address

	// ContractMemoryLimit is the memory limit of each contract execution (in MiB)
	// constant value so all nodes run with the same limit.
	ContractMemoryLimit = uint32(32)
)

var _ paramstypes.ParamSet = &Params{}

// DefaultParams creates default treasury module parameters
func DefaultParams() Params {
	return Params{
		MaxContractSize:     DefaultMaxContractSize,
		MaxContractGas:      DefaultMaxContractGas,
		MaxContractMsgSize:  DefaultMaxContractMsgSize,
		MaxContractDataSize: DefaultMaxContractDataSize,
		EventParams:         DefaultEventParams,
	}
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of treasury module's parameters.
// nolint
func (p *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(KeyMaxContractSize, &p.MaxContractSize, validateMaxContractSize),
		paramstypes.NewParamSetPair(KeyMaxContractGas, &p.MaxContractGas, validateMaxContractGas),
		paramstypes.NewParamSetPair(KeyMaxContractMsgSize, &p.MaxContractMsgSize, validateMaxContractMsgSize),
		paramstypes.NewParamSetPair(KeyMaxContractDataSize, &p.MaxContractDataSize, validateMaxContractDataSize),
		paramstypes.NewParamSetPair(KeyEventParams, &p.EventParams, validateEventParams),
	}
}

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

// String implements fmt.Stringer interface
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Validate params
func (p Params) Validate() error {
	if p.MaxContractSize > EnforcedMaxContractSize {
		return fmt.Errorf("max contract byte size %d must be equal or smaller than %d", p.MaxContractSize, EnforcedMaxContractSize)
	}

	if p.MaxContractGas > EnforcedMaxContractGas {
		return fmt.Errorf("max contract gas %d must be equal or smaller than %d", p.MaxContractGas, EnforcedMaxContractGas)
	}

	if p.MaxContractMsgSize > EnforcedMaxContractMsgSize {
		return fmt.Errorf("max contract msg byte size %d must be equal or smaller than %d", p.MaxContractMsgSize, EnforcedMaxContractMsgSize)
	}

	return nil
}

func validateMaxContractSize(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v > EnforcedMaxContractSize {
		return fmt.Errorf("max contract byte size %d must be equal or smaller than %d", v, EnforcedMaxContractSize)
	}

	return nil
}

func validateMaxContractGas(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v > EnforcedMaxContractGas {
		return fmt.Errorf("max contract gas %d must be equal or smaller than %d", v, EnforcedMaxContractGas)
	}

	return nil
}

func validateMaxContractMsgSize(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v > EnforcedMaxContractMsgSize {
		return fmt.Errorf("max contract msg byte size %d must be equal or smaller than %d", v, EnforcedMaxContractMsgSize)
	}

	return nil
}

func validateMaxContractDataSize(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v > EnforcedMaxContractDataSize {
		return fmt.Errorf("max contract data byte size %d must be equal or smaller than %d", v, EnforcedMaxContractDataSize)
	}

	return nil
}

func validateEventParams(i interface{}) error {
	_, ok := i.(EventParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}
