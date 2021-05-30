package simulation

import (
	"encoding/json"
	"fmt"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/bank/simulation"
	"github.com/cosmos/cosmos-sdk/x/bank/types"

	core "github.com/terra-money/core/types"
)

// RandomGenesisBalances returns a slice of account balances. Each account has
// a balance of simState.InitialStake for sdk.DefaultBondDenom and core.MicroLunaDenom.
func RandomGenesisBalances(simState *module.SimulationState) []types.Balance {
	genesisBalances := []types.Balance{}

	for _, acc := range simState.Accounts {
		genesisBalances = append(genesisBalances, types.Balance{
			Address: acc.Address.String(),
			Coins: sdk.NewCoins(
				sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(simState.InitialStake)),
				sdk.NewCoin(core.MicroLunaDenom, sdk.NewInt(simState.InitialStake)),
			),
		})
	}

	return genesisBalances
}

// RandomizedGenState generates a random GenesisState for bank
func RandomizedGenState(simState *module.SimulationState) {
	var sendEnabledParams types.SendEnabledParams
	simState.AppParams.GetOrGenerate(
		simState.Cdc, string(types.KeySendEnabled), &sendEnabledParams, simState.Rand,
		func(r *rand.Rand) { sendEnabledParams = simulation.RandomGenesisSendParams(r) },
	)

	var defaultSendEnabledParam bool
	simState.AppParams.GetOrGenerate(
		simState.Cdc, string(types.KeyDefaultSendEnabled), &defaultSendEnabledParam, simState.Rand,
		func(r *rand.Rand) { defaultSendEnabledParam = simulation.RandomGenesisDefaultSendParam(r) },
	)

	numAccs := int64(len(simState.Accounts))
	totalSupply := sdk.NewInt(simState.InitialStake * (numAccs + simState.NumBonded))
	totalLunaSupply := sdk.NewInt(simState.InitialStake * numAccs)
	supply := sdk.NewCoins(
		sdk.NewCoin(sdk.DefaultBondDenom, totalSupply),
		sdk.NewCoin(core.MicroLunaDenom, totalLunaSupply),
	)

	bankGenesis := types.GenesisState{
		Params: types.Params{
			SendEnabled:        sendEnabledParams,
			DefaultSendEnabled: defaultSendEnabledParam,
		},
		Balances: RandomGenesisBalances(simState),
		Supply:   supply,
	}

	paramsBytes, err := json.MarshalIndent(&bankGenesis.Params, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated bank parameters:\n%s\n", paramsBytes)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&bankGenesis)
}
