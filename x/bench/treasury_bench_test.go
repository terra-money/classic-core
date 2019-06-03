package bench

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-project/core/types/assets"
	"github.com/terra-project/core/types/util"
	"github.com/terra-project/core/x/treasury"
)

func BenchmarkTreasuryUpdatePerEpoch(b *testing.B) {
	input := createTestInput()

	taxAmount := sdk.NewInt(1000).MulRaw(assets.MicroUnit)

	// Set random prices
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.MicroSDRDenom, sdk.NewDec(1))
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.MicroKRWDenom, sdk.NewDec(10))
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.MicroGBPDenom, sdk.NewDec(100))
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.MicroCNYDenom, sdk.NewDec(1000))

	params := input.treasuryKeeper.GetParams(input.ctx)
	probationEpoch := params.WindowProbation.Int64()

	b.ResetTimer()
	for i := int64(0); i < int64(b.N)+probationEpoch; i++ {

		input.ctx = input.ctx.WithBlockHeight(i*util.BlocksPerEpoch - 1)
		input.mintKeeper.AddSeigniorage(input.ctx, uLunaAmt)

		input.treasuryKeeper.RecordTaxProceeds(input.ctx, sdk.Coins{
			sdk.NewCoin(assets.MicroSDRDenom, taxAmount),
			sdk.NewCoin(assets.MicroKRWDenom, taxAmount),
			sdk.NewCoin(assets.MicroGBPDenom, taxAmount),
			sdk.NewCoin(assets.MicroCNYDenom, taxAmount),
		})

		treasury.EndBlocker(input.ctx, input.treasuryKeeper)
	}
}
