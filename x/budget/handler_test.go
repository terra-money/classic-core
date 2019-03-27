package budget

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandlerMsgSubmitProgram(t *testing.T) {
	input := createTestInput(t)

	h := NewHandler(input.budgetKeeper)

	// Regular submit msg passes
	msg := NewSubmitProgramMsg("test", "testdescription", addrs[0], addrs[1])
	res := h(input.ctx, msg)
	require.True(t, res.IsOK())

	// Everything else should be tested in validateMsg ... so skip
}

func TestHandlerMsgWithdrawProgram(t *testing.T) {
	input := createTestInput(t)

	h := NewHandler(input.budgetKeeper)

	// Submit program
	submitMsg := NewSubmitProgramMsg("test", "testdescription", addrs[0], addrs[1])
	res := h(input.ctx, submitMsg)
	require.True(t, res.IsOK())

	// Withdrawing submitted program works
	withdrawMsg := NewMsgWithdrawProgram(0, addrs[0])
	res = h(input.ctx, withdrawMsg)
	require.True(t, res.IsOK())

	// Withdrawing again doesn't work
	withdrawMsg = NewMsgWithdrawProgram(0, addrs[0])
	res = h(input.ctx, withdrawMsg)
	require.False(t, res.IsOK())

	// Withdrawing from a different submitter address doesn't work
	withdrawMsg = NewMsgWithdrawProgram(0, addrs[2])
	res = h(input.ctx, withdrawMsg)
	require.False(t, res.IsOK())

	// Withdrawing an unsubmitted program doesn't work
	withdrawMsg = NewMsgWithdrawProgram(4, addrs[2])
	res = h(input.ctx, withdrawMsg)
	require.False(t, res.IsOK())
}

func TestHandlerMsgVoteCandidate(t *testing.T) {
	input := createTestInput(t)

	h := NewHandler(input.budgetKeeper)

	// Submit program
	submitMsg := NewSubmitProgramMsg("test", "testdescription", addrs[0], addrs[1])
	res := h(input.ctx, submitMsg)
	require.True(t, res.IsOK())

	// Voting on a submitted program works
	voteMsg := NewMsgVoteProgram(0, true, addrs[0])
	res = h(input.ctx, voteMsg)
	require.True(t, res.IsOK())

	// Voting on an un submitted program doesn't work
	voteMsg = NewMsgVoteProgram(4, true, addrs[0])
	res = h(input.ctx, voteMsg)
	require.False(t, res.IsOK())
}
