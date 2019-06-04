package treasury

import (
	"fmt"

	"github.com/terra-project/core/types/assets"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Params treasury parameters
type Params struct {
	TaxPolicy    PolicyConstraints `json:"tax_policy"`
	RewardPolicy PolicyConstraints `json:"reward_policy"`

	SeigniorageBurdenTarget sdk.Dec `json:"seigniorage_burden_target"`
	MiningIncrement         sdk.Dec `json:"mining_increment"`

	WindowShort     sdk.Int `json:"window_short"`
	WindowLong      sdk.Int `json:"window_long"`
	WindowProbation sdk.Int `json:"window_probation"`
}

// NewParams creates a new param instance
func NewParams(
	taxPolicy, rewardPolicy PolicyConstraints,
	seigniorageBurden sdk.Dec,
	miningIncrement sdk.Dec,
	windowShort, windowLong, windowProbation sdk.Int,
) Params {
	return Params{
		TaxPolicy:               taxPolicy,
		RewardPolicy:            rewardPolicy,
		SeigniorageBurdenTarget: seigniorageBurden,
		MiningIncrement:         miningIncrement,
		WindowShort:             windowShort,
		WindowLong:              windowLong,
		WindowProbation:         windowProbation,
	}
}

// DefaultParams creates default treasury module parameters
func DefaultParams() Params {
	return NewParams(

		// Tax update policy
		PolicyConstraints{
			RateMin:       sdk.NewDecWithPrec(5, 4),                                                 // 0.05%
			RateMax:       sdk.NewDecWithPrec(1, 2),                                                 // 1%
			Cap:           sdk.NewCoin(assets.MicroSDRDenom, sdk.OneInt().MulRaw(assets.MicroUnit)), // 1 SDR Tax cap
			ChangeRateMax: sdk.NewDecWithPrec(25, 5),                                                // 0.025%
		},

		// Reward update policy
		PolicyConstraints{
			RateMin:       sdk.NewDecWithPrec(5, 2),             // 5%
			RateMax:       sdk.NewDecWithPrec(90, 2),            // 90%
			ChangeRateMax: sdk.NewDecWithPrec(25, 3),            // 2.5%
			Cap:           sdk.NewCoin("unused", sdk.ZeroInt()), // UNUSED
		},

		sdk.NewDecWithPrec(67, 2),  // 67%
		sdk.NewDecWithPrec(107, 2), // 1.07 mining increment; exponential growth

		sdk.NewInt(4),
		sdk.NewInt(52),
		sdk.NewInt(12),
	)
}

func validateParams(params Params) error {
	if params.TaxPolicy.RateMax.LT(params.TaxPolicy.RateMin) {
		return fmt.Errorf("treasury TaxPolicy.RateMax %s must be greater than TaxPolicy.RateMin %s",
			params.TaxPolicy.RateMax.String(), params.TaxPolicy.RateMin.String())
	}

	if params.TaxPolicy.RateMin.IsNegative() {
		return fmt.Errorf("treasury parameter TaxPolicy.RateMin must be >= 0, is %s", params.TaxPolicy.RateMin.String())
	}

	if params.RewardPolicy.RateMax.LT(params.RewardPolicy.RateMin) {
		return fmt.Errorf("treasury RewardPolicy.RateMax %s must be greater than RewardPolicy.RateMin %s",
			params.RewardPolicy.RateMax.String(), params.RewardPolicy.RateMin.String())
	}

	if params.RewardPolicy.RateMin.IsNegative() {
		return fmt.Errorf("treasury parameter RewardPolicy.RateMin must be >= 0, is %s", params.RewardPolicy.RateMin.String())
	}

	return nil
}

// implements fmt.Stringer
func (params Params) String() string {
	return fmt.Sprintf(`Treasury Params:
  Tax Policy        : { %v } 
  Reward Policy     : { %v }

  SeigniorageBurdenTarget : %v
  MiningIncrement   : %v

  WindowShort        : %v
  WindowLong         : %v
  `, params.TaxPolicy, params.RewardPolicy, params.SeigniorageBurdenTarget,
		params.MiningIncrement, params.WindowShort, params.WindowLong)
}
