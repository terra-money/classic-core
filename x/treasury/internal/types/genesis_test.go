package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestGenesisValidation(t *testing.T) {
	genState := DefaultGenesisState()
	require.NoError(t, ValidateGenesis(genState))

	genState.TaxRate = sdk.NewDec(-1)
	require.Error(t, ValidateGenesis(genState))

	genState.TaxRate = sdk.NewDecWithPrec(5, 2)
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
