package treasury

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Params treasury parameters
type Params struct {
	TaxMin sdk.Dec `json:"tax_min"` // percentage cap on taxes. Defaults to 2%.
	TaxMax sdk.Dec `json:"tax_max"` // percentage floor on taxes. Defaults to 0.

	RewardMin sdk.Dec `json:"reward_min"` // percentage floor on miner rewards for seigniorage. Defaults to 0.1.
	RewardMax sdk.Dec `json:"reward_max"` // percentage cap on miner rewards for seigniorage. Defaults to 0.9

	LunaTargetIssuance sdk.Int `json:"luna_target"`
	SettlementPeriod   sdk.Int `json:"settlement_period"`
}

// NewParams creates a new param instance
func NewParams(taxMin, taxMax, rewardMin, rewardMax sdk.Dec, lunaTargetIssuance, settlementPeriod sdk.Int) Params {
	return Params{
		TaxMin:             taxMin,
		TaxMax:             taxMax,
		RewardMin:          rewardMin,
		RewardMax:          rewardMax,
		LunaTargetIssuance: lunaTargetIssuance,
		SettlementPeriod:   settlementPeriod,
	}
}

// DefaultParams creates default treasury module parameters
func DefaultParams() Params {
	return NewParams(
		sdk.NewDecWithPrec(1, 3), // 0.1%
		sdk.NewDecWithPrec(2, 2), // 2%
		sdk.NewDecWithPrec(1, 1), // 10%
		sdk.NewDecWithPrec(9, 1), // 90%
		sdk.NewInt(int64(10^9)),
		sdk.NewInt(3000000), // Approx. 1 month
	)
}

func validateParams(params Params) error {
	return nil
}
