package market

import (
	"terra/types/assets"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/stretchr/testify/require"
)

func TestPriceFeedMsg(t *testing.T) {
	_, addrs, _, _ := mock.CreateGenAccounts(1, sdk.Coins{})
	tests := []struct {
		offerCoin  sdk.Coin
		askDenom   string
		expectPass bool
	}{
		{sdk.NewInt64Coin(assets.KRWDenom, 10), assets.LunaDenom, true},
		{sdk.NewInt64Coin(assets.LunaDenom, 10), assets.USDDenom, true},
		{sdk.NewInt64Coin(assets.USDDenom, 10), assets.USDDenom, false},
	}

	for i, tc := range tests {
		msg := NewMsgSwap(addrs[0], tc.offerCoin, tc.askDenom)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}
