package market

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/market/internal/types"
)

// market end block functionality
func EndBlocker(ctx sdk.Context, k Keeper) {
	if !core.IsPeriodLastBlock(ctx, core.BlocksPerDay) {
		return
	}

	// update luna issuance at last block of a day
	updatedIssuance := k.UpdateLastDayIssuance(ctx)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventDaliyIssuanceUpdate,
			sdk.NewAttribute(types.AttributeKeyIssuance, updatedIssuance.String()),
		),
	)
}
