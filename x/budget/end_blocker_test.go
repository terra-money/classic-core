package budget

import (
	"math/rand"
	"testing"
	"time"

	"github.com/terra-project/core/types/assets"
	"github.com/terra-project/core/types/mock"
	"github.com/terra-project/core/types/util"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func TestEndBlockerTallyBasic(t *testing.T) {
	input := createTestInput(t)

	// create test program
	testProgram := generateTestProgram(input.ctx, input.budgetKeeper)

	input.budgetKeeper.StoreProgram(input.ctx, testProgram)

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

	input.budgetKeeper.StoreProgram(input.ctx, testProgram)

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

	input.budgetKeeper.StoreProgram(input.ctx, testProgram)

	// Add a vote each from validators
	for _, addr := range addrs {
		input.budgetKeeper.AddVote(input.ctx, testProgram.ProgramID, addr, true)
	}

	// No claims should have been settled yet
	EndBlocker(input.ctx, input.budgetKeeper)

	claimCount := countClaimPool(input.ctx, input.budgetKeeper)
	require.Equal(t, 0, claimCount)

	// Advance block height by voteperiod - 1, and the program should be settled.
	params := input.budgetKeeper.GetParams(input.ctx)
	input.ctx = input.ctx.WithBlockHeight(params.VotePeriod - 1)
	EndBlocker(input.ctx, input.budgetKeeper)

	claimCount = countClaimPool(input.ctx, input.budgetKeeper)
	require.Equal(t, 1, claimCount)

	input.budgetKeeper.iterateClaimPool(input.ctx, func(recipient sdk.AccAddress, weight sdk.Int) (stop bool) {
		require.Equal(t, input.budgetKeeper.valset.TotalBondedTokens(input.ctx), weight)
		return true
	})

}

func TestEndBlockerClaimDistribution(t *testing.T) {
	input := createTestInput(t)

	// create test program
	testProgram := generateTestProgram(input.ctx, input.budgetKeeper)

	input.budgetKeeper.StoreProgram(input.ctx, testProgram)

	// Add a vote each from validators
	for _, addr := range addrs {
		input.budgetKeeper.AddVote(input.ctx, testProgram.ProgramID, addr, true)
	}

	// No claims should have been settled yet
	EndBlocker(input.ctx, input.budgetKeeper)

	claimCount := countClaimPool(input.ctx, input.budgetKeeper)
	require.Equal(t, 0, claimCount)

	// Advance block height by voteperiod - 1, and the program should be settled.
	params := input.budgetKeeper.GetParams(input.ctx)
	input.ctx = input.ctx.WithBlockHeight(params.VotePeriod - 1)
	EndBlocker(input.ctx, input.budgetKeeper)

	claimCount = countClaimPool(input.ctx, input.budgetKeeper)
	require.Equal(t, 1, claimCount)

	input.mintKeeper.Mint(input.ctx, addrs[0], sdk.NewCoin(assets.MicroLunaDenom, sdk.NewInt(1000)))

	// after 5 week, distribution date reach
	input.ctx = input.ctx.WithBlockHeight(util.BlocksPerEpoch*5 - 1)
	input.treasuryKeeper.SetRewardWeight(input.ctx, sdk.NewDecWithPrec(1, 1))
	EndBlocker(input.ctx, input.budgetKeeper)

	claimCount = countClaimPool(input.ctx, input.budgetKeeper)
	require.Equal(t, 0, claimCount)
}

func TestEndBlockerLegacy(t *testing.T) {
	input := createTestInput(t)

	defaultBudgetParams := DefaultParams()
	defaultBudgetParams.VotePeriod = 1
	input.budgetKeeper.SetParams(input.ctx, defaultBudgetParams)

	ctx := input.ctx.WithBlockHeight(1)

	// Create test program
	testProgram := generateTestProgram(ctx, input.budgetKeeper)

	input.budgetKeeper.StoreProgram(ctx, testProgram)

	// Add a vote each from validators
	for _, addr := range addrs {
		input.budgetKeeper.AddVote(ctx, testProgram.ProgramID, addr, true)
	}

	// Claims should have been settled
	EndBlocker(ctx, input.budgetKeeper)
	claimCount := countClaimPool(input.ctx, input.budgetKeeper)
	require.Equal(t, 1, claimCount)

	ctx = input.ctx.WithBlockHeight(2)

	for _, addr := range addrs {
		input.budgetKeeper.AddVote(ctx, testProgram.ProgramID, addr, false)
	}

	// Program should be legacy
	EndBlocker(ctx, input.budgetKeeper)
	_, err := input.budgetKeeper.GetProgram(ctx, testProgram.ProgramID)
	require.Error(t, err)
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
	input.budgetKeeper.StoreProgram(input.ctx, testProgram)
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
	input.budgetKeeper.StoreProgram(input.ctx, testProgram2)
	input.budgetKeeper.CandQueueInsert(input.ctx, testProgram2.getVotingEndBlock(input.ctx, input.budgetKeeper), testProgram2.ProgramID)

	for i := 0; i < int(minNumTokensToPass.Int64())+1; i++ {
		input.budgetKeeper.AddVote(input.ctx, testProgram2.ProgramID, valAddrs[i], true)
	}

	input.ctx = input.ctx.WithBlockHeight(params.VotePeriod)
	EndBlocker(input.ctx, input.budgetKeeper)
	_, err = input.budgetKeeper.GetProgram(input.ctx, testProgram2.ProgramID)
	require.Nil(t, err)
}

func countClaimPool(ctx sdk.Context, keeper Keeper) (claimCount int) {
	keeper.iterateClaimPool(ctx, func(recipient sdk.AccAddress, weight sdk.Int) (stop bool) {
		claimCount++
		return false
	})

	return claimCount
}
