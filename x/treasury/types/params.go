package types

import (
	"fmt"

	"gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	core "github.com/terra-money/core/types"
)

// Parameter keys
var (
	KeyTaxPolicy               = []byte("TaxPolicy")
	KeyRewardPolicy            = []byte("RewardPolicy")
	KeySeigniorageBurdenTarget = []byte("SeigniorageBurdenTarget")
	KeyMiningIncrement         = []byte("MiningIncrement")
	KeyWindowShort             = []byte("WindowShort")
	KeyWindowLong              = []byte("WindowLong")
	KeyWindowProbation         = []byte("WindowProbation")
)

// Default parameter values
var (
	DefaultTaxPolicy = PolicyConstraints{
		RateMin:       sdk.NewDecWithPrec(5, 4),                                             // 0.05%
		RateMax:       sdk.NewDecWithPrec(1, 2),                                             // 1%
		Cap:           sdk.NewCoin(core.MicroSDRDenom, sdk.OneInt().MulRaw(core.MicroUnit)), // 1 SDR Tax cap
		ChangeRateMax: sdk.NewDecWithPrec(25, 5),                                            // 0.025%
	}
	DefaultRewardPolicy = PolicyConstraints{
		RateMin:       sdk.NewDecWithPrec(5, 2),             // 5%
		RateMax:       sdk.NewDecWithPrec(50, 2),            // 50%
		ChangeRateMax: sdk.NewDecWithPrec(25, 3),            // 2.5%
		Cap:           sdk.NewCoin("unused", sdk.ZeroInt()), // UNUSED
	}
	DefaultSeigniorageBurdenTarget = sdk.NewDecWithPrec(67, 2)  // 67%
	DefaultMiningIncrement         = sdk.NewDecWithPrec(107, 2) // 1.07 mining increment; exponential growth
	DefaultWindowShort             = uint64(4)                  // a month
	DefaultWindowLong              = uint64(52)                 // a year
	DefaultWindowProbation         = uint64(12)                 // 3 month
	DefaultTaxRate                 = sdk.NewDecWithPrec(1, 3)   // 0.1%
	DefaultRewardWeight            = sdk.NewDecWithPrec(5, 2)   // 5%
)

var _ paramstypes.ParamSet = &Params{}

// DefaultParams creates default treasury module parameters
func DefaultParams() Params {
	return Params{
		TaxPolicy:               DefaultTaxPolicy,
		RewardPolicy:            DefaultRewardPolicy,
		SeigniorageBurdenTarget: DefaultSeigniorageBurdenTarget,
		MiningIncrement:         DefaultMiningIncrement,
		WindowShort:             DefaultWindowShort,
		WindowLong:              DefaultWindowLong,
		WindowProbation:         DefaultWindowProbation,
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
// pairs of treasury module's parameters.
// nolint
func (p *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(KeyTaxPolicy, &p.TaxPolicy, validateTaxPolicy),
		paramstypes.NewParamSetPair(KeyRewardPolicy, &p.RewardPolicy, validateRewardPolicy),
		paramstypes.NewParamSetPair(KeySeigniorageBurdenTarget, &p.SeigniorageBurdenTarget, validateSeigniorageBurdenTarget),
		paramstypes.NewParamSetPair(KeyMiningIncrement, &p.MiningIncrement, validateMiningIncrement),
		paramstypes.NewParamSetPair(KeyWindowShort, &p.WindowShort, validateWindowShort),
		paramstypes.NewParamSetPair(KeyWindowLong, &p.WindowLong, validateWindowLong),
		paramstypes.NewParamSetPair(KeyWindowProbation, &p.WindowProbation, validateWindowProbation),
	}
}

// Validate performs basic validation on treasury parameters.
func (p Params) Validate() error {
	if p.TaxPolicy.RateMax.LT(p.TaxPolicy.RateMin) {
		return fmt.Errorf("treasury TaxPolicy.RateMax %s must be greater than TaxPolicy.RateMin %s",
			p.TaxPolicy.RateMax, p.TaxPolicy.RateMin)
	}

	if p.TaxPolicy.RateMin.IsNegative() {
		return fmt.Errorf("treasury parameter TaxPolicy.RateMin must be zero or positive: %s", p.TaxPolicy.RateMin)
	}

	if !p.TaxPolicy.Cap.IsValid() {
		return fmt.Errorf("treasury parameter TaxPolicy.Cap is invalid")
	}

	if p.TaxPolicy.ChangeRateMax.IsNegative() {
		return fmt.Errorf("treasury parameter TaxPolicy.ChangeRateMax must be positive: %s", p.TaxPolicy.ChangeRateMax)
	}

	if p.RewardPolicy.RateMax.LT(p.RewardPolicy.RateMin) {
		return fmt.Errorf("treasury RewardPolicy.RateMax %s must be greater than RewardPolicy.RateMin %s",
			p.RewardPolicy.RateMax, p.RewardPolicy.RateMin)
	}

	if p.RewardPolicy.RateMin.IsNegative() {
		return fmt.Errorf("treasury parameter RewardPolicy.RateMin must be positive: %s", p.RewardPolicy.RateMin)
	}

	if p.RewardPolicy.ChangeRateMax.IsNegative() {
		return fmt.Errorf("treasury parameter RewardPolicy.ChangeRateMax must be positive: %s", p.RewardPolicy.ChangeRateMax)
	}

	if p.SeigniorageBurdenTarget.IsNegative() {
		return fmt.Errorf("treasury parameter SeigniorageBurdenTarget must be positive: %s", p.SeigniorageBurdenTarget)
	}

	if p.MiningIncrement.IsNegative() {
		return fmt.Errorf("treasury parameter MiningIncrement must be positive: %s", p.MiningIncrement)
	}

	if p.WindowLong <= p.WindowShort {
		return fmt.Errorf("treasury parameter WindowLong must be bigger than WindowShort: (%d, %d)", p.WindowLong, p.WindowShort)
	}

	return nil
}

func validateTaxPolicy(i interface{}) error {
	v, ok := i.(PolicyConstraints)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.RateMin.IsNegative() {
		return fmt.Errorf("rate min must be positive: %s", v)
	}

	if v.RateMax.LT(v.RateMin) {
		return fmt.Errorf("rate max must be bigger than rate min: %s", v)
	}

	if !v.Cap.IsValid() {
		return fmt.Errorf("cap is invalid: %s", v)
	}

	if v.ChangeRateMax.IsNegative() {
		return fmt.Errorf("max change rate must be positive: %s", v)
	}

	return nil
}

func validateRewardPolicy(i interface{}) error {
	v, ok := i.(PolicyConstraints)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.RateMin.IsNegative() {
		return fmt.Errorf("rate min must be positive: %s", v)
	}

	if v.RateMax.LT(v.RateMin) {
		return fmt.Errorf("rate max must be bigger than rate min: %s", v)
	}

	if v.ChangeRateMax.IsNegative() {
		return fmt.Errorf("max change rate must be positive: %s", v)
	}

	return nil
}

func validateSeigniorageBurdenTarget(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("seigniorage burden target must be positive: %s", v)
	}

	return nil
}

func validateMiningIncrement(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("mining increment must be positive: %s", v)
	}

	return nil
}

func validateWindowShort(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateWindowLong(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateWindowProbation(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}
