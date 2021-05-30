package msgauth

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	"github.com/terra-money/core/x/msgauth/internal/types"
)

func init() {
	_ = suite.Suite{}
}

func (s *TestSuite) TestMature() {
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

	EndBlocker(s.ctx.WithBlockTime(s.ctx.BlockTime().Add(time.Hour)), s.keeper)
	_, found = s.keeper.GetGrant(s.ctx, granterAddr, granteeAddr, sendAuth.MsgType())
	s.Require().False(found)
}
