package types

import (
	"fmt"

	"gopkg.in/yaml.v2"

	core "github.com/terra-money/core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
)

// DefaultParamspace defines default space for treasury params
const DefaultParamspace = ModuleName

// Parameter keys
var (
	ParamStoreKeyTaxPolicy               = []byte("taxpolicy")
	ParamStoreKeyRewardPolicy            = []byte("rewardpolicy")
	ParamStoreKeySeigniorageBurdenTarget = []byte("seigniorageburdentarget")
	ParamStoreKeyMiningIncrement         = []byte("miningincrement")
	ParamStoreKeyWindowShort             = []byte("windowshort")
	ParamStoreKeyWindowLong              = []byte("windowlong")
	ParamStoreKeyWindowProbation         = []byte("windowprobation")
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
	DefaultWindowShort             = int64(4)                   // a month
	DefaultWindowLong              = int64(52)                  // a year
	DefaultWindowProbation         = int64(12)                  // 3 month
	DefaultTaxRate                 = sdk.NewDecWithPrec(1, 3)   // 0.1%
	DefaultRewardWeight            = sdk.NewDecWithPrec(5, 2)   // 5%
)

var _ subspace.ParamSet = &Params{}

// Params treasury parameters
type Params struct {
	TaxPolicy               PolicyConstraints `json:"tax_policy" yaml:"tax_policy"`
	RewardPolicy            PolicyConstraints `json:"reward_policy" yaml:"reward_policy"`
	SeigniorageBurdenTarget sdk.Dec           `json:"seigniorage_burden_target" yaml:"seigniorage_burden_target"`
	MiningIncrement         sdk.Dec           `json:"mining_increment" yaml:"mining_increment"`
	WindowShort             int64             `json:"window_short" yaml:"window_short"`
	WindowLong              int64             `json:"window_long" yaml:"window_long"`
	WindowProbation         int64             `json:"window_probation" yaml:"window_probation"`
}

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
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// String implements fmt.Stringer interface
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// ValidateBasic performs basic validation on treasury parameters.
func (p Params) ValidateBasic() error {
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

	if p.WindowShort < 0 {
		return fmt.Errorf("treasury parameter WindowShort must be positive: %d", p.WindowShort)
	}

	if p.WindowLong <= p.WindowShort {
		return fmt.Errorf("treasury parameter WindowLong must be bigger than WindowShort: (%d, %d)", p.WindowLong, p.WindowShort)
	}

	if p.WindowProbation < 0 {
		return fmt.Errorf("treasury parameter WindowProbation must be positive: %d", p.WindowProbation)
	}

	return nil
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of treasury module's parameters.
// nolint
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(ParamStoreKeyTaxPolicy, &p.TaxPolicy, validateTaxPolicy),
		params.NewParamSetPair(ParamStoreKeyRewardPolicy, &p.RewardPolicy, validateRewardPolicy),
		params.NewParamSetPair(ParamStoreKeySeigniorageBurdenTarget, &p.SeigniorageBurdenTarget, validateSeigniorageBurdenTarget),
		params.NewParamSetPair(ParamStoreKeyMiningIncrement, &p.MiningIncrement, validateMiningIncrement),
		params.NewParamSetPair(ParamStoreKeyWindowShort, &p.WindowShort, validateWindowShort),
		params.NewParamSetPair(ParamStoreKeyWindowLong, &p.WindowLong, validateWindowLong),
		params.NewParamSetPair(ParamStoreKeyWindowProbation, &p.WindowProbation, validateWindowProbation),
	}
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
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 0 {
		return fmt.Errorf("window short must be positive: %d", v)
	}

	return nil
}

func validateWindowLong(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 0 {
		return fmt.Errorf("window long must be positive: %d", v)
	}

	return nil
}

func validateWindowProbation(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 0 {
		return fmt.Errorf("window probation must be positive: %d", v)
	}

	return nil
}
