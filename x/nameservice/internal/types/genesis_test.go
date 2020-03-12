package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateGenesis(t *testing.T) {
	genesisState := DefaultGenesisState()
	require.NoError(t, ValidateGenesis(genesisState))
}
