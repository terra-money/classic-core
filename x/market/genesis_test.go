package market

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/terra-money/core/x/market/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestExportInitGenesis(t *testing.T) {
	input := keeper.CreateTestInput(t)
	input.MarketKeeper.SetMintPoolDelta(input.Ctx, sdk.NewDec(1123))
	input.MarketKeeper.SetBurnPoolDelta(input.Ctx, sdk.NewDec(3321))
	genesis := ExportGenesis(input.Ctx, input.MarketKeeper)

	newInput := keeper.CreateTestInput(t)
	InitGenesis(newInput.Ctx, newInput.MarketKeeper, genesis)
	newGenesis := ExportGenesis(newInput.Ctx, newInput.MarketKeeper)

	require.Equal(t, genesis, newGenesis)
}
