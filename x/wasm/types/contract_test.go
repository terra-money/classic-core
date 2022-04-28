package types

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestNewEnv(t *testing.T) {
	ctx := sdk.NewContext(nil, tmproto.Header{
		Height: 100,
		Time:   time.Now(),
	}, false, nil)

	require.NotPanics(t, func() {
		_ = NewEnv(ctx, sdk.AccAddress{})
		_ = NewEnv(WithTXCounter(ctx, 100), sdk.AccAddress{})
	})

	require.Panics(t, func() {
		_ = NewEnv(ctx.WithBlockHeight(-1), sdk.AccAddress{})
	})

	require.Panics(t, func() {
		_ = NewEnv(ctx.WithBlockTime(time.Unix(0, 0)), sdk.AccAddress{})
	})
}
