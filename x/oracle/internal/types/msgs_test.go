package types

import (
	"encoding/hex"
	"testing"

	core "github.com/terra-project/core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/stretchr/testify/require"
)

func TestMsgExchangeRatePrevote(t *testing.T) {
	_, addrs, _, _ := mock.CreateGenAccounts(1, sdk.Coins{})

	bz, err := VoteHash("1", sdk.OneDec(), core.MicroSDRDenom, sdk.ValAddress(addrs[0]))
	require.Nil(t, err)

	tests := []struct {
		hash       string
		denom      string
		voter      sdk.AccAddress
		expectPass bool
	}{
		{hex.EncodeToString(bz), "", addrs[0], false},
		{hex.EncodeToString(bz), core.MicroCNYDenom, addrs[0], true},
		{hex.EncodeToString(bz), core.MicroCNYDenom, addrs[0], true},
		{hex.EncodeToString(bz), core.MicroCNYDenom, sdk.AccAddress{}, false},
		{"", core.MicroCNYDenom, addrs[0], false},
	}

	for i, tc := range tests {
		msg := NewMsgExchangeRatePrevote(tc.hash, tc.denom, tc.voter, sdk.ValAddress(tc.voter))
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgExchangeRateVote(t *testing.T) {
	_, addrs, _, _ := mock.CreateGenAccounts(1, sdk.Coins{})

	overflowExchangeRate, _ := sdk.NewDecFromStr("100000000000000000000000000000000000000000000000000000000")

	tests := []struct {
		denom      string
		voter      sdk.AccAddress
		salt       string
		rate       sdk.Dec
		expectPass bool
	}{
		{"", addrs[0], "123", sdk.OneDec(), false},
		{core.MicroCNYDenom, addrs[0], "123", sdk.OneDec().MulInt64(core.MicroUnit), true},
		{core.MicroCNYDenom, addrs[0], "123", sdk.ZeroDec(), true},
		{core.MicroCNYDenom, addrs[0], "123", overflowExchangeRate, false},
		{core.MicroCNYDenom, sdk.AccAddress{}, "123", sdk.OneDec().MulInt64(core.MicroUnit), false},
		{core.MicroCNYDenom, addrs[0], "", sdk.OneDec().MulInt64(core.MicroUnit), false},
	}

	for i, tc := range tests {
		msg := NewMsgExchangeRateVote(tc.rate, tc.salt, tc.denom, tc.voter, sdk.ValAddress(tc.voter))
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgFeederDelegation(t *testing.T) {
	_, addrs, _, _ := mock.CreateGenAccounts(2, sdk.Coins{})

	tests := []struct {
		delegator  sdk.ValAddress
		delegate   sdk.AccAddress
		expectPass bool
	}{
		{sdk.ValAddress(addrs[0]), addrs[1], true},
		{sdk.ValAddress{}, addrs[1], false},
		{sdk.ValAddress(addrs[0]), sdk.AccAddress{}, false},
		{nil, nil, false},
	}

	for i, tc := range tests {
		msg := NewMsgDelegateFeedConsent(tc.delegator, tc.delegate)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgAggregateExchangeRatePrevote(t *testing.T) {
	_, addrs, _, _ := mock.CreateGenAccounts(1, sdk.Coins{})

	exchangeRates := sdk.DecCoins{sdk.NewDecCoinFromDec(core.MicroSDRDenom, sdk.OneDec()), sdk.NewDecCoinFromDec(core.MicroKRWDenom, sdk.NewDecWithPrec(32121, 1))}
	bz, err := VoteHashForAggregate("1", exchangeRates.String(), sdk.ValAddress(addrs[0]))
	require.Nil(t, err)

	tests := []struct {
		hash          string
		exchangeRates sdk.DecCoins
		voter         sdk.AccAddress
		expectPass    bool
	}{
		{hex.EncodeToString(bz), exchangeRates, addrs[0], true},
		{string(bz), exchangeRates, addrs[0], false},
		{hex.EncodeToString(bz), exchangeRates, sdk.AccAddress{}, false},
		{"", exchangeRates, addrs[0], false},
	}

	for i, tc := range tests {
		msg := NewMsgAggregateExchangeRatePrevote(tc.hash, tc.voter, sdk.ValAddress(tc.voter))
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgAggregateExchangeRateVote(t *testing.T) {
	_, addrs, _, _ := mock.CreateGenAccounts(1, sdk.Coins{})

	invalidExchangeRates := "a,b"
	exchangeRates := "1.0foo,1232.132bar"
	abstainExchangeRates := "0.0foo,1232.132bar"
	overFlowExchangeRates := "100000000000000000000000000000000000000000000000000000000.0foo,1232.132bar"

	tests := []struct {
		voter         sdk.AccAddress
		salt          string
		exchangeRates string
		expectPass    bool
	}{
		{addrs[0], "123", exchangeRates, true},
		{addrs[0], "123", invalidExchangeRates, false},
		{addrs[0], "123", abstainExchangeRates, true},
		{addrs[0], "123", overFlowExchangeRates, false},
		{sdk.AccAddress{}, "123", exchangeRates, false},
		{addrs[0], "", exchangeRates, false},
	}

	for i, tc := range tests {
		msg := NewMsgAggregateExchangeRateVote(tc.salt, tc.exchangeRates, tc.voter, sdk.ValAddress(tc.voter))
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}
