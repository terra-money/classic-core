package treasury

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	settlementPeriod = 100000
)

// func calibrateTax(k Keeper) {
// 	k.tk.Set
// }

func EndBlocker(ctx sdk.Context, k Keeper) (resTags sdk.Tags) {
	tags := sdk.NewTags()

	if ctx.BlockHeight()%settlementPeriod == int64(0) {
		k.SettleShares(ctx)
	}

	return tags
}
