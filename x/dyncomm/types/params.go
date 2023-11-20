package types

import (
	"fmt"

	"gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter keys
var (
	KeyMaxZero       = []byte("MaxZero")
	KeySlopeBase     = []byte("SlopeBase")
	KeySlopeVpImpact = []byte("SlopeVpImpact")
	KeyCap           = []byte("Cap")
)

// Default dyncomm parameter values
var (
	DefaultMaxZero       = sdk.NewDecWithPrec(5, 1)  // StrathColes A = 0.5
	DefaultSlopeBase     = sdk.NewDecWithPrec(2, 0)  // StrathColes B = 2
	DefaultSlopeVpImpact = sdk.NewDecWithPrec(10, 0) // StrathColes C = 10
	DefaultCap           = sdk.NewDecWithPrec(2, 1)  // StrathColes D = 20%
)

var _ paramstypes.ParamSet = &Params{}

// DefaultParams creates default dyncomm module parameters
func DefaultParams() Params {
	return Params{
		MaxZero:       DefaultMaxZero,
		SlopeBase:     DefaultSlopeBase,
		SlopeVpImpact: DefaultSlopeVpImpact,
		Cap:           DefaultCap,
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

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of market module's parameters.
func (p *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(KeyMaxZero, &p.MaxZero, validateMaxZero),
		paramstypes.NewParamSetPair(KeySlopeBase, &p.SlopeBase, validateSlopeBase),
		paramstypes.NewParamSetPair(KeySlopeVpImpact, &p.SlopeVpImpact, validateSlopeVpImpact),
		paramstypes.NewParamSetPair(KeyCap, &p.Cap, validateCap),
	}
}

// Validate a set of params
func (p Params) Validate() error {
	if p.SlopeBase.IsNegative() {
		return fmt.Errorf("slope base must be positive or zero, is %s", p.SlopeBase)
	}
	if !p.SlopeVpImpact.IsPositive() {
		return fmt.Errorf("solpe vp impact should be positive, is %d", p.SlopeVpImpact)
	}
	if p.Cap.IsNegative() {
		return fmt.Errorf("cap shall be 0 or positive: %s", p.Cap)
	}
	if p.Cap.GT(sdk.OneDec()) {
		return fmt.Errorf("cap shall be less than 1.0: %s", p.Cap)
	}
	if p.MaxZero.IsNegative() {
		return fmt.Errorf("max zero shall be 0 or positive: %s", p.MaxZero)
	}
	if p.MaxZero.GT(sdk.OneDec()) {
		return fmt.Errorf("max zero shall be less than 1.0: %s", p.MaxZero)
	}

	return nil
}

func validateMaxZero(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("max zero shall be 0 or positive: %s", v)
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("max zero shall be less than 1.0: %s", v)
	}

	return nil
}

func validateSlopeBase(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("slope base must be positive or Zero: %s", v)
	}

	return nil
}

func validateSlopeVpImpact(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if !v.IsPositive() {
		return fmt.Errorf("slope vp impact must be positive: %s", v)
	}

	return nil
}

func validateCap(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("cap shall be 0 or positive: %s", v)
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("cap shall be less than 1.0: %s", v)
	}

	return nil
}
