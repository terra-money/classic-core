package oracle_test

import (
	"sort"
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/core/x/oracle"
	"github.com/terra-money/core/x/oracle/types"
)

func TestFuzz_Tally(t *testing.T) {
	validators := map[string]int64{}

	f := fuzz.New().NilChance(0).Funcs(
		func(e *sdk.Dec, c fuzz.Continue) {
			*e = sdk.NewDec(c.Int63())
		},
		func(e *map[string]int64, c fuzz.Continue) {
			numValidators := c.Intn(100) + 5

			for i := 0; i < numValidators; i++ {
				(*e)[sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address()).String()] = c.Int63n(100)
			}
		},
		func(e *map[string]types.Claim, c fuzz.Continue) {
			for validator, power := range validators {
				addr, err := sdk.ValAddressFromBech32(validator)
				require.NoError(t, err)
				(*e)[validator] = types.NewClaim(power, 0, 0, addr)
			}
		},
		func(e *types.ExchangeRateBallot, c fuzz.Continue) {

			ballot := types.ExchangeRateBallot{}
			for addr, power := range validators {
				addr, _ := sdk.ValAddressFromBech32(addr)

				var rate sdk.Dec
				c.Fuzz(&rate)

				ballot = append(ballot, types.NewVoteForTally(rate, c.RandString(), addr, power))
			}

			sort.Sort(ballot)

			*e = ballot
		},
	)

	// set random denoms and validators
	f.Fuzz(&validators)

	input, _ := setup(t)

	claimMap := map[string]types.Claim{}
	f.Fuzz(&claimMap)

	ballot := types.ExchangeRateBallot{}
	f.Fuzz(&ballot)

	var rewardBand sdk.Dec
	f.Fuzz(&rewardBand)

	require.NotPanics(t, func() {
		oracle.Tally(input.Ctx, ballot, rewardBand, claimMap)
	})
}

func TestFuzz_PickReferenceTerra(t *testing.T) {
	var denoms []string

	f := fuzz.New().NilChance(0).Funcs(
		func(e *[]string, c fuzz.Continue) {
			numDenoms := c.Intn(100) + 5

			for i := 0; i < numDenoms; i++ {
				*e = append(*e, c.RandString())
			}
		},
		func(e *sdk.Dec, c fuzz.Continue) {
			*e = sdk.NewDec(c.Int63())
		},
		func(e *map[string]sdk.Dec, c fuzz.Continue) {
			for _, denom := range denoms {
				var rate sdk.Dec
				c.Fuzz(&rate)

				(*e)[denom] = rate
			}
		},
		func(e *map[string]int64, c fuzz.Continue) {
			numValidator := c.Intn(100) + 5
			for i := 0; i < numValidator; i++ {
				(*e)[sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address()).String()] = int64(c.Intn(100) + 1)
			}
		},
		func(e *map[string]types.ExchangeRateBallot, c fuzz.Continue) {
			validators := map[string]int64{}
			c.Fuzz(&validators)

			for _, denom := range denoms {
				ballot := types.ExchangeRateBallot{}

				for addr, power := range validators {
					addr, _ := sdk.ValAddressFromBech32(addr)

					var rate sdk.Dec
					c.Fuzz(&rate)

					ballot = append(ballot, types.NewVoteForTally(rate, denom, addr, power))
				}

				sort.Sort(ballot)
				(*e)[denom] = ballot
			}
		},
	)

	// set random denoms
	f.Fuzz(&denoms)

	input, _ := setup(t)

	voteTargets := map[string]sdk.Dec{}
	f.Fuzz(&voteTargets)

	voteMap := map[string]types.ExchangeRateBallot{}
	f.Fuzz(&voteMap)

	require.NotPanics(t, func() {
		oracle.PickReferenceTerra(input.Ctx, input.OracleKeeper, voteTargets, voteMap)
	})
}
