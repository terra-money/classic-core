package types

import (
	"fmt"

	core "github.com/terra-project/core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
		RateMax:       sdk.NewDecWithPrec(90, 2),            // 90%
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

// Validate params
func (params Params) Validate() error {
	if params.TaxPolicy.RateMax.LT(params.TaxPolicy.RateMin) {
		return fmt.Errorf("treasury TaxPolicy.RateMax %s must be greater than TaxPolicy.RateMin %s",
			params.TaxPolicy.RateMax, params.TaxPolicy.RateMin)
	}

	if params.TaxPolicy.RateMin.IsNegative() {
		return fmt.Errorf("treasury parameter TaxPolicy.RateMin must be >= 0, is %s", params.TaxPolicy.RateMin)
	}

	if params.RewardPolicy.RateMax.LT(params.RewardPolicy.RateMin) {
		return fmt.Errorf("treasury RewardPolicy.RateMax %s must be greater than RewardPolicy.RateMin %s",
			params.RewardPolicy.RateMax, params.RewardPolicy.RateMin)
	}

	if params.RewardPolicy.RateMin.IsNegative() {
		return fmt.Errorf("treasury parameter RewardPolicy.RateMin must be >= 0, is %s", params.RewardPolicy.RateMin)
	}

	return nil
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of treasury module's parameters.
// nolint
func (params *Params) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{
		{Key: ParamStoreKeyTaxPolicy, Value: &params.TaxPolicy},
		{Key: ParamStoreKeyRewardPolicy, Value: &params.RewardPolicy},
		{Key: ParamStoreKeySeigniorageBurdenTarget, Value: &params.SeigniorageBurdenTarget},
		{Key: ParamStoreKeyMiningIncrement, Value: &params.MiningIncrement},
		{Key: ParamStoreKeyWindowShort, Value: &params.WindowShort},
		{Key: ParamStoreKeyWindowLong, Value: &params.WindowLong},
		{Key: ParamStoreKeyWindowProbation, Value: &params.WindowProbation},
	}
}

// String implements fmt.Stringer interface
func (params Params) String() string {
	return fmt.Sprintf(`Treasury Params:
  Tax Policy        : { %s } 
  Reward Policy     : { %s }

  SeigniorageBurdenTarget : %s
  MiningIncrement         : %s

  WindowShort        : %d
  WindowLong         : %d
  `, params.TaxPolicy, params.RewardPolicy, params.SeigniorageBurdenTarget,
		params.MiningIncrement, params.WindowShort, params.WindowLong)
}
