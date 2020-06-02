package bank

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-project/core/x/bank/internal/types"
)

// NewHookHandler trigger hook to burn the coins in special address
func NewHookHandler(k Keeper, sk types.SupplyKeeper, originHandler sdk.Handler) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		res, err := originHandler(ctx, msg)
		if err != nil {
			return res, err
		}

		burnModuleAcc := sk.GetModuleAccount(ctx, types.BurnModuleName)
		if coins := burnModuleAcc.GetCoins(); !coins.IsZero() {
			err = sk.BurnCoins(ctx, types.BurnModuleName, coins)
		}

		return res, err
	}
}
