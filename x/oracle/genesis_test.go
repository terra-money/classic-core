package oracle

import (
	"github.com/terra-project/core/x/oracle/internal/types"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-project/core/x/oracle/internal/keeper"
)

func TestExportInitGenesis(t *testing.T) {
	input, _ := setup(t)

	input.OracleKeeper.SetOracleDelegate(input.Ctx, keeper.ValAddrs[0], keeper.Addrs[1])
	input.OracleKeeper.AddExchangeRatePrevote(input.Ctx, NewExchangeRatePrevote(VoteHash{123}, "denom", sdk.ValAddress{}, int64(2)))
	input.OracleKeeper.AddExchangeRateVote(input.Ctx, NewExchangeRateVote(sdk.NewDec(1), "denom", sdk.ValAddress{}))
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, "denom", sdk.NewDec(123))
	input.OracleKeeper.AddAggregateExchangeRatePrevote(input.Ctx, NewAggregateExchangeRatePrevote(AggregateVoteHash{123}, sdk.ValAddress{}, int64(2)))
	input.OracleKeeper.AddAggregateExchangeRateVote(input.Ctx, NewAggregateExchangeRateVote(types.ExchangeRateTuples{{Denom: "foo", ExchangeRate: sdk.NewDec(123)}}, sdk.ValAddress{}))
	input.OracleKeeper.SetTobinTax(input.Ctx, "denom", sdk.NewDecWithPrec(123, 3))
	input.OracleKeeper.SetTobinTax(input.Ctx, "denom2", sdk.NewDecWithPrec(123, 3))
	genesis := ExportGenesis(input.Ctx, input.OracleKeeper)

	newInput := keeper.CreateTestInput(t)
	InitGenesis(newInput.Ctx, newInput.OracleKeeper, genesis)
	newGenesis := ExportGenesis(newInput.Ctx, newInput.OracleKeeper)

	require.Equal(t, genesis, newGenesis)
}
