package msgauth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/terra-project/core/x/msgauth/internal/keeper"
	"github.com/terra-project/core/x/msgauth/internal/types"
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
	router        baseapp.Router
	handler       sdk.Handler
}

func (s *TestSuite) SetupTest() {
	s.ctx, s.accountKeeper, s.paramsKeeper, s.bankKeeper, s.keeper, s.router = keeper.SetupTestInput()
	s.handler = NewHandler(s.keeper)
}

func (s *TestSuite) TestGrant() {
	coins := sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(1_000_000_000)))
	cap := types.SendAuthorization{SpendLimit: coins}
	expiration := time.Now().Add(time.Hour)
	msg := types.NewMsgGrantAuthorization(granterAddr, granteeAddr, cap, expiration)

	_, err := s.handler(s.ctx, msg)
	s.Require().NoError(err)

	resCap, resExpiration := s.keeper.GetAuthorization(s.ctx, granteeAddr, granterAddr, bank.MsgSend{}.Type())
	s.Require().Equal(cap, resCap)
	s.Require().Equal(expiration.Unix(), resExpiration)
}

func (s *TestSuite) TestRevoke() {
	coins := sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(1_000_000_000)))
	grantMsg := types.NewMsgGrantAuthorization(granterAddr, granteeAddr, types.SendAuthorization{
		SpendLimit: coins,
	}, time.Now().Add(time.Hour))

	_, err := s.handler(s.ctx, grantMsg)
	s.Require().NoError(err)

	revokeMsg := types.NewMsgRevokeAuthorization(granterAddr, granteeAddr, bank.MsgSend{}.Type())
	_, err = s.handler(s.ctx, revokeMsg)
	s.Require().NoError(err)

	res, expiration := s.keeper.GetAuthorization(s.ctx, granteeAddr, granterAddr, bank.MsgSend{}.Type())
	s.Require().Nil(res)
	s.Require().Equal(int64(0), expiration)
}

func (s *TestSuite) TestExecute() {
	coins := sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(1_000_000_000)))
	s.bankKeeper.SetCoins(s.ctx, granterAddr, coins)

	grantMsg := types.NewMsgGrantAuthorization(granterAddr, granteeAddr, types.SendAuthorization{
		SpendLimit: coins,
	}, time.Now().Add(time.Hour))

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
