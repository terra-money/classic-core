package types

import (
	"fmt"

	"gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"

	core "github.com/terra-money/core/types"
)

// DefaultParamspace nolint
const DefaultParamspace = ModuleName

// Parameter keys
var (
	//Terra liquidity pool(usdr unit) made available per ${poolrecoveryperiod} (usdr unit)
	ParamStoreKeyBasePool = []byte("basepool")
	// The period required to recover BasePool
	ParamStoreKeyPoolRecoveryPeriod = []byte("poolrecoveryperiod")
	// Min spread
	ParamStoreKeyMinStabilitySpread = []byte("minstabilityspread")
)

// Default parameter values
var (
	DefaultBasePool           = sdk.NewDec(250000 * core.MicroUnit) // 250,000sdr = 250,000,000,000usdr
	DefaultPoolRecoveryPeriod = core.BlocksPerDay                   // 14,400
	DefaultMinStabilitySpread = sdk.NewDecWithPrec(2, 2)            // 2%
)

var _ params.ParamSet = &Params{}

// Params market parameters
type Params struct {
	BasePool           sdk.Dec `json:"base_pool" yaml:"base_pool"`
	PoolRecoveryPeriod int64   `json:"pool_recovery_period" yaml:"pool_recovery_period"`
	MinStabilitySpread sdk.Dec `json:"min_spread" yaml:"min_spread"`
}

// DefaultParams creates default market module parameters
func DefaultParams() Params {
	return Params{
		BasePool:           DefaultBasePool,
		PoolRecoveryPeriod: DefaultPoolRecoveryPeriod,
		MinStabilitySpread: DefaultMinStabilitySpread,
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

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of market module's parameters.
// nolint
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(ParamStoreKeyBasePool, &p.BasePool, validateBasePool),
		params.NewParamSetPair(ParamStoreKeyPoolRecoveryPeriod, &p.PoolRecoveryPeriod, validatePoolRecoveryPeriod),
		params.NewParamSetPair(ParamStoreKeyMinStabilitySpread, &p.MinStabilitySpread, validateMinStatbilitySpread),
	}
}

// ValidateBasic a set of params
func (p Params) ValidateBasic() error {
	if p.BasePool.IsNegative() {
		return fmt.Errorf("base pool should be positive or zero, is %s", p.BasePool)
	}
	if p.PoolRecoveryPeriod <= 0 {
		return fmt.Errorf("pool recovery period should be positive, is %d", p.PoolRecoveryPeriod)
	}
	if p.MinStabilitySpread.IsNegative() || p.MinStabilitySpread.GT(sdk.OneDec()) {
		return fmt.Errorf("market minimum stability spead should be a value between [0,1], is %s", p.MinStabilitySpread)
	}

	return nil
}

func validateBasePool(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("base pool must be positive or zero: %s", v)
	}

	return nil
}

func validatePoolRecoveryPeriod(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("pool recovery period must be positive: %d", v)
	}

	return nil
}

func validateMinStatbilitySpread(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("min spread must be positive or zero: %s", v)
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("min spread is too large: %s", v)
	}

	return nil
}
