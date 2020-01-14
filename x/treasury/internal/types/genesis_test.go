package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestGenesisValidation(t *testing.T) {
	genState := DefaultGenesisState()
	require.NoError(t, ValidateGenesis(genState))

	// Error - tax-rate range error
	genState.TaxRate = sdk.NewDec(-1)
	require.Error(t, ValidateGenesis(genState))

	// Valid
	genState.TaxRate = sdk.NewDecWithPrec(1, 2)
	require.NoError(t, ValidateGenesis(genState))

	// Error - reward-weight range error
	genState.RewardWeight = sdk.NewDec(-1)
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
