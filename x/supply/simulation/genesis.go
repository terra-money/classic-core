package simulation

// DONTCOVER

import (
	"fmt"
	core "github.com/terra-money/core/types"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

// RandomizedGenState generates a random GenesisState for supply
func RandomizedGenState(simState *module.SimulationState) {
	numAccs := int64(len(simState.Accounts))
	totalSupply := sdk.NewInt(simState.InitialStake * (numAccs + simState.NumBonded))
	totalLunaSupply := sdk.NewInt(simState.InitialStake * numAccs)
	supplyGenesis := supply.NewGenesisState(sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, totalSupply), sdk.NewCoin(core.MicroLunaDenom, totalLunaSupply)))

	fmt.Printf("Generated supply parameters:\n%s\n", codec.MustMarshalJSONIndent(simState.Cdc, supplyGenesis))
	simState.GenState[supply.ModuleName] = simState.Cdc.MustMarshalJSON(supplyGenesis)
}
