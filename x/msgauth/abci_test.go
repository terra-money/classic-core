package msgauth

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/terra-project/core/x/msgauth/keeper"
	"github.com/terra-project/core/x/msgauth/types"
)

func TestMature(t *testing.T) {
	input := keeper.CreateTestInput(t)
	h := NewHandler(input.AuthorizationKeeper)
	coins := sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(1_000_000_000)))

	// send authorization
	sendAuth := types.NewSendAuthorization(coins)
	msg, err := types.NewMsgGrantAuthorization(keeper.Addrs[0], keeper.Addrs[1], sendAuth, time.Hour)
	require.NoError(t, err)

	_, err = h(input.Ctx, msg)
	require.NoError(t, err)

	grant, found := input.AuthorizationKeeper.GetGrant(input.Ctx, keeper.Addrs[0], keeper.Addrs[1], sendAuth.MsgType())
	require.True(t, found)
	require.Equal(t, sendAuth, grant.GetAuthorization())
	require.Equal(t, input.Ctx.BlockTime().Add(time.Hour), grant.Expiration)

	EndBlocker(input.Ctx.WithBlockTime(input.Ctx.BlockTime().Add(time.Hour)), input.AuthorizationKeeper)
	_, found = input.AuthorizationKeeper.GetGrant(input.Ctx, keeper.Addrs[0], keeper.Addrs[1], sendAuth.MsgType())
	require.False(t, found)
}
