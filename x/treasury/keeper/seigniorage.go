package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/treasury/types"
)

// SettleSeigniorage computes seigniorage and distributes it to oracle and distribution(community-pool) account
func (k Keeper) SettleSeigniorage(ctx sdk.Context) {
	// Mint seigniorage for oracle and community pool
	seigniorageLunaAmt := k.PeekEpochSeigniorage(ctx)
	if seigniorageLunaAmt.LTE(sdk.ZeroInt()) {
		return
	}

	// Settle current epoch seigniorage
	rewardWeight := k.GetRewardWeight(ctx)

	// Align seigniorage to usdr
	seigniorageDecCoin := sdk.NewDecCoin(core.MicroLunaDenom, seigniorageLunaAmt)

	// Mint seigniorage
	seigniorageCoin, _ := seigniorageDecCoin.TruncateDecimal()
	seigniorageCoins := sdk.NewCoins(seigniorageCoin)
	if seigniorageCoins.IsValid() {
		if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, seigniorageCoins); err != nil {
			panic(err)
		}
	}
	seigniorageAmt := seigniorageCoin.Amount

	// Send reward to oracle module
	burnAmt := rewardWeight.MulInt(seigniorageAmt).TruncateInt()
	burnCoins := sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, burnAmt))
	if burnCoins.IsValid() {
		if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, burnCoins); err != nil {
			panic(err)
		}
	}

	// Send left to distribution module
	leftAmt := seigniorageAmt.Sub(burnAmt)
	leftCoins := sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, leftAmt))
	if leftCoins.IsValid() {
		if err := k.bankKeeper.SendCoinsFromModuleToModule(
			ctx,
			types.ModuleName,
			k.distributionModuleName,
			leftCoins,
		); err != nil {
			panic(err)
		}

		// Update distribution community pool
		feePool := k.distrKeeper.GetFeePool(ctx)
		feePool.CommunityPool = feePool.CommunityPool.Add(sdk.NewDecCoinsFromCoins(leftCoins...)...)
		k.distrKeeper.SetFeePool(ctx, feePool)
	}
}
