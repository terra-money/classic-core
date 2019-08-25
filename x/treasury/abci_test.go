package treasury

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/treasury/internal/keeper"
)

func TestEndBlockerIssuanceUpdate(t *testing.T) {
	input := keeper.CreateTestInput(t)

	targetIssuance := sdk.NewInt(1000)
	input.Ctx = input.Ctx.WithBlockHeight(core.BlocksPerEpoch - 1)
	supply := input.SupplyKeeper.GetSupply(input.Ctx)
	supply = supply.SetTotal(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, targetIssuance)))
	input.SupplyKeeper.SetSupply(input.Ctx, supply)
	EndBlocker(input.Ctx, input.TreasuryKeeper)

	issuance := input.TreasuryKeeper.GetHistoricalIssuance(input.Ctx, 0).AmountOf(core.MicroLunaDenom)
	require.Equal(t, targetIssuance, issuance)
}
