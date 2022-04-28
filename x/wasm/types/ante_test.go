package types

import (
	"testing"

	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestTxCount(t *testing.T) {
	ctx := sdk.NewContext(nil, tmproto.Header{}, false, nil)
	_, ok := TXCounter(ctx)
	require.False(t, ok)

	ctx = WithTXCounter(ctx, uint32(10))
	counter, ok := TXCounter(ctx)
	require.True(t, ok)
	require.Equal(t, uint32(10), counter)
}
