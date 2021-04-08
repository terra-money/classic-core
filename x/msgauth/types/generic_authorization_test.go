package types

import (
	"testing"

	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/testutil/testdata"
)

func TestGenericAuthorization(t *testing.T) {
	generic := NewGenericAuthorization("delegate")

	allow, updated, delete := generic.Accept(testdata.NewTestMsg(), tmproto.Header{})
	require.True(t, allow)
	require.Equal(t, NewGenericAuthorization("delegate"), updated)
	require.False(t, delete)
}
