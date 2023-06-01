package simulation

import (
	"encoding/json"
	"fmt"
	"math/rand" //#nosec G404

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/classic-terra/core/v2/x/feeshare/types"
)

// Simulation parameter constants
const (
	enableFeeShareKey  = "enable_fee_share"
	developerSharesKey = "developer_shares"
	allowedDenomsKey   = "allowed_denoms"
)

// GenEnableFeeShare
func GenEnableFeeShare(r *rand.Rand) bool {
	return r.Intn(2) == 1
}

// GenDeveloperShares
func GenDeveloperShares(r *rand.Rand) sdk.Dec {
	// Generate a random float number between 0 and 1
	randNum := r.Float64()

	// precision of 10
	return sdk.MustNewDecFromStr(fmt.Sprintf("%.10f", randNum))
}

// GenAllowedDenoms
func GenAllowedDenoms(r *rand.Rand) []string {
	return []string(nil)
}

// RandomizedGenState generates a random GenesisState for gov
func RandomizedGenState(simState *module.SimulationState) {
	var enableFeeShare bool
	simState.AppParams.GetOrGenerate(
		simState.Cdc, enableFeeShareKey, &enableFeeShare, simState.Rand,
		func(r *rand.Rand) { enableFeeShare = GenEnableFeeShare(r) },
	)

	var developerShares sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, developerSharesKey, &developerShares, simState.Rand,
		func(r *rand.Rand) { developerShares = GenDeveloperShares(r) },
	)

	var allowedDenoms []string
	simState.AppParams.GetOrGenerate(
		simState.Cdc, allowedDenomsKey, &allowedDenoms, simState.Rand,
		func(r *rand.Rand) { allowedDenoms = GenAllowedDenoms(r) },
	)

	feeshareGenesis := types.NewGenesisState(
		types.Params{
			EnableFeeShare:  enableFeeShare,
			DeveloperShares: developerShares,
			AllowedDenoms:   allowedDenoms,
		},
		[]types.FeeShare(nil),
	)

	bz, err := json.MarshalIndent(&feeshareGenesis.Params, "", " ")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Selected randomly generated feeshare parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(feeshareGenesis)
}
