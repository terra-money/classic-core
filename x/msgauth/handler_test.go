package msgauth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/crypto/ed25519"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/terra-money/core/x/msgauth/internal/keeper"
	"github.com/terra-money/core/x/msgauth/internal/types"
)

var (
	granteePub    = ed25519.GenPrivKey().PubKey()
	granterPub    = ed25519.GenPrivKey().PubKey()
	recipientPub  = ed25519.GenPrivKey().PubKey()
	granteeAddr   = sdk.AccAddress(granteePub.Address())
	granterAddr   = sdk.AccAddress(granterPub.Address())
	recipientAddr = sdk.AccAddress(recipientPub.Address())
)

type TestSuite struct {
	suite.Suite
	ctx           sdk.Context
	accountKeeper auth.AccountKeeper
	paramsKeeper  params.Keeper
	bankKeeper    bank.Keeper
	keeper        Keeper
	router        sdk.Router
	handler       sdk.Handler
}

func (s *TestSuite) SetupTest() {
	s.ctx, s.accountKeeper, s.paramsKeeper, s.bankKeeper, s.keeper, s.router = keeper.SetupTestInput()
	s.handler = NewHandler(s.keeper)
}

func (s *TestSuite) TestGrant() {
	coins := sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(1_000_000_000)))

	// send authorization
	sendAuth := types.SendAuthorization{SpendLimit: coins}
	msg := types.NewMsgGrantAuthorization(granterAddr, granteeAddr, sendAuth, time.Hour)

	_, err := s.handler(s.ctx, msg)
	s.Require().NoError(err)

	grant, found := s.keeper.GetGrant(s.ctx, granterAddr, granteeAddr, sendAuth.MsgType())
	s.Require().True(found)
	s.Require().Equal(sendAuth, grant.Authorization)
	s.Require().Equal(s.ctx.BlockTime().Add(time.Hour), grant.Expiration)

	// generic authorization
	genericAuth := types.NewGenericAuthorization("swap")
	msg = types.NewMsgGrantAuthorization(granterAddr, granteeAddr, genericAuth, time.Hour)

	_, err = s.handler(s.ctx, msg)
	s.Require().NoError(err)

	grant, found = s.keeper.GetGrant(s.ctx, granterAddr, granteeAddr, "swap")
	s.Require().True(found)
	s.Require().Equal(genericAuth, grant.Authorization)
	s.Require().Equal(s.ctx.BlockTime().Add(time.Hour), grant.Expiration)

	// test not allowed to grant
	genericAuth = types.NewGenericAuthorization("now allowed msg")
	msg = types.NewMsgGrantAuthorization(granterAddr, granteeAddr, genericAuth, time.Hour)

	_, err = s.handler(s.ctx, msg)
	s.Require().Error(err)
}

func (s *TestSuite) TestRevoke() {
	coins := sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(1_000_000_000)))
	grantMsg := types.NewMsgGrantAuthorization(granterAddr, granteeAddr, types.SendAuthorization{
		SpendLimit: coins,
	}, time.Hour)

	_, err := s.handler(s.ctx, grantMsg)
	s.Require().NoError(err)

	revokeMsg := types.NewMsgRevokeAuthorization(granterAddr, granteeAddr, bank.MsgSend{}.Type())
	_, err = s.handler(s.ctx, revokeMsg)
	s.Require().NoError(err)

	_, found := s.keeper.GetGrant(s.ctx, granteeAddr, granterAddr, bank.MsgSend{}.Type())
	s.Require().False(found)
}

func (s *TestSuite) TestExecute() {
	coins := sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(1_000_000_000)))
	s.bankKeeper.SetCoins(s.ctx, granterAddr, coins)

	grantMsg := types.NewMsgGrantAuthorization(granterAddr, granteeAddr, types.SendAuthorization{
		SpendLimit: coins,
	}, time.Hour)

	_, err := s.handler(s.ctx, grantMsg)
	s.Require().NoError(err)

	execMsg := types.NewMsgExecAuthorized(granteeAddr, []sdk.Msg{
		bank.NewMsgSend(granterAddr, granteeAddr, coins),
	})

	_, err = s.handler(s.ctx, execMsg)
	s.Require().NoError(err)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
