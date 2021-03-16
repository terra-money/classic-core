package types

import (
	"testing"

	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/types"

	core "github.com/terra-project/core/types"
)

func TestSendAuthorization(t *testing.T) {
	addrs := []sdk.AccAddress{
		sdk.AccAddress([]byte("addr1_______________")),
		sdk.AccAddress([]byte("addr2_______________")),
	}

	halfCoins := sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, sdk.NewInt(50)))
	coins := sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, sdk.NewInt(100)))
	sendAuthorization := NewSendAuthorization(coins)

	allow, _, _ := sendAuthorization.Accept(testdata.NewTestMsg(addrs[0]), tmproto.Header{})
	require.False(t, allow)

	allow, _, deleted := sendAuthorization.Accept(types.NewMsgSend(addrs[0], addrs[1], coins), tmproto.Header{})
	require.True(t, allow)
	require.True(t, deleted)

	allow, updated, deleted := sendAuthorization.Accept(types.NewMsgSend(addrs[0], addrs[1], halfCoins), tmproto.Header{})
	require.True(t, allow)
	require.False(t, deleted)
	require.Equal(t, updated, NewSendAuthorization(halfCoins))

	allow, _, _ = sendAuthorization.Accept(types.NewMsgSend(addrs[0], addrs[1], coins.Add(coins...)), tmproto.Header{})
	require.False(t, allow)
}
