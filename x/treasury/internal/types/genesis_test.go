package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestGenesisValidation(t *testing.T) {
	genState := DefaultGenesisState()
	require.NoError(t, ValidateGenesis(genState))

	// Error - at least one tax-rate should be given
	genState.TaxRates = []sdk.Dec{}
	require.Error(t, ValidateGenesis(genState))

	// Error - tax-rate range error
	genState.TaxRates = []sdk.Dec{sdk.NewDec(-1)}
	require.Error(t, ValidateGenesis(genState))

	// Valid
	genState.TaxRates = []sdk.Dec{sdk.NewDecWithPrec(1, 2)}
	require.NoError(t, ValidateGenesis(genState))

	// Error - at least one reward-weight should be given
	genState.RewardWeights = []sdk.Dec{}
	require.Error(t, ValidateGenesis(genState))

	// Error - reward-weight range error
	genState.RewardWeights = []sdk.Dec{sdk.NewDec(-1)}
	require.Error(t, ValidateGenesis(genState))
}

func TestGenesisEqual(t *testing.T) {
	genState1 := DefaultGenesisState()
	genState2 := DefaultGenesisState()

	require.True(t, genState1.Equal(genState2))
}

func TestGenesisEmpty(t *testing.T) {
	genState := GenesisState{}
	require.True(t, genState.IsEmpty())
}
