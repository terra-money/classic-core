package market

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/market/internal/keeper"
)

func TestOracleThreshold(t *testing.T) {
	input := keeper.CreateTestInput(t)

	targetIssuance := sdk.NewInt(1000000)
	supply := input.SupplyKeeper.GetSupply(input.Ctx)
	supply = supply.SetTotal(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, targetIssuance)))
	input.SupplyKeeper.SetSupply(input.Ctx, supply)

	input.Ctx = input.Ctx.WithBlockHeight(core.BlocksPerDay - 1)
	EndBlocker(input.Ctx, input.MarketKeeper)
	issuance := input.MarketKeeper.GetPrevDayIssuance(input.Ctx).AmountOf(core.MicroLunaDenom)
	require.Equal(t, targetIssuance, issuance)
}
