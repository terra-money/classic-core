package keeper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/terra-money/core/x/msgauth/internal/types"
)

type TestSuite struct {
	suite.Suite
	ctx           sdk.Context
	accountKeeper auth.AccountKeeper
	paramsKeeper  params.Keeper
	bankKeeper    bank.Keeper
	keeper        Keeper
	router        sdk.Router
}

func (s *TestSuite) SetupTest() {
	s.ctx, s.accountKeeper, s.paramsKeeper, s.bankKeeper, s.keeper, s.router = SetupTestInput()
}

func (s *TestSuite) TestKeeper() {
	err := s.bankKeeper.SetCoins(s.ctx, granterAddr, sdk.NewCoins(sdk.NewInt64Coin("steak", 10000)))
	s.Require().Nil(err)
	s.Require().True(s.bankKeeper.GetCoins(s.ctx, granterAddr).AmountOf("steak").Equal(sdk.NewInt(10000)))

	s.T().Log("verify that no authorization returns nil")
	_, found := s.keeper.GetGrant(s.ctx, granterAddr, granteeAddr, bank.MsgSend{}.Type())
	s.Require().False(found)
	now := s.ctx.BlockHeader().Time
	s.Require().NotNil(now)

	newCoins := sdk.NewCoins(sdk.NewInt64Coin("steak", 100))

	s.T().Log("verify if authorization is accepted")
	x := types.NewAuthorizationGrant(types.SendAuthorization{SpendLimit: newCoins}, now.Add(time.Hour))
	s.keeper.SetGrant(s.ctx, granterAddr, granteeAddr, x)
	grant, _ := s.keeper.GetGrant(s.ctx, granterAddr, granteeAddr, bank.MsgSend{}.Type())
	s.Require().NotNil(grant)
	s.Require().Equal(grant.Authorization.MsgType(), bank.MsgSend{}.Type())

	s.T().Log("verify fetching authorization with wrong msg type fails")
	_, found = s.keeper.GetGrant(s.ctx, granterAddr, granteeAddr, bank.MsgMultiSend{}.Type())
	s.Require().False(found)

	s.T().Log("verify fetching authorization with wrong grantee fails")
	_, found = s.keeper.GetGrant(s.ctx, granterAddr, recipientAddr, bank.MsgSend{}.Type())
	s.Require().False(found)

	grants := s.keeper.GetGrants(s.ctx, granterAddr, granteeAddr)
	s.Require().Equal(1, len(grants))

	s.keeper.IterateGrants(s.ctx, func(
		granter, grantee sdk.AccAddress, grant types.AuthorizationGrant,
	) bool {
		s.Require().Equal(granterAddr, granter)
		s.Require().Equal(granteeAddr, grantee)
		s.Require().Equal(x, grant)
		return false
	})

	s.T().Log("")

	s.T().Log("verify revoke fails with wrong information")
	s.keeper.RevokeGrant(s.ctx, granterAddr, recipientAddr, bank.MsgSend{}.Type())
	_, found = s.keeper.GetGrant(s.ctx, granterAddr, granteeAddr, bank.MsgSend{}.Type())
	s.Require().True(found)

	s.T().Log("verify revoke executes with correct information")
	s.keeper.RevokeGrant(s.ctx, granterAddr, granteeAddr, bank.MsgSend{}.Type())
	_, found = s.keeper.GetGrant(s.ctx, granterAddr, granteeAddr, bank.MsgSend{}.Type())
	s.Require().False(found)

}

func (s *TestSuite) TestKeeperFees() {
	err := s.bankKeeper.SetCoins(s.ctx, granterAddr, sdk.NewCoins(sdk.NewInt64Coin("steak", 10000)))
	s.Require().Nil(err)
	s.Require().True(s.bankKeeper.GetCoins(s.ctx, granterAddr).AmountOf("steak").Equal(sdk.NewInt(10000)))

	now := s.ctx.BlockHeader().Time
	s.Require().NotNil(now)

	smallCoins := sdk.NewCoins(sdk.NewInt64Coin("steak", 2))
	coins := sdk.NewCoins(sdk.NewInt64Coin("steak", 20))
	largeCoins := sdk.NewCoins(sdk.NewInt64Coin("steak", 123))
	//lotCoin := sdk.NewCoins(sdk.NewInt64Coin("steak", 4567))

	msgs := types.MsgExecAuthorized{
		Grantee: granteeAddr,
		Msgs: []sdk.Msg{
			bank.MsgSend{
				Amount:      smallCoins,
				FromAddress: granterAddr,
				ToAddress:   recipientAddr,
			},
		},
	}

	s.T().Log("verify dispatch fails with invalid authorization")
	error := s.keeper.DispatchActions(s.ctx, granteeAddr, msgs.Msgs)
	s.Require().Error(error)

	s.T().Log("verify dispatch executes with correct information")
	// grant authorization

	s.keeper.SetGrant(s.ctx, granterAddr, granteeAddr, types.NewAuthorizationGrant(types.SendAuthorization{SpendLimit: coins}, now))
	grant, found := s.keeper.GetGrant(s.ctx, granterAddr, granteeAddr, bank.MsgSend{}.Type())
	s.Require().NotNil(grant)
	s.Require().True(found)
	s.Require().Equal(grant.Authorization.MsgType(), bank.MsgSend{}.Type())
	error = s.keeper.DispatchActions(s.ctx, granteeAddr, msgs.Msgs)
	s.Require().Nil(error)

	_, found = s.keeper.GetGrant(s.ctx, granterAddr, granteeAddr, bank.MsgSend{}.Type())
	s.Require().True(found)

	s.T().Log("verify dispatch fails with overlimit")

	msgs = types.MsgExecAuthorized{
		Grantee: granteeAddr,
		Msgs: []sdk.Msg{
			bank.MsgSend{
				Amount:      largeCoins,
				FromAddress: granterAddr,
				ToAddress:   recipientAddr,
			},
		},
	}

	error = s.keeper.DispatchActions(s.ctx, granteeAddr, msgs.Msgs)
	s.Require().Error(error)

	_, found = s.keeper.GetGrant(s.ctx, granterAddr, granteeAddr, bank.MsgSend{}.Type())
	s.Require().True(found)

	s.T().Log("verify dispatch success and revoke grant which is out of limit")

	msgs = types.MsgExecAuthorized{
		Grantee: granteeAddr,
		Msgs: []sdk.Msg{
			bank.MsgSend{
				Amount:      coins.Sub(smallCoins),
				FromAddress: granterAddr,
				ToAddress:   recipientAddr,
			},
		},
	}

	error = s.keeper.DispatchActions(s.ctx, granteeAddr, msgs.Msgs)
	s.Require().NoError(error)

	_, found = s.keeper.GetGrant(s.ctx, granterAddr, granteeAddr, bank.MsgSend{}.Type())
	s.Require().False(found)
}

func (s *TestSuite) TestGrantQueue() {
	now := s.ctx.BlockTime()
	s.keeper.InsertGrantQueue(s.ctx, granterAddr, granteeAddr, bank.MsgSend{}.Type(), now.Add(time.Hour))
	s.keeper.InsertGrantQueue(s.ctx, granterAddr, granteeAddr, bank.MsgSend{}.Type()+"2", now.Add(time.Hour))

	ggmPairs := []types.GGMPair{
		{
			GranterAddress: granterAddr,
			GranteeAddress: granteeAddr,
			MsgType:        bank.MsgSend{}.Type(),
		},
		{
			GranterAddress: granterAddr,
			GranteeAddress: granteeAddr,
			MsgType:        bank.MsgSend{}.Type() + "2",
		},
	}

	timeSlice := s.keeper.GetGrantQueueTimeSlice(s.ctx, now)
	s.Require().Equal(0, len(timeSlice))

	timeSlice = s.keeper.GetGrantQueueTimeSlice(s.ctx, now.Add(time.Hour))
	s.Require().Equal(ggmPairs, timeSlice)

	allPairs := s.keeper.DequeueAllMatureGrantQueue(s.ctx.WithBlockTime(now))
	s.Require().Equal(0, len(allPairs))

	allPairs = s.keeper.DequeueAllMatureGrantQueue(s.ctx.WithBlockTime(now.Add(time.Hour)))
	s.Require().Equal(ggmPairs, allPairs)

	s.keeper.InsertGrantQueue(s.ctx, granterAddr, granteeAddr, bank.MsgSend{}.Type(), now.Add(time.Hour))
	s.keeper.RevokeFromGrantQueue(s.ctx, granterAddr, granteeAddr, bank.MsgSend{}.Type(), now.Add(time.Hour))
	timeSlice = s.keeper.GetGrantQueueTimeSlice(s.ctx, now.Add(time.Hour))
	s.Require().Equal(0, len(timeSlice))
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
