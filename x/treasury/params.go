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

	EpochShort     sdk.Int `json:"epoch_short"`
	EpochLong      sdk.Int `json:"epoch_long"`
	EpochProbation sdk.Int `json:"epoch_probation"`

	OracleClaimShare sdk.Dec `json:"oracle_share"`
	BudgetClaimShare sdk.Dec `json:"budget_share"`
}

// NewParams creates a new param instance
func NewParams(
	taxPolicy, rewardPolicy PolicyConstraints,
	seigniorageBurden sdk.Dec,
	miningIncrement sdk.Dec,
	epochShort, epochLong, epochProbation sdk.Int,
	oracleShare, budgetShare sdk.Dec,
) Params {
	return Params{
		TaxPolicy:               taxPolicy,
		RewardPolicy:            rewardPolicy,
		SeigniorageBurdenTarget: seigniorageBurden,
		MiningIncrement:         miningIncrement,
		EpochShort:              epochShort,
		EpochLong:               epochLong,
		EpochProbation:          epochProbation,
		OracleClaimShare:        oracleShare,
		BudgetClaimShare:        budgetShare,
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
			RateMax:       sdk.NewDecWithPrec(20, 2),            // 20%
			ChangeRateMax: sdk.NewDecWithPrec(25, 3),            // 2.5%
			Cap:           sdk.NewCoin("unused", sdk.ZeroInt()), // UNUSED
		},

		sdk.NewDecWithPrec(67, 2),  // 67%
		sdk.NewDecWithPrec(107, 2), // 1.07 mining increment; exponential growth

		sdk.NewInt(4),
		sdk.NewInt(52),
		sdk.NewInt(12),

		sdk.NewDecWithPrec(1, 1), // 10%
		sdk.NewDecWithPrec(9, 1), // 90%
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

	shareSum := params.OracleClaimShare.Add(params.BudgetClaimShare)
	if !shareSum.Equal(sdk.OneDec()) {
		return fmt.Errorf("treasury parameter ClaimShares must sum to 1, but sums to %s", shareSum.String())
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

  EpochShort        : %v
  EpochLong         : %v

  OracleClaimShare  : %v
  BudgetClaimShare  : %v
  `, params.TaxPolicy, params.RewardPolicy, params.SeigniorageBurdenTarget,
		params.MiningIncrement, params.EpochShort, params.EpochLong,
		params.OracleClaimShare, params.BudgetClaimShare)
}
