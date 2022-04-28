package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestMsgStoreCode(t *testing.T) {
	addrs := []sdk.AccAddress{
		sdk.AccAddress([]byte("addr1_______________")),
	}

	tests := []struct {
		sender       sdk.AccAddress
		wasmByteCode []byte
		expectPass   bool
	}{
		{addrs[0], []byte{}, false},
		{sdk.AccAddress{}, []byte{1, 2, 3}, false},
		{addrs[0], make([]byte, EnforcedMaxContractSize+1), false},
		{addrs[0], []byte{1, 2, 3}, true},
	}

	for i, tc := range tests {
		msg := NewMsgStoreCode(tc.sender, tc.wasmByteCode)
		require.Equal(t, RouterKey, msg.Route())
		require.Equal(t, TypeMsgStoreCode, msg.Type())
		if !tc.sender.Empty() {
			require.Equal(t, tc.sender, msg.GetSigners()[0])
		} else {
			require.Panics(t, func() {
				msg.GetSigners()
			})
		}

		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgMigrateCode(t *testing.T) {
	addrs := []sdk.AccAddress{
		sdk.AccAddress([]byte("addr1_______________")),
	}

	tests := []struct {
		codeID       uint64
		sender       sdk.AccAddress
		wasmByteCode []byte
		expectPass   bool
	}{
		{1, addrs[0], []byte{}, false},
		{2, sdk.AccAddress{}, []byte{1, 2, 3}, false},
		{3, addrs[0], make([]byte, EnforcedMaxContractSize+1), false},
		{1, addrs[0], []byte{1, 2, 3}, true},
	}

	for i, tc := range tests {
		msg := NewMsgMigrateCode(tc.codeID, tc.sender, tc.wasmByteCode)
		require.Equal(t, RouterKey, msg.Route())
		require.Equal(t, TypeMsgMigrateCode, msg.Type())
		if !tc.sender.Empty() {
			require.Equal(t, tc.sender, msg.GetSigners()[0])
		} else {
			require.Panics(t, func() {
				msg.GetSigners()
			})
		}

		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgInstantiateCode(t *testing.T) {
	addrs := []sdk.AccAddress{
		sdk.AccAddress([]byte("addr1_______________")),
	}

	tests := []struct {
		creator    sdk.AccAddress
		admin      sdk.AccAddress
		codeID     uint64
		initMsg    []byte
		initCoins  sdk.Coins
		expectPass bool
	}{
		{sdk.AccAddress{}, sdk.AccAddress{}, 0, []byte("{}"), sdk.Coins{}, false},
		{addrs[0], sdk.AccAddress{}, 0, make([]byte, EnforcedMaxContractMsgSize+1), sdk.Coins{}, false},
		{addrs[0], sdk.AccAddress{}, 0, []byte("{}"), sdk.Coins{{Amount: sdk.NewInt(1)}}, false},
		{addrs[0], sdk.AccAddress{}, 0, []byte("{invalid json}"), sdk.Coins{}, false},
		{addrs[0], sdk.AccAddress{}, 0, []byte("{}"), sdk.Coins{}, true},
	}

	for i, tc := range tests {
		msg := NewMsgInstantiateContract(tc.creator, tc.admin, tc.codeID, tc.initMsg, tc.initCoins)
		require.Equal(t, RouterKey, msg.Route())
		require.Equal(t, TypeMsgInstantiateContract, msg.Type())
		if !tc.creator.Empty() {
			require.Equal(t, tc.creator, msg.GetSigners()[0])
		} else {
			require.Panics(t, func() {
				msg.GetSigners()
			})
		}

		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgExecuteContract(t *testing.T) {
	addrs := []sdk.AccAddress{
		sdk.AccAddress([]byte("addr1_______________")),
		sdk.AccAddress([]byte("addr2_______________")),
	}

	tests := []struct {
		sender     sdk.AccAddress
		contract   sdk.AccAddress
		msg        []byte
		coins      sdk.Coins
		expectPass bool
	}{
		{sdk.AccAddress{}, addrs[1], []byte("{}"), sdk.Coins{}, false},
		{addrs[0], sdk.AccAddress{}, []byte("{}"), sdk.Coins{}, false},
		{addrs[0], addrs[1], make([]byte, EnforcedMaxContractMsgSize+1), sdk.Coins{}, false},
		{addrs[0], addrs[1], []byte("{}"), sdk.Coins{{Amount: sdk.NewInt(1)}}, false},
		{addrs[0], addrs[1], []byte("{invalid json}"), sdk.Coins{}, false},
		{addrs[0], addrs[1], []byte("{}"), sdk.Coins{}, true},
	}

	for i, tc := range tests {
		msg := NewMsgExecuteContract(tc.sender, tc.contract, tc.msg, tc.coins)
		require.Equal(t, RouterKey, msg.Route())
		require.Equal(t, TypeMsgExecuteContract, msg.Type())
		if !tc.sender.Empty() {
			require.Equal(t, tc.sender, msg.GetSigners()[0])
		} else {
			require.Panics(t, func() {
				msg.GetSigners()
			})
		}

		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgMigrateContract(t *testing.T) {
	addrs := []sdk.AccAddress{
		sdk.AccAddress([]byte("addr1_______________")),
		sdk.AccAddress([]byte("addr2_______________")),
	}

	tests := []struct {
		admin      sdk.AccAddress
		contract   sdk.AccAddress
		codeID     uint64
		msg        json.RawMessage
		expectPass bool
	}{
		{sdk.AccAddress{}, addrs[1], 1, []byte("{}"), false},
		{addrs[0], sdk.AccAddress{}, 1, []byte("{}"), false},
		{addrs[0], addrs[1], 0, []byte("{}"), false},
		{addrs[0], addrs[1], 1, make([]byte, EnforcedMaxContractMsgSize+1), false},
		{addrs[0], addrs[1], 1, []byte("{invalid json}"), false},
		{addrs[0], addrs[1], 1, []byte("{}"), true},
	}

	for i, tc := range tests {
		msg := NewMsgMigrateContract(tc.admin, tc.contract, tc.codeID, tc.msg)
		require.Equal(t, RouterKey, msg.Route())
		require.Equal(t, TypeMsgMigrateContract, msg.Type())
		if !tc.admin.Empty() {
			require.Equal(t, tc.admin, msg.GetSigners()[0])
		} else {
			require.Panics(t, func() {
				msg.GetSigners()
			})
		}

		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgUpdateContractAdmin(t *testing.T) {
	addrs := []sdk.AccAddress{
		sdk.AccAddress([]byte("addr1_______________")),
		sdk.AccAddress([]byte("addr2_______________")),
		sdk.AccAddress([]byte("addr3_______________")),
	}

	tests := []struct {
		admin      sdk.AccAddress
		newAdmin   sdk.AccAddress
		contract   sdk.AccAddress
		expectPass bool
	}{
		{sdk.AccAddress{}, addrs[1], addrs[2], false},
		{addrs[0], sdk.AccAddress{}, addrs[2], false},
		{addrs[0], addrs[1], sdk.AccAddress{}, false},
		{addrs[0], addrs[1], addrs[2], true},
	}

	for i, tc := range tests {
		msg := NewMsgUpdateContractAdmin(tc.admin, tc.newAdmin, tc.contract)
		require.Equal(t, RouterKey, msg.Route())
		require.Equal(t, TypeMsgUpdateContractAdmin, msg.Type())
		if !tc.admin.Empty() {
			require.Equal(t, tc.admin, msg.GetSigners()[0])
		} else {
			require.Panics(t, func() {
				msg.GetSigners()
			})
		}

		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgClearContractMigratable(t *testing.T) {
	addrs := []sdk.AccAddress{
		sdk.AccAddress([]byte("addr1_______________")),
		sdk.AccAddress([]byte("addr2_______________")),
		sdk.AccAddress([]byte("addr3_______________")),
	}

	tests := []struct {
		admin      sdk.AccAddress
		contract   sdk.AccAddress
		expectPass bool
	}{
		{sdk.AccAddress{}, addrs[1], false},
		{addrs[0], sdk.AccAddress{}, false},
		{addrs[0], addrs[1], true},
	}

	for i, tc := range tests {
		msg := NewMsgClearContractAdmin(tc.admin, tc.contract)
		require.Equal(t, RouterKey, msg.Route())
		require.Equal(t, TypeMsgClearContractAdmin, msg.Type())
		if !tc.admin.Empty() {
			require.Equal(t, tc.admin, msg.GetSigners()[0])
		} else {
			require.Panics(t, func() {
				msg.GetSigners()
			})
		}

		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}
