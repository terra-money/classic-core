package keeper

import (
	"github.com/terra-money/core/x/treasury/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BurnCoinsFromBurnAccount burn all coins from burn account
func (k Keeper) BurnCoinsFromBurnAccount(ctx sdk.Context) {
	burnAddress := k.accountKeeper.GetModuleAddress(types.BurnModuleName)
	if coins := k.bankKeeper.GetAllBalances(ctx, burnAddress); !coins.IsZero() {
		err := k.bankKeeper.BurnCoins(ctx, types.BurnModuleName, coins)
		if err != nil {
			panic(err)
		}
	}

	return
}
