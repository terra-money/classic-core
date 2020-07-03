package msgauth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k Keeper) {

	// clears all the mature grants
	matureGrants := k.DequeueAllMatureGrantQueue(ctx)
	for _, grant := range matureGrants {
		k.RevokeGrant(ctx, grant.GranterAddress, grant.GranteeAddress, grant.MsgType)
	}
}
