package oracle

import (
	"terra/types/assets"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BenchmarkOracleFeedVotePerBlock(b *testing.B) {
	input := createTestInput()

	defaultOracleParams := DefaultParams()
	defaultOracleParams.VotePeriod = sdk.OneInt()
	input.oracleKeeper.SetParams(input.ctx, defaultOracleParams)

	h := NewHandler(input.oracleKeeper)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := input.ctx.WithBlockHeight(int64(i))
		msg1 := NewMsgPriceFeed(assets.KRWDenom, sdk.NewDec(1), sdk.AccAddress(input.valset.Validators[0].Address.Bytes()))
		res1 := h(ctx, msg1)
		if !res1.IsOK() {
			panic(res1.Log)
		}

		msg2 := NewMsgPriceFeed(assets.KRWDenom, sdk.NewDec(2), sdk.AccAddress(input.valset.Validators[1].Address.Bytes()))
		res2 := h(ctx, msg2)
		if !res2.IsOK() {
			panic(res2.Log)
		}

		msg3 := NewMsgPriceFeed(assets.KRWDenom, sdk.NewDec(3), sdk.AccAddress(input.valset.Validators[2].Address.Bytes()))
		res3 := h(ctx, msg3)
		if !res3.IsOK() {
			panic(res3.Log)
		}

		EndBlocker(ctx, input.oracleKeeper)
	}
}
