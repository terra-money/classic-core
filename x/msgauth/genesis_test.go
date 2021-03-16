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

func TestGenesisExportImport(t *testing.T) {
	input := keeper.CreateTestInput(t)
	coins := sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(1_000_000_000)))

	now := input.Ctx.BlockHeader().Time
	grant, err := types.NewAuthorizationGrant(types.NewSendAuthorization(coins), now.Add(time.Hour))
	require.NoError(t, err)

	input.AuthorizationKeeper.SetGrant(input.Ctx, keeper.Addrs[0], keeper.Addrs[1], banktypes.TypeMsgSend, grant)
	genesis := ExportGenesis(input.Ctx, input.AuthorizationKeeper)

	// Clear keeper
	input.AuthorizationKeeper.RevokeGrant(input.Ctx, keeper.Addrs[0], keeper.Addrs[1], banktypes.TypeMsgSend)

	InitGenesis(input.Ctx, input.AuthorizationKeeper, genesis)
	newGenesis := ExportGenesis(input.Ctx, input.AuthorizationKeeper)

	require.Equal(t, genesis, newGenesis)
}
