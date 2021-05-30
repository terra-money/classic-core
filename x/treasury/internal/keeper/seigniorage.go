package keeper

import (
	"github.com/terra-money/core/x/treasury/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	core "github.com/terra-money/core/types"
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
	err := k.supplyKeeper.MintCoins(ctx, types.ModuleName, seigniorageCoins)
	if err != nil {
		panic(err)
	}
	seigniorageAmt := seigniorageCoin.Amount

	// Send reward to oracle module
	oracleRewardAmt := rewardWeight.MulInt(seigniorageAmt).TruncateInt()
	oracleRewardCoins := sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, oracleRewardAmt))
	err = k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.oracleModuleName, oracleRewardCoins)
	if err != nil {
		panic(err)
	}

	// Send left to distribution module
	leftAmt := seigniorageAmt.Sub(oracleRewardAmt)
	leftCoins := sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, leftAmt))
	err = k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.distributionModuleName, leftCoins)
	if err != nil {
		panic(err)
	}

	// Update distribution community pool
	feePool := k.distrKeeper.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(sdk.NewDecCoinsFromCoins(leftCoins...)...)
	k.distrKeeper.SetFeePool(ctx, feePool)
}
