package budget

import (
	"math/rand"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func TestKeeperProgramID(t *testing.T) {
	input := createTestInput(t)

	// Program ids start at 0 and increment by 1 on each request
	numTests := 10
	for i := 0; i < numTests; i++ {
		id := input.budgetKeeper.NewProgramID(input.ctx)
		require.Equal(t, id, uint64(i))
	}
}

func TestKeeperDeposit(t *testing.T) {
	input := createTestInput(t)

	// Set the balance to equal the default deposit
	deposit := sdk.Coins{input.budgetKeeper.GetParams(input.ctx).Deposit}
	input.bankKeeper.SetCoins(input.ctx, addrs[0], deposit)

	// addr0 has enough coins to pay the deposit
	err := input.budgetKeeper.PayDeposit(input.ctx, addrs[0])
	require.Nil(t, err)

	// Doesn't have enough coins to pay the deposit
	err = input.budgetKeeper.PayDeposit(input.ctx, addrs[0])
	require.NotNil(t, err)

	// Refund works
	err = input.budgetKeeper.RefundDeposit(input.ctx, addrs[0])
	require.Nil(t, err)

	// After refund, addr0's balance equals the deposit he paid previously
	balance := input.bankKeeper.GetCoins(input.ctx, addrs[0])
	require.Equal(t, balance, deposit)
}

func TestKeeperParams(t *testing.T) {
	input := createTestInput(t)

	defaultParams := DefaultParams()
	input.budgetKeeper.SetParams(input.ctx, defaultParams)

	retrievedParams := input.budgetKeeper.GetParams(input.ctx)
	require.Equal(t, defaultParams, retrievedParams)
}

func TestKeeperProgram(t *testing.T) {
	input := createTestInput(t)

	maxTests := 30
	idCeiling := 10

	// We create a program bitmap to keep track of programs that have been stored /
	// deleted from the context store. We compare the bitmap to the store at the end
	// of the test to verify state correctness.
	programBitmap := make([]bool, idCeiling)

	// just a random test program...
	testProgram := NewProgram("", "", addrs[0], addrs[1], 0)

	rand.Seed(int64(time.Now().Nanosecond()))
	numTests := rand.Int() % maxTests
	for i := 0; i < numTests; i++ {
		programID := uint64(rand.Int63() % int64(idCeiling))
		action := rand.Int() % 2
		if action == 0 {
			programBitmap[programID] = true
			input.budgetKeeper.SetProgram(input.ctx, programID, testProgram)
		} else {
			programBitmap[programID] = false
			input.budgetKeeper.DeleteProgram(input.ctx, programID)
		}
	}

	// Count live programs in the bitmap
	expectedLivePrgmCount := 0
	for programID, live := range programBitmap {
		if live {
			expectedLivePrgmCount++

			// Make sure bitmap positives are also reflected in the store
			_, err := input.budgetKeeper.GetProgram(input.ctx, uint64(programID))
			require.Nil(t, err)
		}
	}

	actualLivePrgmCount := 0
	input.budgetKeeper.IteratePrograms(input.ctx, false, func(programID uint64, program Program) (stop bool) {
		require.True(t, programBitmap[programID])
		actualLivePrgmCount++
		return false
	})

	// Count of live programs should match in the context store and bitmap
	require.Equal(t, expectedLivePrgmCount, actualLivePrgmCount)
}

func TestKeeperVote(t *testing.T) {
	input := createTestInput(t)

	maxTests := 30
	idCeiling := 10
	voterCeiling := 10

	voters := make([]sdk.AccAddress, voterCeiling)
	for i := 0; i < voterCeiling; i++ {
		voters = append(voters, sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()))
	}

	// We create a vote bitmap to keep track of programs that have been stored /
	// deleted from the context store. We compare the bitmap to the store at the end
	// of the test to verify state correctness.
	voteBitmap := make([]map[string]bool, idCeiling)
	for j := 0; j < idCeiling; j++ {
		voteBitmap[j] = make(map[string]bool)
	}

	rand.Seed(int64(time.Now().Nanosecond()))
	numTests := rand.Int() % maxTests
	for i := 0; i < numTests; i++ {
		programID := uint64(rand.Int63() % int64(idCeiling))
		voterIndex := uint64(rand.Int() % voterCeiling)
		voterAddress := voters[voterIndex]

		action := rand.Int() % 2
		if action == 0 {
			voteBitmap[programID][voterAddress.String()] = true
			input.budgetKeeper.AddVote(input.ctx, programID, voterAddress, true)
		} else {
			voteBitmap[programID][voterAddress.String()] = false
			input.budgetKeeper.DeleteVote(input.ctx, programID, voterAddress)
		}
	}

	// Count live programs in the bitmap
	expectedLiveVoteCount := 0
	for programID, votes := range voteBitmap {
		for voterAddrStr, live := range votes {
			voterAddr, err := sdk.AccAddressFromBech32(voterAddrStr)
			require.Nil(t, err)

			if live {
				expectedLiveVoteCount++

				_, err2 := input.budgetKeeper.GetVote(input.ctx, uint64(programID), voterAddr)
				require.Nil(t, err2)
			}
		}
	}

	// Match live programs in the store
	actualLiveVoteCount := 0
	input.budgetKeeper.IterateVotes(input.ctx,
		func(programID uint64, voterAddress sdk.AccAddress, option bool) (stop bool) {
			require.True(t, voteBitmap[programID][voterAddress.String()])
			actualLiveVoteCount++
			return false
		})

	// Count of live programs should match in the context store and bitmap
	require.Equal(t, expectedLiveVoteCount, actualLiveVoteCount)
}

func TestKeeperCandidateQueue(t *testing.T) {
	input := createTestInput(t)

	// Insert a program in the queue
	input.budgetKeeper.CandQueueInsert(input.ctx, 0, 0)

	// Check if it exists; it should.
	exists := input.budgetKeeper.CandQueueHas(input.ctx, 0, 0)
	require.True(t, exists)

	// Not with a different endblock num though.
	exists = input.budgetKeeper.CandQueueHas(input.ctx, 1, 0)
	require.False(t, exists)

	// Nor with a different programID.
	exists = input.budgetKeeper.CandQueueHas(input.ctx, 0, 1)
	require.False(t, exists)

	// delete works too.
	input.budgetKeeper.CandQueueRemove(input.ctx, 0, 0)
	exists = input.budgetKeeper.CandQueueHas(input.ctx, 0, 0)
	require.False(t, exists)

	// test iterator
	numTests := 30
	counter := 0
	for i := 0; i < numTests; i++ {
		programID := input.budgetKeeper.NewProgramID(input.ctx)
		input.budgetKeeper.CandQueueInsert(input.ctx, int64(i), programID)
		input.ctx = input.ctx.WithBlockHeight(int64(i))
		input.budgetKeeper.CandQueueIterateExpired(input.ctx, input.ctx.BlockHeight(),
			func(programID uint64) (stop bool) {
				counter++

				input.budgetKeeper.CandQueueRemove(input.ctx, int64(i), uint64(i))
				return false
			})
	}

	require.Equal(t, numTests, counter)
}
