package treasury

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/terra-project/core/x/treasury/internal/keeper"
)

func TestExportInitGenesis(t *testing.T) {
	input := keeper.CreateTestInput(t)
	input.TreasuryKeeper.RecordEpochInitialIssuance(input.Ctx)
	input.TreasuryKeeper.SetRewardWeight(input.Ctx, sdk.NewDec(1123))
	input.TreasuryKeeper.SetTaxCap(input.Ctx, "foo", sdk.NewInt(1234))
	input.TreasuryKeeper.SetTaxRate(input.Ctx, sdk.NewDec(5435))
	input.TreasuryKeeper.SetEpochTaxProceeds(input.Ctx, sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(923))))
	input.TreasuryKeeper.SetTR(input.Ctx, int64(0), sdk.NewDec(123))
	input.TreasuryKeeper.SetTR(input.Ctx, int64(1), sdk.NewDec(345))
	input.TreasuryKeeper.SetTR(input.Ctx, int64(2), sdk.NewDec(567))
	input.TreasuryKeeper.SetSR(input.Ctx, int64(0), sdk.NewDec(123))
	input.TreasuryKeeper.SetSR(input.Ctx, int64(1), sdk.NewDec(345))
	input.TreasuryKeeper.SetSR(input.Ctx, int64(2), sdk.NewDec(567))
	input.TreasuryKeeper.SetTSL(input.Ctx, int64(0), sdk.NewInt(123))
	input.TreasuryKeeper.SetTSL(input.Ctx, int64(1), sdk.NewInt(345))
	input.TreasuryKeeper.SetTSL(input.Ctx, int64(2), sdk.NewInt(567))
	input.TreasuryKeeper.SetCumulatedHeight(input.Ctx, int64(123))
	genesis := ExportGenesis(input.Ctx, input.TreasuryKeeper)

	newInput := keeper.CreateTestInput(t)
	InitGenesis(newInput.Ctx, newInput.TreasuryKeeper, genesis)
	newGenesis := ExportGenesis(newInput.Ctx, newInput.TreasuryKeeper)

	require.Equal(t, genesis, newGenesis)

	// Make epoch initial issuance to zero
	tmp := genesis.EpochInitialIssuance
	genesis.EpochInitialIssuance = sdk.Coins{}

	newInput = keeper.CreateTestInput(t)
	InitGenesis(newInput.Ctx, newInput.TreasuryKeeper, genesis)
	newGenesis = ExportGenesis(newInput.Ctx, newInput.TreasuryKeeper)

	// Return back epoch initial issuance
	genesis.EpochInitialIssuance = tmp
	require.Equal(t, genesis, newGenesis)
}
