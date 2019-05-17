package treasury

import (
	"fmt"

	"github.com/terra-project/core/types"
	"github.com/terra-project/core/types/util"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all treasury state that must be provided at genesis
type GenesisState struct {
	Params              Params        `json:"params"` // treasury params
	GenesisTaxRate      sdk.Dec       `json:"tax_rate"`
	GenesisRewardWeight sdk.Dec       `json:"reward_weight"`
	Claims              []types.Claim `json:"claims"`
}

// NewGenesisState constructs a new genesis state
func NewGenesisState(params Params, taxRate, rewardWeight sdk.Dec, claims []types.Claim) GenesisState {
	return GenesisState{
		Params:              params,
		GenesisTaxRate:      taxRate,
		GenesisRewardWeight: rewardWeight,
		Claims:              claims,
	}
}

// DefaultGenesisState returns raw genesis message for testing
func DefaultGenesisState() GenesisState {
	params := DefaultParams()
	return GenesisState{
		Params:              params,
		GenesisTaxRate:      sdk.NewDecWithPrec(1, 3), // 0.1%
		GenesisRewardWeight: sdk.NewDecWithPrec(5, 2), // 5%
		Claims:              []types.Claim{},
	}
}

// InitGenesis new treasury genesis
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetParams(ctx, data.Params)
	keeper.SetTaxRate(ctx, data.GenesisTaxRate)
	keeper.setTaxCap(ctx, data.Params.TaxPolicy.Cap.Denom, data.Params.TaxPolicy.Cap.Amount)
	keeper.SetRewardWeight(ctx, data.GenesisRewardWeight)

	for _, claim := range data.Claims {
		keeper.AddClaim(ctx, claim)
	}

}

// ExportGenesis returns a GenesisState for a given context and keeper. The
// GenesisState will contain the pool, and validator/delegator distribution info's
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	params := k.GetParams(ctx)
	taxRate := k.GetTaxRate(ctx, sdk.ZeroInt())
	rewardWeight := k.GetRewardWeight(ctx, util.GetEpoch(ctx))

	var claims []types.Claim
	k.IterateClaims(ctx, func(claim types.Claim) (stop bool) {
		claims = append(claims, claim)

		return false
	})

	return NewGenesisState(params, taxRate, rewardWeight, claims)
}

// ValidateGenesis validates the provided treasury genesis state to ensure the
// expected invariants holds. (i.e. params in correct bounds, no duplicate validators)
func ValidateGenesis(data GenesisState) error {
	if data.GenesisTaxRate.GT(data.Params.TaxPolicy.RateMax) ||
		data.GenesisTaxRate.LT(data.Params.TaxPolicy.RateMin) {
		return fmt.Errorf("Genesis tax rate must be between %s and %s, is %s",
			data.Params.TaxPolicy.RateMin, data.Params.TaxPolicy.RateMax, data.GenesisTaxRate)
	}

	if data.GenesisRewardWeight.GT(data.Params.RewardPolicy.RateMax) ||
		data.GenesisRewardWeight.LT(data.Params.RewardPolicy.RateMin) {
		return fmt.Errorf("Genesis reward rate must be between %s and %s, is %s",
			data.Params.RewardPolicy.RateMin, data.Params.RewardPolicy.RateMax, data.GenesisRewardWeight)
	}

	return validateParams(data.Params)
}
