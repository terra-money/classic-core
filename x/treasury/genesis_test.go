package treasury

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/terra-project/core/x/treasury/internal/keeper"
)

func TestExportInitGenesis(t *testing.T) {
	input := keeper.CreateTestInput(t)
	input.TreasuryKeeper.SetHistoricalIssuance(input.Ctx, 0, sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(33))))
	input.TreasuryKeeper.SetHistoricalIssuance(input.Ctx, 1, sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(123))))
	input.TreasuryKeeper.SetHistoricalIssuance(input.Ctx, 2, sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(131))))
	input.TreasuryKeeper.SetRewardWeight(input.Ctx, sdk.NewDec(1123))
	input.TreasuryKeeper.SetTaxCap(input.Ctx, "foo", sdk.NewInt(1234))
	input.TreasuryKeeper.SetTaxRate(input.Ctx, sdk.NewDec(5435))
	input.TreasuryKeeper.SetTaxProceeds(input.Ctx, 0, sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(923))))
	input.TreasuryKeeper.SetTaxProceeds(input.Ctx, 1, sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(123))))
	input.TreasuryKeeper.SetTaxProceeds(input.Ctx, 2, sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(623))))
	genesis := ExportGenesis(input.Ctx, input.TreasuryKeeper)

	newInput := keeper.CreateTestInput(t)
	InitGenesis(newInput.Ctx, newInput.TreasuryKeeper, genesis)
	newGenesis := ExportGenesis(newInput.Ctx, newInput.TreasuryKeeper)

	require.Equal(t, genesis, newGenesis)
}
