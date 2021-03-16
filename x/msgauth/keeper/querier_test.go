package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/terra-project/core/x/msgauth/types"
)

func TestQueryGrants(t *testing.T) {
	input := CreateTestInput(t)
	ctx := sdk.WrapSDKContext(input.Ctx)

	querier := NewQuerier(input.AuthorizationKeeper)

	now := input.Ctx.BlockHeader().Time

	// register grant
	grant, err := types.NewAuthorizationGrant(types.NewSendAuthorization(sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(123)))), now)
	require.NoError(t, err)

	input.AuthorizationKeeper.SetGrant(input.Ctx, Addrs[0], Addrs[1], banktypes.TypeMsgSend, grant)

	// query grants
	res, err := querier.Grants(ctx, &types.QueryGrantsRequest{
		Granter: Addrs[0].String(),
		Grantee: Addrs[1].String(),
	})
	require.NoError(t, err)
	require.Equal(t, &types.QueryGrantsResponse{Grants: types.AuthorizationGrants{grant}}, res)
}

func TestQueryAllGrants(t *testing.T) {
	input := CreateTestInput(t)
	ctx := sdk.WrapSDKContext(input.Ctx)

	querier := NewQuerier(input.AuthorizationKeeper)
	now := input.Ctx.BlockHeader().Time

	// register grants
	grant, err := types.NewAuthorizationGrant(types.NewSendAuthorization(sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(123)))), now)
	require.NoError(t, err)

	grant2, err := types.NewAuthorizationGrant(types.NewGenericAuthorization(banktypes.TypeMsgSend+"2"), now)
	require.NoError(t, err)

	input.AuthorizationKeeper.SetGrant(input.Ctx, Addrs[0], Addrs[1], banktypes.TypeMsgSend, grant)
	input.AuthorizationKeeper.SetGrant(input.Ctx, Addrs[0], Addrs[1], banktypes.TypeMsgSend+"2", grant2)

	res, err := querier.Grants(ctx, &types.QueryGrantsRequest{
		Granter: Addrs[0].String(),
		Grantee: Addrs[1].String(),
	})
	require.NoError(t, err)
	require.Equal(t, &types.QueryGrantsResponse{Grants: types.AuthorizationGrants{grant, grant2}}, res)
}
