package keeper

import (
	"github.com/terra-money/core/x/market/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BurnCoinsFromBurnAccount burn all coins from burn account
func (k Keeper) BurnCoinsFromBurnAccount(ctx sdk.Context) {
	burnAddress := k.AccountKeeper.GetModuleAddress(types.BurnModuleName)
	if coins := k.BankKeeper.GetAllBalances(ctx, burnAddress); !coins.IsZero() {
		err := k.BankKeeper.BurnCoins(ctx, types.BurnModuleName, coins)
		if err != nil {
			panic(err)
		}
	}
}
