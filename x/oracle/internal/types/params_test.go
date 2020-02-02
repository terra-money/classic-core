package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestParamsEqual(t *testing.T) {
	p1 := DefaultParams()
	err := p1.ValidateBasic()
	require.NoError(t, err)

	// minus vote period
	p1.VotePeriod = -1
	err = p1.ValidateBasic()
	require.Error(t, err)

	// small vote threshold
	p2 := DefaultParams()
	p2.VoteThreshold = sdk.ZeroDec()
	err = p2.ValidateBasic()
	require.Error(t, err)

	// negative reward band
	p3 := DefaultParams()
	p3.RewardBand = sdk.NewDecWithPrec(-1, 2)
	err = p3.ValidateBasic()
	require.Error(t, err)

	// negative slash fraction
	p4 := DefaultParams()
	p4.SlashFraction = sdk.NewDec(-1)
	err = p4.ValidateBasic()
	require.Error(t, err)

	// negative min valid per window
	p5 := DefaultParams()
	p5.MinValidPerWindow = sdk.NewDec(-1)
	err = p5.ValidateBasic()
	require.Error(t, err)

	// small slash window
	p6 := DefaultParams()
	p6.SlashWindow = int64(1)
	err = p6.ValidateBasic()
	require.Error(t, err)

	// small distribution window
	p7 := DefaultParams()
	p7.RewardDistributionWindow = int64(1)
	err = p7.ValidateBasic()
	require.Error(t, err)

	p8 := DefaultParams()
	require.NotNil(t, p8.ParamSetPairs())
	require.NotNil(t, p8.String())
}
