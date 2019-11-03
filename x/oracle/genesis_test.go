package oracle

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/oracle/internal/keeper"
)

func TestExportInitGenesis(t *testing.T) {
	input, h := setup(t)

	makePrevoteAndVote(t, input, h, 1, core.MicroSDRDenom, randomExchangeRate, 0)
	makePrevoteAndVote(t, input, h, 1, core.MicroSDRDenom, randomExchangeRate, 1)
	makePrevoteAndVote(t, input, h, 1, core.MicroSDRDenom, randomExchangeRate.MulInt64(10000), 2)

	EndBlocker(input.Ctx.WithBlockHeight(1), input.OracleKeeper)

	input.OracleKeeper.SetOracleDelegate(input.Ctx, keeper.ValAddrs[0], keeper.Addrs[1])
	input.OracleKeeper.AddExchangeRatePrevote(input.Ctx, NewExchangeRatePrevote("1234", "denom", sdk.ValAddress{}, int64(2)))
	input.OracleKeeper.AddExchangeRateVote(input.Ctx, NewExchangeRateVote(sdk.NewDec(1), "denom", sdk.ValAddress{}))
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, "denom", sdk.NewDec(123))
	genesis := ExportGenesis(input.Ctx, input.OracleKeeper)

	newInput := keeper.CreateTestInput(t)
	InitGenesis(newInput.Ctx, newInput.OracleKeeper, genesis)
	newGenesis := ExportGenesis(newInput.Ctx, newInput.OracleKeeper)

	require.Equal(t, genesis, newGenesis)
}
