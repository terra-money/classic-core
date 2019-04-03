package budget

import (
	"math/rand"
	"terra/types/mock"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func TestEndBlockerTallyBasic(t *testing.T) {
	input := createTestInput(t)

	// create test program
	testProgram := generateTestProgram(input.ctx, input.budgetKeeper)

	input.budgetKeeper.SetProgram(input.ctx, testProgram.ProgramID, testProgram)

	// Add validators and their votes; to keep things simple, let's assume each validator holds 1 token
	valset := mock.NewMockValSet()
	for _, addr := range addrs {
		valAddr := sdk.ValAddress(addr)
		validator := mock.NewMockValidator(valAddr, sdk.OneInt())
		valset.Validators = append(valset.Validators, validator)

		input.budgetKeeper.AddVote(input.ctx, testProgram.ProgramID, addr, true)
	}
	input.budgetKeeper.valset = valset

	actualVotePower, actualTotalPower := tally(input.ctx, input.budgetKeeper, testProgram.ProgramID)

	// totalPower and votepower should match the number of validators (uniform, single weighted)
	require.Equal(t, actualTotalPower, sdk.NewInt(int64(len(addrs))))
	require.Equal(t, actualVotePower, sdk.NewInt(int64(len(addrs))))
}

func TestEndBlockerTallyRandom(t *testing.T) {
	input := createTestInput(t)

	// create test program

	testProgram := generateTestProgram(input.ctx, input.budgetKeeper)

	input.budgetKeeper.SetProgram(input.ctx, testProgram.ProgramID, testProgram)

	rand.Seed(int64(time.Now().Nanosecond()))
	numValidators := rand.Int() % 100 // cap validator count by a 100

	totalPower := 0
	votePower := 0
	valset := mock.NewMockValSet()
	for i := 0; i < numValidators; i++ {
		valAccAddr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
		valAddr := sdk.ValAddress(valAccAddr)
		valPower := rand.Int() % 10000
		validator := mock.NewMockValidator(valAddr, sdk.NewInt(int64(valPower)))
		valset.Validators = append(valset.Validators, validator)

		totalPower += valPower

		option := rand.Int() % 2
		choice := option == 1
		if choice {
			votePower += valPower
		} else {
			votePower -= valPower
		}

		input.budgetKeeper.AddVote(input.ctx, testProgram.ProgramID, valAccAddr, choice)
	}
	input.budgetKeeper.valset = valset

	actualVotePower, actualTotalPower := tally(input.ctx, input.budgetKeeper, testProgram.ProgramID)

	require.Equal(t, actualTotalPower, sdk.NewInt(int64(totalPower)))
	require.Equal(t, actualVotePower, sdk.NewInt(int64(votePower)))
}

func TestEndBlockerTiming(t *testing.T) {
	input := createTestInput(t)

	// create test program
	testProgram := generateTestProgram(input.ctx, input.budgetKeeper)

	input.budgetKeeper.SetProgram(input.ctx, testProgram.ProgramID, testProgram)

	// Add a vote each from validators
	for _, addr := range addrs {
		input.budgetKeeper.AddVote(input.ctx, testProgram.ProgramID, addr, true)
	}

	// No claims should have been settled yet
	claims, _ := EndBlocker(input.ctx, input.budgetKeeper)
	require.Equal(t, 0, len(claims))

	// Advance block height by voteperiod, and the program should be settled.
	params := input.budgetKeeper.GetParams(input.ctx)
	input.ctx = input.ctx.WithBlockHeight(params.VotePeriod)
	claims, _ = EndBlocker(input.ctx, input.budgetKeeper)

	require.Equal(t, 1, len(claims))
	require.Equal(t, input.budgetKeeper.valset.TotalBondedTokens(input.ctx), claims[0].Weight)
}

func TestEndBlockerPassOrReject(t *testing.T) {
	input := createTestInput(t)

	// add a hundred validators with 1 stakable token each
	valset := mock.NewMockValSet()
	valAddrs := []sdk.AccAddress{}
	for i := 0; i < 100; i++ {
		valAccAddr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
		validator := mock.NewMockValidator(sdk.ValAddress(valAccAddr), sdk.OneInt())
		valset.Validators = append(valset.Validators, validator)
		valAddrs = append(valAddrs, valAccAddr)
	}
	input.budgetKeeper.valset = valset

	// Compute minimum validator support
	activeThreshold := input.budgetKeeper.GetParams(input.ctx).ActiveThreshold
	minNumTokensToPass := activeThreshold.MulInt(sdk.NewInt(100)).TruncateInt()

	// create test program
	testProgram := generateTestProgram(input.ctx, input.budgetKeeper)
	input.budgetKeeper.SetProgram(input.ctx, testProgram.ProgramID, testProgram)
	input.budgetKeeper.CandQueueInsert(input.ctx, testProgram.getVotingEndBlock(input.ctx, input.budgetKeeper), testProgram.ProgramID)

	// vote slightly such that the sum falls short of the threshold; tally should fail and program not activated.
	for i := 0; i < int(minNumTokensToPass.Int64())-1; i++ {
		input.budgetKeeper.AddVote(input.ctx, testProgram.ProgramID, valAddrs[i], true)
	}

	params := input.budgetKeeper.GetParams(input.ctx)
	input.ctx = input.ctx.WithBlockHeight(params.VotePeriod)
	EndBlocker(input.ctx, input.budgetKeeper)
	_, err := input.budgetKeeper.GetProgram(input.ctx, testProgram.ProgramID)
	require.NotNil(t, err)

	input.budgetKeeper.DeleteProgram(input.ctx, testProgram.ProgramID)

	// vote above the threshold; the tally should now pass
	testProgram2 := generateTestProgram(input.ctx, input.budgetKeeper)
	input.budgetKeeper.SetProgram(input.ctx, testProgram2.ProgramID, testProgram2)
	input.budgetKeeper.CandQueueInsert(input.ctx, testProgram2.getVotingEndBlock(input.ctx, input.budgetKeeper), testProgram2.ProgramID)

	for i := 0; i < int(minNumTokensToPass.Int64())+1; i++ {
		input.budgetKeeper.AddVote(input.ctx, testProgram2.ProgramID, valAddrs[i], true)
	}

	input.ctx = input.ctx.WithBlockHeight(params.VotePeriod)
	EndBlocker(input.ctx, input.budgetKeeper)
	_, err = input.budgetKeeper.GetProgram(input.ctx, testProgram2.ProgramID)
	require.Nil(t, err)
}
