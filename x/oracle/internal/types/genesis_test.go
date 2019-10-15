package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenesisValidation(t *testing.T) {
	genState := DefaultGenesisState()
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
