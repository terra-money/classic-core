package keeper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/msgauth/types"
)

func TestKeeper(t *testing.T) {
	input := CreateTestInput(t)

	// "verify that no authorization returns nil"
	_, found := input.AuthorizationKeeper.GetGrant(input.Ctx, Addrs[0], Addrs[1], banktypes.TypeMsgSend)
	require.False(t, found)

	now := input.Ctx.BlockHeader().Time

	// "verify if authorization is accepted"
	grant, err := types.NewAuthorizationGrant(types.NewSendAuthorization(InitCoins), now.Add(time.Hour))
	require.NoError(t, err)

	input.AuthorizationKeeper.SetGrant(input.Ctx, Addrs[0], Addrs[1], banktypes.TypeMsgSend, grant)
	grant2, _ := input.AuthorizationKeeper.GetGrant(input.Ctx, Addrs[0], Addrs[1], banktypes.TypeMsgSend)
	require.Equal(t, grant, grant2)
	require.Equal(t, grant2.GetAuthorization().MsgType(), banktypes.TypeMsgSend)

	// "verify fetching authorization with wrong msg type fails"
	_, found = input.AuthorizationKeeper.GetGrant(input.Ctx, Addrs[0], Addrs[1], banktypes.MsgMultiSend{}.Type())
	require.False(t, found)

	// "verify fetching authorization with wrong grantee fails"
	_, found = input.AuthorizationKeeper.GetGrant(input.Ctx, Addrs[0], Addrs[2], banktypes.TypeMsgSend)
	require.False(t, found)

	grants := input.AuthorizationKeeper.GetGrants(input.Ctx, Addrs[0], Addrs[1])
	require.Equal(t, 1, len(grants))

	allGrants := input.AuthorizationKeeper.GetAllGrants(input.Ctx, Addrs[0])
	require.Equal(t, 1, len(allGrants))

	input.AuthorizationKeeper.IterateGrants(input.Ctx, func(
		granter, grantee sdk.AccAddress, grant2 types.AuthorizationGrant,
	) bool {
		require.Equal(t, Addrs[0], granter)
		require.Equal(t, Addrs[1], grantee)
		require.Equal(t, grant, grant2)
		return false
	})

	// "verify revoke fails with wrong information"
	input.AuthorizationKeeper.RevokeGrant(input.Ctx, Addrs[0], Addrs[2], banktypes.TypeMsgSend)
	_, found = input.AuthorizationKeeper.GetGrant(input.Ctx, Addrs[0], Addrs[1], banktypes.TypeMsgSend)
	require.True(t, found)

	// "verify revoke executes with correct information"
	input.AuthorizationKeeper.RevokeGrant(input.Ctx, Addrs[0], Addrs[1], banktypes.TypeMsgSend)
	_, found = input.AuthorizationKeeper.GetGrant(input.Ctx, Addrs[0], Addrs[1], banktypes.TypeMsgSend)
	require.False(t, found)

}

func TestKeeperFees(t *testing.T) {
	input := CreateTestInput(t)

	now := input.Ctx.BlockHeader().Time

	smallCoins := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 2))
	coins := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 20))
	largeCoins := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 123))

	msgs := []sdk.Msg{banktypes.NewMsgSend(Addrs[0], Addrs[2], smallCoins)}

	// "verify dispatch fails with invalid authorization"
	err := input.AuthorizationKeeper.DispatchActions(input.Ctx, Addrs[1], msgs)
	require.Error(t, err)

	// "verify dispatch executes with correct information"
	// grant authorization
	grant, err := types.NewAuthorizationGrant(types.NewSendAuthorization(coins), now)
	require.NoError(t, err)
	input.AuthorizationKeeper.SetGrant(input.Ctx, Addrs[0], Addrs[1], banktypes.TypeMsgSend, grant)

	grant2, found := input.AuthorizationKeeper.GetGrant(input.Ctx, Addrs[0], Addrs[1], banktypes.TypeMsgSend)
	require.True(t, found)
	require.Equal(t, grant, grant2)
	require.Equal(t,
		grant.Authorization.GetCachedValue().(types.AuthorizationI).MsgType(),
		banktypes.TypeMsgSend,
	)

	err = input.AuthorizationKeeper.DispatchActions(input.Ctx, Addrs[1], msgs)
	require.NoError(t, err)

	_, found = input.AuthorizationKeeper.GetGrant(input.Ctx, Addrs[0], Addrs[1], banktypes.TypeMsgSend)
	require.True(t, found)

	// "verify dispatch fails with overlimit"
	msgs = []sdk.Msg{banktypes.NewMsgSend(Addrs[0], Addrs[2], largeCoins)}
	err = input.AuthorizationKeeper.DispatchActions(input.Ctx, Addrs[1], msgs)
	require.Error(t, err)

	_, found = input.AuthorizationKeeper.GetGrant(input.Ctx, Addrs[0], Addrs[1], banktypes.TypeMsgSend)
	require.True(t, found)

	// "verify dispatch success and revoke grant which is out of limit"
	msgs = []sdk.Msg{banktypes.NewMsgSend(Addrs[0], Addrs[2], coins.Sub(smallCoins))}
	err = input.AuthorizationKeeper.DispatchActions(input.Ctx, Addrs[1], msgs)
	require.NoError(t, err)

	_, found = input.AuthorizationKeeper.GetGrant(input.Ctx, Addrs[0], Addrs[1], banktypes.TypeMsgSend)
	require.False(t, found)
}

func TestGrantQueue(t *testing.T) {
	input := CreateTestInput(t)

	now := input.Ctx.BlockTime()
	input.AuthorizationKeeper.InsertGrantQueue(input.Ctx, Addrs[0], Addrs[1], banktypes.TypeMsgSend, now.Add(time.Hour))
	input.AuthorizationKeeper.InsertGrantQueue(input.Ctx, Addrs[0], Addrs[1], banktypes.TypeMsgSend+"2", now.Add(time.Hour))

	ggmPairs := []types.GGMPair{
		{
			GranterAddress: Addrs[0].String(),
			GranteeAddress: Addrs[1].String(),
			MsgType:        banktypes.TypeMsgSend,
		},
		{
			GranterAddress: Addrs[0].String(),
			GranteeAddress: Addrs[1].String(),
			MsgType:        banktypes.TypeMsgSend + "2",
		},
	}

	timeSlice := input.AuthorizationKeeper.GetGrantQueueTimeSlice(input.Ctx, now)
	require.Equal(t, 0, len(timeSlice.Pairs))

	timeSlice = input.AuthorizationKeeper.GetGrantQueueTimeSlice(input.Ctx, now.Add(time.Hour))
	require.Equal(t, ggmPairs, timeSlice.Pairs)

	allPairs := input.AuthorizationKeeper.DequeueAllMatureGrantQueue(input.Ctx.WithBlockTime(now))
	require.Equal(t, 0, len(allPairs.Pairs))

	allPairs = input.AuthorizationKeeper.DequeueAllMatureGrantQueue(input.Ctx.WithBlockTime(now.Add(time.Hour)))
	require.Equal(t, ggmPairs, allPairs.Pairs)

	input.AuthorizationKeeper.InsertGrantQueue(input.Ctx, Addrs[0], Addrs[1], banktypes.TypeMsgSend, now.Add(time.Hour))
	input.AuthorizationKeeper.RevokeFromGrantQueue(input.Ctx, Addrs[0], Addrs[1], banktypes.TypeMsgSend, now.Add(time.Hour))
	timeSlice = input.AuthorizationKeeper.GetGrantQueueTimeSlice(input.Ctx, now.Add(time.Hour))
	require.Equal(t, 0, len(timeSlice.Pairs))
}
