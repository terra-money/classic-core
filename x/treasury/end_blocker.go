package treasury

import (
	"terra/types/assets"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker called to adjust macro weights (tax, mining reward) and settle outstanding claims.
func EndBlocker(ctx sdk.Context, k Keeper) {
	settlementPeriod := k.GetParams(ctx).SettlementPeriod
	if sdk.NewInt(ctx.BlockHeight()).Mod(settlementPeriod).Equal(sdk.ZeroInt()) {
		k.settleClaims(ctx)

		tax := updateTaxes(ctx, k)
		k.pk.SetTax(ctx, tax)

		rewardWeight := updateRewardWeight(ctx, k)
		k.SetRewardWeight(ctx, rewardWeight)
	}
}

// ProcessOracleRewardees process rewardees from the oracle module
func ProcessOracleRewardees(ctx sdk.Context, k Keeper, rewardees map[string]sdk.Int) {
	for voterAddress, rewardWeight := range rewardees {
		addr, err := sdk.AccAddressFromBech32(voterAddress)
		if err != nil {
			continue
		}

		oracleRewardClaim := NewClaim(ctx.BlockHeight(), OracleClaimClass, sdk.NewDecFromInt(rewardWeight), addr)
		k.AddClaim(ctx, oracleRewardClaim)
	}
}

func updateTaxes(ctx sdk.Context, k Keeper) sdk.Dec {
	params := k.GetParams(ctx)

	targetIssuance := k.GetParams(ctx).LunaTargetIssuance
	currentIssuance := k.GetIssuance(ctx, assets.LunaDenom)

	excessIssuance := currentIssuance.Sub(targetIssuance)
	if excessIssuance.Equal(sdk.ZeroInt()) {
		return params.TaxMin
	}

	return params.TaxMax.Sub(params.TaxMin).QuoInt(currentIssuance).MulInt(excessIssuance)
}

func updateRewardWeight(ctx sdk.Context, k Keeper) sdk.Dec {
	params := k.GetParams(ctx)

	targetIssuance := k.GetParams(ctx).LunaTargetIssuance
	currentIssuance := k.GetIssuance(ctx, assets.LunaDenom)

	excessIssuance := currentIssuance.Sub(targetIssuance)
	if excessIssuance.Equal(sdk.ZeroInt()) {
		return params.RewardMin
	}

	return params.RewardMax.Sub(params.RewardMin).QuoInt(currentIssuance).MulInt(excessIssuance)
}
