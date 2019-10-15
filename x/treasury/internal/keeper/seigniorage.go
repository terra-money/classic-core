package keeper

import (
	"fmt"

	"github.com/terra-project/core/x/treasury/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	core "github.com/terra-project/core/types"
)

// SettleSeigniorage
func (k Keeper) SettleSeigniorage(ctx sdk.Context) {
	// Mint seigniorage for oracle and community pool
	epoch := core.GetEpoch(ctx)
	seigniorageLunaAmt := k.PeekEpochSeigniorage(ctx, epoch)
	if seigniorageLunaAmt.LTE(sdk.ZeroInt()) {
		return
	}

	// Settle current epoch seigniorage
	rewardWeight := k.GetRewardWeight(ctx, epoch)

	// Align seigniorage to usdr
	seigniorageLunaDecCoin := sdk.NewDecCoin(core.MicroLunaDenom, seigniorageLunaAmt)
	seigniorageDecCoin, err := k.marketKeeper.ComputeInternalSwap(ctx, seigniorageLunaDecCoin, core.MicroSDRDenom)
	if err != nil {
		k.Logger(ctx).Error(fmt.Sprintf("[Treasury] Failed to swap seigniorage to usdr, %s", err.Error()))
		return
	}

	// Mint seigniorage
	seigniorageCoin, _ := seigniorageDecCoin.TruncateDecimal()
	seigniorageCoins := sdk.NewCoins(seigniorageCoin)
	err = k.supplyKeeper.MintCoins(ctx, types.ModuleName, seigniorageCoins)
	if err != nil {
		panic(err)
	}
	seigniorageAmt := seigniorageCoin.Amount

	// Send reward to oracle module
	oracleRewardAmt := rewardWeight.MulInt(seigniorageAmt).TruncateInt()
	oracleRewardCoins := sdk.NewCoins(sdk.NewCoin(core.MicroSDRDenom, oracleRewardAmt))
	err = k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.oracleModuleName, oracleRewardCoins)
	if err != nil {
		panic(err)
	}

	// Send left to distribution module
	leftAmt := seigniorageAmt.Sub(oracleRewardAmt)
	leftCoins := sdk.NewCoins(sdk.NewCoin(core.MicroSDRDenom, leftAmt))
	err = k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.distributionModuleName, leftCoins)
	if err != nil {
		panic(err)
	}

	// Update distribution community p9ol
	feePool := k.distrKeeper.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(sdk.NewDecCoins(leftCoins))
	k.distrKeeper.SetFeePool(ctx, feePool)
}
