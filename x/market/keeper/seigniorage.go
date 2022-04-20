package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/market/types"
)

// SettleSeigniorage settle seigniorage to the registered addresses
// and burn left coins. The recipient addresses can be registered
// via SeigniorageRouteChangeProposal.
func (k Keeper) SettleSeigniorage(ctx sdk.Context) {
	moduleAddr := k.AccountKeeper.GetModuleAddress(types.ModuleName)
	collectedCoins := k.BankKeeper.GetAllBalances(ctx, moduleAddr)

	// no coins, then no actions are required
	if collectedCoins.Empty() {
		return
	}

	// only Luna will be distributed as seigniorage
	seigniorageAmount := collectedCoins.AmountOf(core.MicroLunaDenom)
	routes := k.GetSeigniorageRoutes(ctx)

	var burnCoins sdk.Coins

	// If seigniorageAmount is zero, then just burn all collected coins
	if !seigniorageAmount.IsZero() {
		leftSeigniorageAmount := sdk.NewInt(seigniorageAmount.Int64())
		for _, route := range routes {
			routeAmount := route.Weight.MulInt(seigniorageAmount).TruncateInt()
			if routeAmount.IsPositive() {
				coins := sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, routeAmount))
				recipient, err := sdk.AccAddressFromBech32(route.Address)
				if err != nil {
					panic(err)
				}

				// transfer weight * seigniorage amount LUNA token to the recipient address
				if route.Address == types.AlternateCommunityPoolAddress.String() {
					// If the given address is the predefined alternate address,
					// fund community pool because community pool does not have
					// its own address,
					// - https://github.com/cosmos/cosmos-sdk/issues/10811
					err = k.DistributionKeeper.FundCommunityPool(ctx, coins, moduleAddr)
				} else {
					err = k.BankKeeper.SendCoins(ctx, moduleAddr, recipient, coins)
				}
				if err != nil {
					panic(err)
				}

				leftSeigniorageAmount = leftSeigniorageAmount.Sub(routeAmount)
			}
		}

		for _, coin := range collectedCoins {
			// replace Luna amount to burn amount
			if coin.Denom == core.MicroLunaDenom {
				coin.Amount = leftSeigniorageAmount
			}

			if coin.Amount.IsPositive() {
				burnCoins = append(burnCoins, coin)
			}
		}
	} else {
		burnCoins = collectedCoins
	}

	// burn all left coins
	if !burnCoins.Empty() {
		err := k.BankKeeper.BurnCoins(ctx, types.ModuleName, burnCoins)
		if err != nil {
			panic(err)
		}
	}

	return
}
