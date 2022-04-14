package market

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/terra-money/core/x/market/keeper"
	"github.com/terra-money/core/x/market/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestExportInitGenesis(t *testing.T) {
	input := keeper.CreateTestInput(t)
	input.MarketKeeper.SetTerraPoolDelta(input.Ctx, sdk.NewDec(1123))
	input.MarketKeeper.SetSeigniorageRoutes(input.Ctx, []types.SeigniorageRoute{
		{Address: keeper.Addrs[0].String(), Weight: sdk.NewDecWithPrec(1, 1)},
		{Address: keeper.Addrs[1].String(), Weight: sdk.NewDecWithPrec(2, 1)},
	})
	genesis := ExportGenesis(input.Ctx, input.MarketKeeper)

	newInput := keeper.CreateTestInput(t)
	InitGenesis(newInput.Ctx, newInput.MarketKeeper, genesis)
	newGenesis := ExportGenesis(newInput.Ctx, newInput.MarketKeeper)

	require.Equal(t, genesis, newGenesis)
}
