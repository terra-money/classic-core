package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParams(t *testing.T) {
	params := DefaultParams()
	require.NoError(t, params.Validate())

	params = DefaultParams()
	params.MaxContractGas = EnforcedMaxContractGas + 1
	require.Error(t, params.Validate())

	params = DefaultParams()
	params.MaxContractMsgSize = EnforcedMaxContractMsgSize + 1
	require.Error(t, params.Validate())

	params = DefaultParams()
	params.MaxContractSize = EnforcedMaxContractSize + 1
	require.Error(t, params.Validate())
}
