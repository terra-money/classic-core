package types

import (
	"fmt"

	"gopkg.in/yaml.v2"

	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Max params for static check
const (
	EnforcedMaxContractSize    = uint64(3000 * 1024) // 3MB
	EnforcedMaxContractGas     = uint64(100_000_000) // 100,000,000
	EnforcedMaxContractMsgSize = uint64(20 * 1024)   // 10KB
)

// Parameter keys
var (
	KeyMaxContractSize    = []byte("MaxContractSize")
	KeyMaxContractGas     = []byte("MaxContractGas")
	KeyMaxContractMsgSize = []byte("MaxContractMsgSize")
)

// Default parameter values
const (
	DefaultMaxContractSize    = uint64(600 * 1024) // 600 KB
	DefaultMaxContractGas     = uint64(20_000_000) // 20,000,000
	DefaultMaxContractMsgSize = uint64(4 * 1024)   // 4KB

	// ContractMemoryLimit is the memory limit of each contract execution (in MiB)
	// constant value so all nodes run with the same limit.
	ContractMemoryLimit = uint32(32)
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
