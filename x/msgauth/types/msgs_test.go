package types

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgGrantAuthorization(t *testing.T) {
	addrs := []sdk.AccAddress{
		sdk.AccAddress([]byte("addr1_______________")),
		sdk.AccAddress([]byte("addr2_______________")),
	}

	RegisterGrantableMsgType("swap")

	tests := []struct {
		granter       sdk.AccAddress
		grantee       sdk.AccAddress
		authorization AuthorizationI
		period        time.Duration
		expectedErr   string
	}{
		{addrs[0], addrs[1], NewGenericAuthorization("swap"), 100, ""},
		{addrs[0], addrs[0], NewGenericAuthorization("swap"), 100, "can not be grantee == granter: invalid request"},
		{sdk.AccAddress{}, addrs[1], NewGenericAuthorization("swap"), 100, "Invalid granter address (empty address string is not allowed): invalid address"},
		{addrs[0], sdk.AccAddress{}, NewGenericAuthorization("swap"), 100, "Invalid grantee address (empty address string is not allowed): invalid address"},
		{addrs[0], addrs[1], NewGenericAuthorization("swap"), 0, "period of authorization should be positive time duration"},
	}

	for _, tc := range tests {
		msg, err := NewMsgGrantAuthorization(tc.granter, tc.grantee, tc.authorization, tc.period)
		require.NoError(t, err)

		if tc.expectedErr == "" {
			require.NoError(t, msg.ValidateBasic())
		} else {
			require.EqualError(t, msg.ValidateBasic(), tc.expectedErr)
		}
	}
}

func TestMsgRevokeAuthorization(t *testing.T) {
	addrs := []sdk.AccAddress{
		sdk.AccAddress([]byte("addr1_______________")),
		sdk.AccAddress([]byte("addr2_______________")),
	}

	tests := []struct {
		granter     sdk.AccAddress
		grantee     sdk.AccAddress
		expectedErr string
	}{
		{addrs[0], addrs[1], ""},
		{sdk.AccAddress{}, addrs[1], "Invalid granter address (empty address string is not allowed): invalid address"},
		{addrs[0], sdk.AccAddress{}, "Invalid grantee address (empty address string is not allowed): invalid address"},
	}

	for _, tc := range tests {
		msg := NewMsgRevokeAuthorization(tc.granter, tc.grantee, "")
		if tc.expectedErr == "" {
			require.Nil(t, msg.ValidateBasic())
		} else {
			require.EqualError(t, msg.ValidateBasic(), tc.expectedErr)
		}
	}
}

func TestMsgExecAuthorization(t *testing.T) {
	addrs := []sdk.AccAddress{
		sdk.AccAddress([]byte("addr2_______________")),
	}

	tests := []struct {
		grantee     sdk.AccAddress
		msgs        []sdk.Msg
		expectedErr string
	}{
		{addrs[0], []sdk.Msg{testdata.NewTestMsg(addrs[0])}, ""},
		{sdk.AccAddress{}, []sdk.Msg{testdata.NewTestMsg(addrs[0])}, "Invalid grantee address (empty address string is not allowed): invalid address"},
		{addrs[0], []sdk.Msg{}, "cannot execute empty msgs: invalid request"},
	}

	for _, tc := range tests {
		msg, err := NewMsgExecAuthorized(tc.grantee, tc.msgs)
		require.NoError(t, err)

		if tc.expectedErr == "" {
			require.NoError(t, msg.ValidateBasic())
		} else {
			require.EqualError(t, msg.ValidateBasic(), tc.expectedErr)
		}
	}
}
