package types

import (
	"testing"

	"github.com/stretchr/testify/require"
	core "github.com/terra-project/core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestGenesisValidation(t *testing.T) {
	genState := DefaultGenesisState()
	require.NoError(t, ValidateGenesis(genState))

	// Error - tax_rate range error
	genState.TaxRate = sdk.NewDec(-1)
	require.Error(t, ValidateGenesis(genState))

	// Valid
	genState.TaxRate = sdk.NewDecWithPrec(1, 2)
	require.NoError(t, ValidateGenesis(genState))

	// Error - reward_weight range error
	genState.RewardWeight = sdk.NewDec(-1)
	require.Error(t, ValidateGenesis(genState))

	// Valid
	genState.RewardWeight = sdk.NewDecWithPrec(5, 2)
	require.NoError(t, ValidateGenesis(genState))

	// Error - cumulated_height range error
	genState.CumulatedHeight = -1
	require.Error(t, ValidateGenesis(genState))

	// Error - cumulated_height indicates 2 epoch, but stored TRs is smaller than 2
	genState.CumulatedHeight = 2 * core.BlocksPerWeek
	require.Error(t, ValidateGenesis(genState))

	dummyDec := sdk.NewDec(10)
	dummyInt := sdk.NewInt(10)

	// Error - cumulated_height indicates 2 epoch, but stored SRs is smaller than 2
	genState.TRs = []sdk.Dec{dummyDec, dummyDec}
	require.Error(t, ValidateGenesis(genState))

	// Error - cumulated_height indicates 2 epoch, but stored TSLs is smaller than 2
	genState.SRs = []sdk.Dec{dummyDec, dummyDec}
	require.Error(t, ValidateGenesis(genState))

	// Valid
	genState.TSLs = []sdk.Int{dummyInt, dummyInt}
	require.NoError(t, ValidateGenesis(genState))
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
