package treasury

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/terra-project/core/x/treasury/internal/types"
)

// NewTreasuryPolicyUpdateHandler custom gov proposal handler
func NewTreasuryPolicyUpdateHandler(k Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) sdk.Error {
		switch c := content.(type) {
		case TaxRateUpdateProposal:
			return handleTaxRateUpdateProposal(ctx, k, c)
		case RewardWeightUpdateProposal:
			return handleRewardWeightUpdateProposal(ctx, k, c)

		default:
			errMsg := fmt.Sprintf("unrecognized distr proposal content type: %T", c)
			return sdk.ErrUnknownRequest(errMsg)
		}
	}
}

// handleTaxRateUpdateProposal is a handler for updating tax rate
func handleTaxRateUpdateProposal(ctx sdk.Context, k Keeper, p TaxRateUpdateProposal) sdk.Error {
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

// handleRewardWeightUpdateProposal is a handler for updating reward weight
func handleRewardWeightUpdateProposal(ctx sdk.Context, k Keeper, p RewardWeightUpdateProposal) sdk.Error {
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
