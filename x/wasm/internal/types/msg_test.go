package types

import (
	"testing"

	core "github.com/terra-money/core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/stretchr/testify/require"
)

func TestMsgStoreCode(t *testing.T) {
	_, addrs, _, _ := mock.CreateGenAccounts(1, sdk.Coins{})

	tests := []struct {
		sender       sdk.AccAddress
		wasmByteCode core.Base64Bytes
		expectPass   bool
	}{
		{addrs[0], []byte{}, false},
		{sdk.AccAddress{}, []byte{1, 2, 3}, false},
		{addrs[0], make([]byte, EnforcedMaxContractSize+1), false},
		{addrs[0], []byte{1, 2, 3}, true},
	}

	for i, tc := range tests {
		msg := NewMsgStoreCode(tc.sender, tc.wasmByteCode)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgInstantiateCode(t *testing.T) {
	_, addrs, _, _ := mock.CreateGenAccounts(1, sdk.Coins{})

	tests := []struct {
		creator    sdk.AccAddress
		codeID     uint64
		initMsg    core.Base64Bytes
		initCoins  sdk.Coins
		expectPass bool
	}{
		{sdk.AccAddress{}, 0, []byte{}, sdk.Coins{}, false},
		{addrs[0], 0, make([]byte, EnforcedMaxContractMsgSize+1), sdk.Coins{}, false},
		{addrs[0], 0, []byte{}, sdk.Coins{{Amount: sdk.NewInt(1)}}, false},
		{addrs[0], 0, []byte{}, sdk.Coins{}, true},
	}

	for i, tc := range tests {
		msg := NewMsgInstantiateContract(tc.creator, tc.codeID, tc.initMsg, tc.initCoins, true)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgExecuteContract(t *testing.T) {
	_, addrs, _, _ := mock.CreateGenAccounts(2, sdk.Coins{})

	tests := []struct {
		sender     sdk.AccAddress
		contract   sdk.AccAddress
		msg        core.Base64Bytes
		coins      sdk.Coins
		expectPass bool
	}{
		{sdk.AccAddress{}, addrs[1], []byte{}, sdk.Coins{}, false},
		{addrs[0], sdk.AccAddress{}, []byte{}, sdk.Coins{}, false},
		{addrs[0], addrs[1], make([]byte, EnforcedMaxContractMsgSize+1), sdk.Coins{}, false},
		{addrs[0], addrs[1], []byte{}, sdk.Coins{{Amount: sdk.NewInt(1)}}, false},
		{addrs[0], addrs[1], []byte{}, sdk.Coins{}, true},
	}

	for i, tc := range tests {
		msg := NewMsgExecuteContract(tc.sender, tc.contract, tc.msg, tc.coins)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgMigrateContract(t *testing.T) {
	_, addrs, _, _ := mock.CreateGenAccounts(2, sdk.Coins{})

	tests := []struct {
		owner      sdk.AccAddress
		contract   sdk.AccAddress
		codeID     uint64
		msg        core.Base64Bytes
		expectPass bool
	}{
		{sdk.AccAddress{}, addrs[1], 1, []byte{}, false},
		{addrs[0], sdk.AccAddress{}, 1, []byte{}, false},
		{addrs[0], addrs[1], 0, []byte{}, false},
		{addrs[0], addrs[1], 1, make([]byte, EnforcedMaxContractMsgSize+1), false},
		{addrs[0], addrs[1], 1, []byte{}, true},
	}

	for i, tc := range tests {
		msg := NewMsgMigrateContract(tc.owner, tc.contract, tc.codeID, tc.msg)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgUpdateContractOwner(t *testing.T) {
	_, addrs, _, _ := mock.CreateGenAccounts(3, sdk.Coins{})

	tests := []struct {
		owner      sdk.AccAddress
		newOwner   sdk.AccAddress
		contract   sdk.AccAddress
		expectPass bool
	}{
		{sdk.AccAddress{}, addrs[1], addrs[2], false},
		{addrs[0], sdk.AccAddress{}, addrs[2], false},
		{addrs[0], addrs[1], sdk.AccAddress{}, false},
		{addrs[0], addrs[1], addrs[2], true},
	}

	for i, tc := range tests {
		msg := NewMsgUpdateContractOwner(tc.owner, tc.newOwner, tc.contract)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}
