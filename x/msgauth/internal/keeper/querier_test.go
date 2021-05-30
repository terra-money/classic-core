package keeper

import (
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/core/x/msgauth/internal/types"
)

func init() {
	_ = suite.Suite{}
}

func (s *TestSuite) TestNewQuerier() {
	querier := NewQuerier(s.keeper)

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	_, err := querier(s.ctx, []string{"INVALID_PATH"}, query)
	s.Require().Error(err)
}

func (s *TestSuite) TestQueryGrant() {
	querier := NewQuerier(s.keeper)
	now := s.ctx.BlockHeader().Time

	// register grant
	grant := types.NewAuthorizationGrant(types.NewSendAuthorization(sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(123)))), now)
	s.keeper.SetGrant(s.ctx, granterAddr, granteeAddr, grant)

	params := types.NewQueryGrantParams(granterAddr, granteeAddr, types.SendAuthorization{}.MsgType())
	bz, err := s.keeper.cdc.MarshalJSON(params)
	s.Require().NoError(err)

	query := abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, err := querier(s.ctx, []string{types.QueryGrant}, query)
	s.Require().NoError(err)

	var resGrant types.AuthorizationGrant
	s.keeper.cdc.MustUnmarshalJSON(res, &resGrant)
	s.Require().Equal(grant, resGrant)
}

func (s *TestSuite) TestQueryGrants() {
	querier := NewQuerier(s.keeper)
	now := s.ctx.BlockHeader().Time

	// register grants
	grant := types.NewAuthorizationGrant(types.NewSendAuthorization(sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(123)))), now)
	grant2 := types.NewAuthorizationGrant(types.NewGenericAuthorization(types.SendAuthorization{}.MsgType()+"2"), now)
	s.keeper.SetGrant(s.ctx, granterAddr, granteeAddr, grant)
	s.keeper.SetGrant(s.ctx, granterAddr, granteeAddr, grant2)

	params := types.NewQueryGrantsParams(granterAddr, granteeAddr)
	bz, err := s.keeper.cdc.MarshalJSON(params)
	s.Require().NoError(err)

	query := abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, err := querier(s.ctx, []string{types.QueryGrants}, query)
	s.Require().NoError(err)

	var resGrants []types.AuthorizationGrant
	s.keeper.cdc.MustUnmarshalJSON(res, &resGrants)
	s.Require().Equal([]types.AuthorizationGrant{grant, grant2}, resGrants)
}
