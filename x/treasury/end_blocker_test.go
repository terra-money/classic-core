package treasury

import (
	"math/rand"
	"terra/types"
	"terra/types/assets"
	"terra/types/util"
	"terra/x/treasury/tags"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestEndBlockerTiming(t *testing.T) {
	input := createTestInput(t)
	input = reset(input)

	// First endblocker should fail
	tTags := EndBlocker(input.ctx, input.treasuryKeeper)
	require.True(t, len(tTags.ToKVPairs()) == 0)

	// Subsequent endblocker should settle, but NOT update policy
	params := input.treasuryKeeper.GetParams(input.ctx)
	for i := int64(1); i < params.EpochProbation.Int64(); i++ {
		if i%params.EpochShort.Int64() == 0 {
			input.ctx = input.ctx.WithBlockHeight(i * util.GetBlocksPerEpoch())
			input.mintKeeper.AddSeigniorage(input.ctx, lunaAmt)

			tTags := EndBlocker(input.ctx, input.treasuryKeeper)

			require.Equal(t, 4, len(tTags))
		}
	}

	// After probationary period, we should also be updating policy variables
	for i := params.EpochProbation.Int64(); i < params.EpochProbation.Int64()+12; i++ {
		if i%params.EpochShort.Int64() == 0 {
			input.ctx = input.ctx.WithBlockHeight(i * util.GetBlocksPerEpoch())
			input.mintKeeper.AddSeigniorage(input.ctx, lunaAmt)

			tTags := EndBlocker(input.ctx, input.treasuryKeeper)

			require.Equal(t, tTags.ToKVPairs()[4].GetValue(), []byte(tags.ActionPolicyUpdate))
		}
	}
}

func randomMacroVariables() (tax, seigniorage, lunaIssuance sdk.Int) {
	rand.Seed(int64(time.Now().Nanosecond()))

	tax = sdk.NewInt(rand.Int63())
	seigniorage = sdk.NewInt(rand.Int63())
	lunaIssuance = sdk.NewInt(rand.Int63() % 1000)

	return
}

func reset(input testInput) testInput {

	// Set blocknum back to 0
	input.ctx = input.ctx.WithBlockHeight(0)

	// Reset oracle price
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.SDRDenom, sdk.NewDec(1))

	// Reset genesis
	InitGenesis(input.ctx, input.treasuryKeeper, DefaultGenesisState())

	// Give everyone some luna
	for _, addr := range addrs {
		input.mintKeeper.Mint(input.ctx, addr, sdk.NewCoin(assets.LunaDenom, lunaAmt))
	}

	return input
}

// updatePolicy updates
func updatePolicy(input testInput, startIndex int,
	taxRevenues, seigniorageRevenues []sdk.Int) (taxRate, rewardWeight sdk.Dec, err sdk.Error, ctx sdk.Context) {

	if len(taxRevenues) != len(seigniorageRevenues) {
		err = sdk.ErrInternal("lengths of inputs should be the same")
	}

	params := input.treasuryKeeper.GetParams(input.ctx)
	blocksPerEpoch := util.GetBlocksPerEpoch()

	for i := 0; i < len(taxRevenues); i++ {
		input.ctx = input.ctx.WithBlockHeight(params.EpochShort.Int64() * int64(i+startIndex) * blocksPerEpoch)

		taxRevenue := taxRevenues[i]
		input.treasuryKeeper.RecordTaxProceeds(input.ctx, sdk.Coins{sdk.NewCoin(assets.SDRDenom, taxRevenue)})

		seigniorageRevenue := seigniorageRevenues[i]
		input.mintKeeper.AddSeigniorage(input.ctx, seigniorageRevenue)

		// Call endblocker
		EndBlocker(input.ctx, input.treasuryKeeper)
	}

	taxRate = input.treasuryKeeper.GetTaxRate(input.ctx, util.GetEpoch(input.ctx))
	rewardWeight = input.treasuryKeeper.GetRewardWeight(input.ctx, util.GetEpoch(input.ctx))
	ctx = input.ctx

	return
}

func generateRandomMacroInputs() (taxRevenues, seigniorageRevenues []sdk.Int) {
	rand.Seed(int64(time.Now().Nanosecond()))

	taxRevenues = []sdk.Int{}
	seigniorageRevenues = []sdk.Int{}

	numPeriods := rand.Int63() % 30 // bound to less than 30 periods
	for i := 0; i < int(numPeriods); i++ {
		taxRevenues = append(taxRevenues, sdk.NewInt(rand.Int63()))
		seigniorageRevenues = append(seigniorageRevenues, sdk.NewInt(rand.Int63()))
	}

	return
}

func TestEndBlockerUpdatePolicy(t *testing.T) {
	input := createTestInput(t)
	input = reset(input)

	taxRevenues, seigniorageRevenues := generateRandomMacroInputs()
	newTaxRate, newSeigniorageWeight, err, ctx := updatePolicy(input, 1, taxRevenues, seigniorageRevenues)
	require.Nil(t, err)

	input.ctx = ctx
	taxRate := input.treasuryKeeper.GetTaxRate(input.ctx, util.GetEpoch(input.ctx))
	rewardWeight := input.treasuryKeeper.GetRewardWeight(input.ctx, util.GetEpoch(input.ctx))

	require.Equal(t, taxRate, newTaxRate)
	require.Equal(t, rewardWeight, newSeigniorageWeight)
}

func TestEndBlockerSettleClaims(t *testing.T) {
	input := createTestInput(t)
	input = reset(input)

	tests := []struct {
		claims      []types.Claim
		seigniorage sdk.Int
		sdrRewards  []int64
	}{
		// Test 1: no claims
		{[]types.Claim{}, sdk.NewInt(1000), []int64{0, 0, 0}},

		// Test 2: two claims of the same class
		{[]types.Claim{
			types.NewClaim(types.OracleClaimClass, sdk.NewInt(10), addrs[0]),
			types.NewClaim(types.OracleClaimClass, sdk.NewInt(90), addrs[1]),
		}, sdk.NewInt(1000), []int64{10, 90, 0}},

		// Test 3: similar to case 2, except different class
		{[]types.Claim{
			types.NewClaim(types.BudgetClaimClass, sdk.NewInt(10), addrs[0]),
			types.NewClaim(types.BudgetClaimClass, sdk.NewInt(90), addrs[1]),
		}, sdk.NewInt(1000), []int64{90, 810, 0}},

		// Test 4: Many claims of different classes
		{[]types.Claim{
			types.NewClaim(types.OracleClaimClass, sdk.NewInt(10), addrs[0]),
			types.NewClaim(types.BudgetClaimClass, sdk.NewInt(10), addrs[0]),
			types.NewClaim(types.OracleClaimClass, sdk.NewInt(10), addrs[1]),
			types.NewClaim(types.OracleClaimClass, sdk.NewInt(80), addrs[2]),
			types.NewClaim(types.BudgetClaimClass, sdk.NewInt(90), addrs[2]),
		}, sdk.NewInt(1000), []int64{100, 10, 890}},
	}

	params := input.treasuryKeeper.GetParams(input.ctx)
	blocksPerEpoch := util.GetBlocksPerEpoch()

	for i, tc := range tests {

		// Advance blockcount
		input.ctx = input.ctx.WithBlockHeight(params.EpochShort.Int64() * blocksPerEpoch * int64(i))

		// clear SDR balances for testing; keep luna for policy update safety
		for _, addr := range addrs {
			input.bankKeeper.SetCoins(input.ctx, addr, sdk.Coins{sdk.NewCoin(assets.LunaDenom, lunaAmt)})
		}

		// Reset reward weight
		input.treasuryKeeper.SetRewardWeight(input.ctx, sdk.ZeroDec())
		input.mintKeeper.AddSeigniorage(input.ctx, tc.seigniorage)

		// Call endblocker
		for _, claim := range tc.claims {
			input.treasuryKeeper.AddClaim(input.ctx, claim)
		}
		EndBlocker(input.ctx, input.treasuryKeeper)

		for j, addr := range addrs {
			balance := input.bankKeeper.GetCoins(input.ctx, addr).AmountOf(assets.SDRDenom)
			require.Equal(t, balance, sdk.NewInt(tc.sdrRewards[j]), "test: %v", i)
		}

		counter := 0
		input.treasuryKeeper.IterateClaims(input.ctx, func(claim types.Claim) bool {
			counter++
			return false
		})

		require.Equal(t, 0, counter, "Claims expected to be cleared after treasury update")
	}
}
