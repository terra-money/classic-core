package treasury

import (
	"fmt"
	"terra/types/assets"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Params treasury parameters
type Params struct {
	TaxRateMin sdk.Dec `json:"tax_rate_min"` // percentage cap on taxes. Defaults to 2%.
	TaxRateMax sdk.Dec `json:"tax_rate_max"` // percentage floor on taxes. Defaults to 0.

	TaxCap sdk.Coin `json:"tax_cap"` // Tax Cap in TerraSDR

	RewardMin sdk.Dec `json:"reward_min"` // percentage floor on miner rewards for seigniorage. Defaults to 0.1.
	RewardMax sdk.Dec `json:"reward_max"` // percentage cap on miner rewards for seigniorage. Defaults to 0.9

	EpochLong  sdk.Int `json:"epoch_long"`
	EpochShort sdk.Int `json:"epoch_short"`

	OracleClaimShare sdk.Dec `json:"oracle_share"`
	BudgetClaimShare sdk.Dec `json:"budget_share"`
}

// NewParams creates a new param instance
func NewParams(taxRateMin, taxRateMax, rewardMin, rewardMax, oracleClaimShare, budgetClaimShare sdk.Dec,
	epochLong, epochShort sdk.Int, taxCap sdk.Coin) Params {
	return Params{
		TaxRateMin:       taxRateMin,
		TaxRateMax:       taxRateMax,
		TaxCap:           taxCap,
		RewardMin:        rewardMin,
		RewardMax:        rewardMax,
		OracleClaimShare: oracleClaimShare,
		BudgetClaimShare: budgetClaimShare,
		EpochLong:        epochLong,
		EpochShort:       epochShort,
	}
}

// DefaultParams creates default treasury module parameters
func DefaultParams() Params {
	return NewParams(
		sdk.NewDecWithPrec(1, 3),                   // 0.1%
		sdk.NewDecWithPrec(2, 2),                   // 2%
		sdk.NewDecWithPrec(5, 2),                   // 5%
		sdk.NewDecWithPrec(9, 1),                   // 90%
		sdk.NewDecWithPrec(1, 1),                   // 10%
		sdk.NewDecWithPrec(9, 1),                   // 90%
		sdk.NewInt(52),                             // Approx. 1 year
		sdk.NewInt(4),                              // Approx. 1 month
		sdk.NewCoin(assets.SDRDenom, sdk.OneInt()), // 1 TerraSDR as cap
	)
}

func validateParams(params Params) error {
	if params.TaxRateMax.LT(params.TaxRateMin) {
		return fmt.Errorf("treasury parameter TaxRateMax (%s) must be greater than TaxRateMin (%s)",
			params.TaxRateMax.String(), params.TaxRateMin.String())
	}

	if params.TaxRateMin.IsNegative() {
		return fmt.Errorf("treasury parameter TaxRateMin must be >= 0, is %s", params.TaxRateMin.String())
	}

	if params.RewardMax.LT(params.RewardMin) {
		return fmt.Errorf("treasury parameter RewardMax (%s) must be greater than RewardMin (%s)",
			params.RewardMax.String(), params.RewardMin.String())
	}

	if params.RewardMin.IsNegative() {
		return fmt.Errorf("treasury parameter RewardMin must be >= 0, is %s", params.RewardMin.String())
	}

	shareSum := params.OracleClaimShare.Add(params.BudgetClaimShare)
	if !shareSum.Equal(sdk.OneDec()) {
		return fmt.Errorf("treasury parameter ClaimShares must sum to 1, but sums to %s", shareSum.String())
	}

	return nil
}

func (params Params) String() string {
	return fmt.Sprintf(`Treasury Params:
  Tax Rate Min: %s
  Tax Rate Max: %s
 
  Tax Cap: %s

  Mining Reward Weight Min: %v
  Mining Reward Weight Max: %v

  Oracle Reward Weight: %v
  Budget Reward Weight: %v 

  Epoch Long: %v 
  Epoch Short %v
  `, params.TaxRateMin, params.TaxRateMax, params.TaxCap,
		params.RewardMin, params.RewardMax, params.EpochLong,
		params.EpochShort)
}
