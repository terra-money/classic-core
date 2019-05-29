package oracle

import (
	"encoding/hex"
	"testing"

	"github.com/terra-project/core/types/assets"

	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/stretchr/testify/require"
)

func TestMsgPriceFeed(t *testing.T) {
	_, addrs, _, _ := mock.CreateGenAccounts(1, sdk.Coins{})

	salt := "1"
	bz, err := VoteHash("1", sdk.OneDec(), assets.MicroSDRDenom, types.ValAddress(addrs[0]))
	require.Nil(t, err)

	tests := []struct {
		hash       string
		denom      string
		voter      sdk.AccAddress
		salt       string
		price      sdk.Dec
		expectPass bool
	}{
		{hex.EncodeToString(bz), "", addrs[0], salt, sdk.OneDec(), false},
		{hex.EncodeToString(bz), assets.MicroCNYDenom, addrs[0], salt, sdk.OneDec().MulInt64(assets.MicroUnit), true},
		{hex.EncodeToString(bz), assets.MicroCNYDenom, addrs[0], "123", sdk.ZeroDec(), true},
		{hex.EncodeToString(bz), assets.MicroCNYDenom, sdk.AccAddress{}, salt, sdk.OneDec().MulInt64(assets.MicroUnit), false},
	}

	for i, tc := range tests {
		msg := NewMsgPriceFeed(tc.hash, tc.salt, tc.denom, tc.voter, sdk.ValAddress(tc.voter), tc.price)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}
