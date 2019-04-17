package bench

import (
	"fmt"
	"math/rand"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-project/core/x/budget"
)

func BenchmarkSubmitAndVoteProgramsPerBlock(b *testing.B) {
	const numOfPrograms = 5
	input := createTestInput()

	defaultBudgetParams := budget.DefaultParams()
	defaultBudgetParams.VotePeriod = 1
	input.budgetKeeper.SetParams(input.ctx, defaultBudgetParams)

	h := budget.NewHandler(input.budgetKeeper)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := input.ctx.WithBlockHeight(int64(i))

		for p := 0; p < numOfPrograms; p++ {
			// registers programs
			msg := budget.NewMsgSubmitProgram(fmt.Sprintf("test-%d-%d", i, p), "description", addrs[p], addrs[p+1])
			res := h(ctx, msg)

			if !res.IsOK() {
				panic(res.Log)
			}

			for v := 0; v < numOfValidators; v++ {
				// votes to registered program
				voteMsg := budget.NewMsgVoteProgram(uint64(i*numOfPrograms+p+1), rand.Intn(2) == 0, addrs[v])
				res := h(ctx, voteMsg)

				if !res.IsOK() {
					panic(res.Log)
				}
			}
		}

		budget.EndBlocker(ctx, input.budgetKeeper)
	}

}

func BenchmarkSubmitAndWithdrawProgramPerBlock(b *testing.B) {
	input := createTestInput()

	h := budget.NewHandler(input.budgetKeeper)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := input.ctx.WithBlockHeight(int64(i))

		for v := 0; v < numOfValidators; v++ {
			var msg sdk.Msg

			if i%2 == 0 {
				msg = budget.NewMsgSubmitProgram(fmt.Sprintf("test-%d-%d", i, v), "description", addrs[v], addrs[v])
			} else {
				msg = budget.NewMsgWithdrawProgram(uint64((i/2)*numOfValidators+v+1), addrs[v])
			}

			res := h(ctx, msg)
			if !res.IsOK() {
				panic(res.Log)
			}
		}

		budget.EndBlocker(ctx, input.budgetKeeper)
	}
}
