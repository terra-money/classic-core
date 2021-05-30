package bank

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/core/x/bank/internal/types"
)

func EndBlocker(ctx sdk.Context, k Keeper, sk types.SupplyKeeper) {
	burnModuleAcc := sk.GetModuleAccount(ctx, types.BurnModuleName)
	if coins := burnModuleAcc.GetCoins(); !coins.IsZero() {
		// ignore error; error never happens
		_ = sk.BurnCoins(ctx, types.BurnModuleName, coins)
	}
}
