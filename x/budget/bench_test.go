package budget

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BenchmarkOneSubmitterSubmitProgramPerBlock(b *testing.B) {
	input := createTestInput()

	h := NewHandler(input.budgetKeeper)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := input.ctx.WithBlockHeight(int64(i))
		msg := NewMsgSubmitProgram(fmt.Sprintf("test-%d", i), "description", addrs[0], addrs[1])
		res := h(ctx, msg)

		if !res.IsOK() {
			panic("program submission broken")
		}

		EndBlocker(ctx, input.budgetKeeper)
	}
}

func BenchmarkOneSubmitterSubmitAndWithdrawProgram(b *testing.B) {
	input := createTestInput()

	h := NewHandler(input.budgetKeeper)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := input.ctx.WithBlockHeight(int64(i))

		var res sdk.Result
		var msg sdk.Msg
		if i%2 == 0 {
			msg = NewMsgSubmitProgram(fmt.Sprintf("test-%d", i), "description", addrs[0], addrs[1])
		} else {
			msg = NewMsgWithdrawProgram(uint64(i/2+1), addrs[0])
		}

		res = h(ctx, msg)

		if !res.IsOK() {
			panic(res.Log)
		}

		EndBlocker(ctx, input.budgetKeeper)
	}
}
