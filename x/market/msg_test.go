package market

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
		offerCoin  sdk.Coin
		askDenom   string
		expectPass bool
	}{
		{sdk.NewInt64Coin(assets.MicroKRWDenom, sdk.NewInt(10).MulRaw(assets.MicroUnit).Int64()), assets.MicroLunaDenom, true},
		{sdk.NewInt64Coin(assets.MicroLunaDenom, sdk.NewInt(10).MulRaw(assets.MicroUnit).Int64()), assets.MicroUSDDenom, true},
		{sdk.NewInt64Coin(assets.MicroUSDDenom, sdk.NewInt(10).MulRaw(assets.MicroUnit).Int64()), assets.MicroUSDDenom, false},
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
