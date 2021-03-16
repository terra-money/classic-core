package msgauth

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-project/core/x/msgauth/keeper"
	"github.com/terra-project/core/x/msgauth/types"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	// clears all the mature grants
	matureGrants := k.DequeueAllMatureGrantQueue(ctx)
	for _, grant := range matureGrants.Pairs {
		granter, err := sdk.AccAddressFromBech32(grant.GranterAddress)
		if err != nil {
			panic(fmt.Sprintf("invalid granter address %s not found", grant.GranterAddress))
		}

		grantee, err := sdk.AccAddressFromBech32(grant.GranteeAddress)
		if err != nil {
			panic(fmt.Sprintf("invalid grantee address %s not found", grant.GranteeAddress))
		}

		k.RevokeGrant(ctx, granter, grantee, grant.MsgType)
	}
}
