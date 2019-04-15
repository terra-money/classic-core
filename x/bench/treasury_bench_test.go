package bench

import (
	"testing"
)

func BenchmarkTreasuryUpdatePerBlock(b *testing.B) {
	// input := createTestInput()

	// defaultTreasuryParams := treasury.DefaultParams()

	// h := oracle.NewHandler(input.oracleKeeper)

	// denoms := []string{
	// 	assets.MicroSDRDenom,
	// 	assets.MicroKRWDenom,
	// 	assets.MicroUSDDenom,
	// 	assets.MicroCNYDenom,
	// 	assets.MicroJPYDenom,
	// 	assets.MicroGBPDenom,
	// 	assets.MicroEURDenom,
	// }

	// b.ResetTimer()
	// for i := 0; i < b.N; i++ {
	// 	ctx := input.ctx.WithBlockHeight(int64(i))

	// 	for j := 0; j < numOfValidators; j++ {
	// 		for d := 0; d < len(denoms); d++ {
	// 			voteMsg := oracle.NewMsgPriceFeed(denoms[d], sdk.NewDec(1), addrs[j])

	// 			res := h(ctx, voteMsg)
	// 			if !res.IsOK() {
	// 				panic(res.Log)
	// 			}
	// 		}
	// 	}

	// 	oracle.EndBlocker(ctx, input.oracleKeeper)
	// }
}
