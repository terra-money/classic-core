package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/terra-project/core/x/msgauth/types"
)

func TestLegacyNewLegacyQuerier(t *testing.T) {
	input := CreateTestInput(t)
	querier := NewLegacyQuerier(input.AuthorizationKeeper, input.Cdc)

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	_, err := querier(input.Ctx, []string{"INVALID_PATH"}, query)
	require.Error(t, err)
}

func TestLegacyQueryGrants(t *testing.T) {
	input := CreateTestInput(t)
	querier := NewLegacyQuerier(input.AuthorizationKeeper, input.Cdc)
	now := input.Ctx.BlockHeader().Time

	// register grant
	grant, err := types.NewAuthorizationGrant(types.NewSendAuthorization(sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(123)))), now)
	require.NoError(t, err)

	input.AuthorizationKeeper.SetGrant(input.Ctx, Addrs[0], Addrs[1], banktypes.TypeMsgSend, grant)

	params := types.NewQueryGrantsParams(Addrs[0], Addrs[1])
	bz, err := input.Cdc.MarshalJSON(params)
	require.NoError(t, err)

	query := abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, err := querier(input.Ctx, []string{types.QueryGrants}, query)
	require.NoError(t, err)

	var resGrant types.AuthorizationGrants
	input.Cdc.MustUnmarshalJSON(res, &resGrant)
	require.Equal(t, types.AuthorizationGrants{grant}, resGrant)
}

func TestLegacyQueryAllGrants(t *testing.T) {
	input := CreateTestInput(t)
	querier := NewLegacyQuerier(input.AuthorizationKeeper, input.Cdc)
	now := input.Ctx.BlockHeader().Time

	// register grants
	grant, err := types.NewAuthorizationGrant(types.NewSendAuthorization(sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(123)))), now)
	require.NoError(t, err)

	grant2, err := types.NewAuthorizationGrant(types.NewGenericAuthorization(banktypes.TypeMsgSend+"2"), now)
	require.NoError(t, err)

	input.AuthorizationKeeper.SetGrant(input.Ctx, Addrs[0], Addrs[1], banktypes.TypeMsgSend, grant)
	input.AuthorizationKeeper.SetGrant(input.Ctx, Addrs[0], Addrs[1], banktypes.TypeMsgSend+"2", grant2)

	params := types.NewQueryGrantsParams(Addrs[0], Addrs[1])
	bz, err := input.Cdc.MarshalJSON(params)
	require.NoError(t, err)

	query := abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, err := querier(input.Ctx, []string{types.QueryAllGrants}, query)
	require.NoError(t, err)

	var resGrants types.AuthorizationGrants
	input.Cdc.MustUnmarshalJSON(res, &resGrants)
	require.Equal(t, types.AuthorizationGrants{grant, grant2}, resGrants)
}
