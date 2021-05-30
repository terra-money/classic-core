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

func (s *TestSuite) TestGenesisExportImport() {
	coins := sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(1_000_000_000)))

	now := s.ctx.BlockHeader().Time
	grant := NewAuthorizationGrant(types.SendAuthorization{SpendLimit: coins}, now.Add(time.Hour))
	s.keeper.SetGrant(s.ctx, granterAddr, granteeAddr, grant)
	genesis := ExportGenesis(s.ctx, s.keeper)

	// Clear keeper
	s.keeper.RevokeGrant(s.ctx, granterAddr, granteeAddr, grant.Authorization.MsgType())

	InitGenesis(s.ctx, s.keeper, genesis)
	newGenesis := ExportGenesis(s.ctx, s.keeper)

	s.Require().Equal(genesis, newGenesis)
}
