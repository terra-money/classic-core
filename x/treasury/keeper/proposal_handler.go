package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-project/core/x/treasury/types"
)

// HandleTaxRateUpdateProposal is a handler for updating tax rate
func HandleTaxRateUpdateProposal(ctx sdk.Context, k Keeper, p *types.TaxRateUpdateProposal) error {
	taxPolicy := k.TaxPolicy(ctx)
	taxRate := k.GetTaxRate(ctx)
	newTaxRate := taxPolicy.Clamp(taxRate, p.TaxRate)

	// Set the new tax rate to the store
	k.SetTaxRate(ctx, newTaxRate)

	// Emit gov handler events
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeTaxRateUpdate,
			sdk.NewAttribute(types.AttributeKeyTaxRate, newTaxRate.String()),
		),
	)

	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("updated tax-rate to %s", newTaxRate))
	return nil
}

// HandleRewardWeightUpdateProposal is a handler for updating reward weight
func HandleRewardWeightUpdateProposal(ctx sdk.Context, k Keeper, p *types.RewardWeightUpdateProposal) error {
	rewardPolicy := k.RewardPolicy(ctx)
	rewardWeight := k.GetRewardWeight(ctx)
	newRewardWeight := rewardPolicy.Clamp(rewardWeight, p.RewardWeight)

	// Set the new reward rate to the store
	k.SetRewardWeight(ctx, newRewardWeight)

	// Emit gov handler events
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeRewardWeightUpdate,
			sdk.NewAttribute(types.AttributeKeyRewardWeight, newRewardWeight.String()),
		),
	)

	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("updated reward-weight to %s", newRewardWeight))
	return nil
}
