package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-project/core/x/treasury/types"
)

// BurnCoinsFromBurnAccount burn all coins from burn account
func (k Keeper) BurnCoinsFromBurnAccount(ctx sdk.Context) {
	burnAddress := k.accountKeeper.GetModuleAddress(types.BurnModuleName)
	if coins := k.bankKeeper.GetAllBalances(ctx, burnAddress); !coins.IsZero() {
		fmt.Println(coins)
		err := k.bankKeeper.BurnCoins(ctx, types.BurnModuleName, coins)
		if err != nil {
			panic(err)
		}
	}

	return
}
