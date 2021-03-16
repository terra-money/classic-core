package types

import (
	"fmt"

	"gopkg.in/yaml.v2"

	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Max params for static check
const (
	EnforcedMaxContractSize    = uint64(500 * 1024)  // 500KB
	EnforcedMaxContractGas     = uint64(100_000_000) // 100,000,000
	EnforcedMaxContractMsgSize = uint64(10 * 1024)   // 10KB
)

// Parameter keys
var (
	KeyMaxContractSize    = []byte("maxcontractsize")
	KeyMaxContractGas     = []byte("maxcontractgas")
	KeyMaxContractMsgSize = []byte("maxcontractmsgsize")
)

// Default parameter values
const (
	DefaultMaxContractSize    = EnforcedMaxContractSize // 500 KB
	DefaultMaxContractGas     = EnforcedMaxContractGas  // 100,000,000
	DefaultMaxContractMsgSize = uint64(1 * 1024)        // 1KB
)

// Constant gas parameters
const (
	GasMultiplier      = uint64(100)    // Please note that all gas prices returned to the wasmer engine should have this multiplied
	CompileCostPerByte = uint64(2)      // sdk gas cost per bytes
	InstanceCost       = uint64(40_000) // sdk gas cost for executing wasmer engine
	HumanizeCost       = uint64(5)      // sdk gas cost to convert canonical address to human address
	CanonicalizeCost   = uint64(4)      // sdk gas cost to convert human address to canonical address
)

var _ paramstypes.ParamSet = &Params{}

// DefaultParams creates default treasury module parameters
func DefaultParams() Params {
	return Params{
		MaxContractSize:    DefaultMaxContractSize,
		MaxContractGas:     DefaultMaxContractGas,
		MaxContractMsgSize: DefaultMaxContractMsgSize,
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
