package oracle

import (
	"github.com/terra-project/core/types/assets"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/stretchr/testify/require"
)

func TestMsgPriceFeed(t *testing.T) {
	_, addrs, _, _ := mock.CreateGenAccounts(1, sdk.Coins{})
	tests := []struct {
		denom      string
		price      sdk.Dec
		voter      sdk.AccAddress
		expectPass bool
	}{
		{"", sdk.OneDec(), addrs[0], false},
		{assets.MicroCNYDenom, sdk.OneDec().MulInt64(assets.MicroUnit), addrs[0], true},
		{assets.MicroCNYDenom, sdk.ZeroDec().MulInt64(assets.MicroUnit), addrs[0], false},
		{assets.MicroCNYDenom, sdk.OneDec().MulInt64(assets.MicroUnit), sdk.AccAddress{}, false},
	}

	for i, tc := range tests {
		msg := NewMsgPriceFeed(tc.denom, tc.price, tc.voter, sdk.ValAddress(tc.voter))
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}
