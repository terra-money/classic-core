package types

import (
	"testing"

	core "github.com/terra-money/core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/stretchr/testify/require"
)

func TestMsgSwap(t *testing.T) {
	_, addrs, _, _ := mock.CreateGenAccounts(1, sdk.Coins{})

	overflowOfferAmt, _ := sdk.NewIntFromString("100000000000000000000000000000000000000000000000000000000")

	tests := []struct {
		trader     sdk.AccAddress
		offerCoin  sdk.Coin
		askDenom   string
		expectPass bool
	}{
		{addrs[0], sdk.NewCoin(core.MicroLunaDenom, sdk.OneInt()), core.MicroSDRDenom, true},
		{sdk.AccAddress{}, sdk.NewCoin(core.MicroLunaDenom, sdk.OneInt()), core.MicroSDRDenom, false},
		{addrs[0], sdk.NewCoin(core.MicroLunaDenom, sdk.ZeroInt()), core.MicroSDRDenom, false},
		{addrs[0], sdk.NewCoin(core.MicroLunaDenom, overflowOfferAmt), core.MicroSDRDenom, false},
		{addrs[0], sdk.NewCoin(core.MicroLunaDenom, sdk.OneInt()), core.MicroLunaDenom, false},
	}

	for i, tc := range tests {
		msg := NewMsgSwap(tc.trader, tc.offerCoin, tc.askDenom)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgSwapSend(t *testing.T) {
	_, addrs, _, _ := mock.CreateGenAccounts(1, sdk.Coins{})

	overflowOfferAmt, _ := sdk.NewIntFromString("100000000000000000000000000000000000000000000000000000000")

	tests := []struct {
		fromAddress sdk.AccAddress
		toAddress   sdk.AccAddress
		offerCoin   sdk.Coin
		askDenom    string
		expectPass  bool
	}{
		{addrs[0], addrs[0], sdk.NewCoin(core.MicroLunaDenom, sdk.OneInt()), core.MicroSDRDenom, true},
		{addrs[0], sdk.AccAddress{}, sdk.NewCoin(core.MicroLunaDenom, sdk.OneInt()), core.MicroSDRDenom, false},
		{sdk.AccAddress{}, addrs[0], sdk.NewCoin(core.MicroLunaDenom, sdk.OneInt()), core.MicroSDRDenom, false},
		{addrs[0], addrs[0], sdk.NewCoin(core.MicroLunaDenom, sdk.ZeroInt()), core.MicroSDRDenom, false},
		{addrs[0], addrs[0], sdk.NewCoin(core.MicroLunaDenom, overflowOfferAmt), core.MicroSDRDenom, false},
		{addrs[0], addrs[0], sdk.NewCoin(core.MicroLunaDenom, sdk.OneInt()), core.MicroLunaDenom, false},
	}

	for i, tc := range tests {
		msg := NewMsgSwapSend(tc.fromAddress, tc.toAddress, tc.offerCoin, tc.askDenom)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}
