package msgauth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/terra-project/core/x/msgauth/keeper"
	"github.com/terra-project/core/x/msgauth/types"
)

func TestGrant(t *testing.T) {
	input := keeper.CreateTestInput(t)
	h := NewHandler(input.AuthorizationKeeper)
	coins := sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(1_000_000_000)))

	// send authorization
	sendAuth := types.NewSendAuthorization(coins)
	msg, err := types.NewMsgGrantAuthorization(keeper.Addrs[0], keeper.Addrs[1], sendAuth, time.Hour)

	_, err = h(input.Ctx, msg)
	require.NoError(t, err)

	grant, found := input.AuthorizationKeeper.GetGrant(input.Ctx, keeper.Addrs[0], keeper.Addrs[1], sendAuth.MsgType())
	require.True(t, found)
	require.Equal(t, sendAuth, grant.GetAuthorization())
	require.Equal(t, input.Ctx.BlockTime().Add(time.Hour), grant.Expiration)

	// generic authorization
	genericAuth := types.NewGenericAuthorization("swap")
	msg, err = types.NewMsgGrantAuthorization(keeper.Addrs[0], keeper.Addrs[1], genericAuth, time.Hour)
	require.NoError(t, err)

	_, err = h(input.Ctx, msg)
	require.NoError(t, err)

	grant, found = input.AuthorizationKeeper.GetGrant(input.Ctx, keeper.Addrs[0], keeper.Addrs[1], "swap")
	require.True(t, found)
	require.Equal(t, genericAuth, grant.GetAuthorization())
	require.Equal(t, input.Ctx.BlockTime().Add(time.Hour), grant.Expiration)

	// test not allowed to grant
	genericAuth = types.NewGenericAuthorization("now allowed msg")
	msg, err = types.NewMsgGrantAuthorization(keeper.Addrs[0], keeper.Addrs[1], genericAuth, time.Hour)
	require.NoError(t, err)

	_, err = h(input.Ctx, msg)
	require.Error(t, err)
}

func TestRevoke(t *testing.T) {
	input := keeper.CreateTestInput(t)
	h := NewHandler(input.AuthorizationKeeper)
	coins := sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(1_000_000_000)))
	grantMsg, err := types.NewMsgGrantAuthorization(keeper.Addrs[0], keeper.Addrs[1], types.NewSendAuthorization(coins), time.Hour)

	_, err = h(input.Ctx, grantMsg)
	require.NoError(t, err)

	revokeMsg := types.NewMsgRevokeAuthorization(keeper.Addrs[0], keeper.Addrs[1], banktypes.TypeMsgSend)
	_, err = h(input.Ctx, revokeMsg)
	require.NoError(t, err)

	_, found := input.AuthorizationKeeper.GetGrant(input.Ctx, keeper.Addrs[1], keeper.Addrs[0], banktypes.TypeMsgSend)
	require.False(t, found)
}

func TestExecute(t *testing.T) {
	input := keeper.CreateTestInput(t)
	h := NewHandler(input.AuthorizationKeeper)
	coins := sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(1_000_000_000)))
	input.BankKeeper.SetBalances(input.Ctx, keeper.Addrs[0], coins)

	grantMsg, err := types.NewMsgGrantAuthorization(keeper.Addrs[0], keeper.Addrs[1], types.NewSendAuthorization(coins), time.Hour)
	require.NoError(t, err)

	_, err = h(input.Ctx, grantMsg)
	require.NoError(t, err)

	execMsg, err := types.NewMsgExecAuthorized(keeper.Addrs[1], []sdk.Msg{
		banktypes.NewMsgSend(keeper.Addrs[0], keeper.Addrs[1], coins),
	})

	_, err = h(input.Ctx, execMsg)
	require.NoError(t, err)
}
