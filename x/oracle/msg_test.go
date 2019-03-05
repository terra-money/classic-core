package oracle

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/stretchr/testify/require"
)

func TestPriceFeedMsg(t *testing.T) {
	_, addrs, _, _ := mock.CreateGenAccounts(1, sdk.Coins{})
	tests := []struct {
		denom      string
		price      sdk.Dec
		feeder     sdk.AccAddress
		expectPass bool
	}{
		{"", sdk.OneDec(), addrs[0], false},
		{"terra", sdk.OneDec(), addrs[0], true},
		{"terra", sdk.ZeroDec(), addrs[0], false},
		{"terra", sdk.OneDec(), addrs[0], false},
		{"terra", sdk.OneDec(), sdk.AccAddress{}, false},
	}

	for i, tc := range tests {
		msg := NewPriceFeedMsg(tc.denom, tc.price, tc.feeder)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}
