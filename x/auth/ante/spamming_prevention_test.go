package ante_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/terra-project/core/x/auth/ante"
	"github.com/terra-project/core/x/oracle"
)

func TestOracleSpamming(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasmtest")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	viper.Set(flags.FlagHome, tempDir)

	_, ctx := createTestApp()
	_, _, addr1 := types.KeyTestPubAddr()
	_, _, addr2 := types.KeyTestPubAddr()

	spd := ante.NewSpammingPreventionDecorator(dummyOracleKeeper{
		feeders: map[string]string{
			sdk.ValAddress(addr1).String(): addr1.String(),
			sdk.ValAddress(addr2).String(): addr2.String(),
		},
	})

	// normal so ok
	ctx = ctx.WithBlockHeight(100)
	require.NoError(t, spd.CheckOracleSpamming(ctx, []sdk.Msg{
		oracle.NewMsgAggregateExchangeRatePrevote(oracle.AggregateVoteHash{}, addr1, sdk.ValAddress(addr1)),
		oracle.NewMsgAggregateExchangeRateVote("", "", addr1, sdk.ValAddress(addr1)),
	}))

	// do it again is blocked
	require.Error(t, spd.CheckOracleSpamming(ctx, []sdk.Msg{
		oracle.NewMsgAggregateExchangeRatePrevote(oracle.AggregateVoteHash{}, addr1, sdk.ValAddress(addr1)),
		oracle.NewMsgAggregateExchangeRateVote("", "", addr1, sdk.ValAddress(addr1)),
	}))

	// next block; can put oracle again
	ctx = ctx.WithBlockHeight(101)
	require.NoError(t, spd.CheckOracleSpamming(ctx, []sdk.Msg{
		oracle.NewMsgAggregateExchangeRatePrevote(oracle.AggregateVoteHash{}, addr1, sdk.ValAddress(addr1)),
		oracle.NewMsgAggregateExchangeRateVote("", "", addr1, sdk.ValAddress(addr1)),
	}))

	// catch wrong feeder
	ctx = ctx.WithBlockHeight(102)
	require.Error(t, spd.CheckOracleSpamming(ctx, []sdk.Msg{
		oracle.NewMsgAggregateExchangeRatePrevote(oracle.AggregateVoteHash{}, addr2, sdk.ValAddress(addr1)),
		oracle.NewMsgAggregateExchangeRateVote("", "", addr1, sdk.ValAddress(addr1)),
	}))

	// catch wrong feeder
	ctx = ctx.WithBlockHeight(103)
	require.Error(t, spd.CheckOracleSpamming(ctx, []sdk.Msg{
		oracle.NewMsgAggregateExchangeRatePrevote(oracle.AggregateVoteHash{}, addr1, sdk.ValAddress(addr1)),
		oracle.NewMsgAggregateExchangeRateVote("", "", addr2, sdk.ValAddress(addr1)),
	}))
}

type dummyOracleKeeper struct {
	feeders map[string]string
}

func (ok dummyOracleKeeper) ValidateFeeder(ctx sdk.Context, feederAddr sdk.AccAddress, validatorAddr sdk.ValAddress, checkBonded bool) error {
	if val, ok := ok.feeders[validatorAddr.String()]; ok && val == feederAddr.String() {
		return nil
	}

	return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "cannot ensure feeder right")
}
