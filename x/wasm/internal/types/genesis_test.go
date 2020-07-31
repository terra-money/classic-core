package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenesisValidation(t *testing.T) {
	genState := DefaultGenesisState()
	require.NoError(t, ValidateGenesis(genState))

	genState.Params.MaxContractSize = EnforcedMaxContractSize + 1
	require.Error(t, ValidateGenesis(genState))

	genState = DefaultGenesisState()
	genState.Params.MaxContractGas = EnforcedMaxContractGas + 1
	require.Error(t, ValidateGenesis(genState))

	genState = DefaultGenesisState()
	genState.Params.MaxContractMsgSize = EnforcedMaxContractMsgSize + 1
	require.Error(t, ValidateGenesis(genState))

	genState = DefaultGenesisState()
	genState.Codes = []Code{{}, {}}
	genState.LastCodeID = 2
	require.NoError(t, ValidateGenesis(genState))

	genState.LastCodeID = 1
	require.Error(t, ValidateGenesis(genState))

	genState = DefaultGenesisState()
	genState.Contracts = []Contract{{}, {}}
	genState.LastInstanceID = 2
	require.NoError(t, ValidateGenesis(genState))

	genState.LastInstanceID = 1
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
