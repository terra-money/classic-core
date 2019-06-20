package plan

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/terra-project/core/types"
	"github.com/terra-project/core/types/assets"
	"github.com/terra-project/core/types/util"
)

var (
	genesisTime = time.Unix(1556085600, 0)

	preseedSchedule types.VestingSchedule
	seedSchedule    types.VestingSchedule
	seedSchedule2   types.VestingSchedule
	privateSchedule types.VestingSchedule

	preseedAccounts []types.BaseGradedVestingAccount
	seedAccounts    []types.BaseGradedVestingAccount
	normalAccounts  []types.BaseGradedVestingAccount

	normalAddress = [...]string{
		"terra1wpplgwx5q2ph7z2vqm9m0t2jgr6qyjkwhxvff3",
		"terra1dp0taj85ruc299rkdvzp4z5pfg6z6swaed74e6",
	}

	preseedCoins = []sdk.Coins{
		sdk.NewCoins(
			sdk.NewCoin(assets.MicroLunaDenom, sdk.NewInt(79411554440)),
			sdk.NewCoin(assets.MicroSDRDenom, sdk.NewInt(3393296)),
		),
		sdk.NewCoins(
			sdk.NewCoin(assets.MicroLunaDenom, sdk.NewInt(20000000000000)),
		),
		sdk.NewCoins(
			sdk.NewCoin(assets.MicroLunaDenom, sdk.NewInt(2000007013377)),
			sdk.NewCoin(assets.MicroSDRDenom, sdk.NewInt(13015)),
		),
		sdk.NewCoins(
			sdk.NewCoin(assets.MicroLunaDenom, sdk.NewInt(401836579451)),
			sdk.NewCoin(assets.MicroSDRDenom, sdk.NewInt(833906)),
		),
		sdk.NewCoins(
			sdk.NewCoin(assets.MicroLunaDenom, sdk.NewInt(4176551929047)),
		),
	}

	preseedOriginalVesting = []sdk.Coins{
		sdk.NewCoins(
			sdk.NewCoin(assets.MicroLunaDenom, sdk.NewInt(30000000000000)),
		),
		sdk.NewCoins(
			sdk.NewCoin(assets.MicroLunaDenom, sdk.NewInt(20000000000000)),
		),
		sdk.NewCoins(
			sdk.NewCoin(assets.MicroLunaDenom, sdk.NewInt(10000000000000)),
		),
		sdk.NewCoins(
			sdk.NewCoin(assets.MicroLunaDenom, sdk.NewInt(5000000000000)),
		),
		sdk.NewCoins(
			sdk.NewCoin(assets.MicroLunaDenom, sdk.NewInt(4500000000000)),
		),
	}

	seedCoins = []sdk.Coins{
		sdk.NewCoins(
			sdk.NewCoin(assets.MicroLunaDenom, sdk.NewInt(1729185508867)),
			sdk.NewCoin(assets.MicroSDRDenom, sdk.NewInt(1753376)),
		),
		sdk.NewCoins(
			sdk.NewCoin(assets.MicroLunaDenom, sdk.NewInt(9347826000000)),
		),
		sdk.NewCoins(
			sdk.NewCoin(assets.MicroLunaDenom, sdk.NewInt(3913044896495)),
		),
		sdk.NewCoins(
			sdk.NewCoin(assets.MicroLunaDenom, sdk.NewInt(7373288)),
			sdk.NewCoin(assets.MicroSDRDenom, sdk.NewInt(59867)),
		),
		sdk.NewCoins(
			sdk.NewCoin(assets.MicroLunaDenom, sdk.NewInt(3826087000000)),
		),
		sdk.NewCoins(
			sdk.NewCoin(assets.MicroLunaDenom, sdk.NewInt(3913044998248)),
		),
	}

	seedOriginalVesting = []sdk.Coins{
		sdk.NewCoins(
			sdk.NewCoin(assets.MicroLunaDenom, sdk.NewInt(13043478000000)),
		),
		sdk.NewCoins(
			sdk.NewCoin(assets.MicroLunaDenom, sdk.NewInt(9347826000000)),
		),
		sdk.NewCoins(
			sdk.NewCoin(assets.MicroLunaDenom, sdk.NewInt(4347826000000)),
		),
		sdk.NewCoins(
			sdk.NewCoin(assets.MicroLunaDenom, sdk.NewInt(4347826000000)),
		),
		sdk.NewCoins(
			sdk.NewCoin(assets.MicroLunaDenom, sdk.NewInt(3826087000000)),
		),
		sdk.NewCoins(
			sdk.NewCoin(assets.MicroLunaDenom, sdk.NewInt(4347826000000)),
		),
	}

	normalCoins = []sdk.Coins{
		sdk.NewCoins(
			sdk.NewCoin(assets.MicroLunaDenom, sdk.NewInt(350000000000)),
			sdk.NewCoin(assets.MicroSDRDenom, sdk.NewInt(62500000000)),
		),
		sdk.NewCoins(
			sdk.NewCoin(assets.MicroLunaDenom, sdk.NewInt(695561902462109)),
			sdk.NewCoin(assets.MicroSDRDenom, sdk.NewInt(999916140381807)),
		),
	}

	normalOriginalVesting = []sdk.Coins{
		sdk.NewCoins(
			sdk.NewCoin(assets.MicroLunaDenom, sdk.NewInt(350000000000)),
			sdk.NewCoin(assets.MicroSDRDenom, sdk.NewInt(62500000000)),
		),
		sdk.NewCoins(
			sdk.NewCoin(assets.MicroLunaDenom, sdk.NewInt(1000000000000000)),
		),
	}
)

func init() {

	config := sdk.GetConfig()
	config.SetCoinType(330)
	config.SetFullFundraiserPath("44'/330'/0'/0/0")
	config.SetBech32PrefixForAccount(util.Bech32PrefixAccAddr, util.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(util.Bech32PrefixValAddr, util.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(util.Bech32PrefixConsAddr, util.Bech32PrefixConsPub)
	config.Seal()

	preseedSchedule = types.NewVestingSchedule(assets.MicroLunaDenom, []types.Schedule{
		types.NewSchedule(genesisTime.AddDate(0, 1, 0).Unix(), sdk.NewDecWithPrec(10, 2)),
		types.NewSchedule(genesisTime.AddDate(0, 2, 0).Unix(), sdk.NewDecWithPrec(10, 2)),
		types.NewSchedule(genesisTime.AddDate(0, 3, 0).Unix(), sdk.NewDecWithPrec(10, 2)),
		types.NewSchedule(genesisTime.AddDate(0, 12, 0).Unix(), sdk.NewDecWithPrec(70, 2)),
	})

	seedSchedule = types.NewVestingSchedule(assets.MicroLunaDenom, []types.Schedule{
		types.NewSchedule(genesisTime.AddDate(0, 1, 0).Unix(), sdk.NewDecWithPrec(10, 2)),
		types.NewSchedule(genesisTime.AddDate(0, 2, 0).Unix(), sdk.NewDecWithPrec(10, 2)),
		types.NewSchedule(genesisTime.AddDate(0, 3, 0).Unix(), sdk.NewDecWithPrec(10, 2)),
		types.NewSchedule(genesisTime.AddDate(0, 10, 0).Unix(), sdk.NewDecWithPrec(70, 2)),
	})

	seedSchedule2 = types.NewVestingSchedule(assets.MicroLunaDenom, []types.Schedule{
		types.NewSchedule(genesisTime.AddDate(0, 1, 0).Unix(), sdk.NewDecWithPrec(47, 3)),
		types.NewSchedule(genesisTime.AddDate(0, 2, 0).Unix(), sdk.NewDecWithPrec(47, 3)),
		types.NewSchedule(genesisTime.AddDate(0, 3, 0).Unix(), sdk.NewDecWithPrec(47, 3)),
		types.NewSchedule(genesisTime.AddDate(0, 10, 0).Unix(), sdk.NewDecWithPrec(326, 3)),
		types.NewSchedule(genesisTime.AddDate(0, 18, 0).Unix(), sdk.NewDecWithPrec(533, 3)),
	})

	privateSchedule = types.NewVestingSchedule(assets.MicroLunaDenom, []types.Schedule{
		types.NewSchedule(genesisTime.AddDate(0, 4, 0).Unix(), sdk.NewDecWithPrec(16, 2)),
		types.NewSchedule(genesisTime.AddDate(0, 5, 0).Unix(), sdk.NewDecWithPrec(16, 2)),
		types.NewSchedule(genesisTime.AddDate(0, 6, 0).Unix(), sdk.NewDecWithPrec(16, 2)),
		types.NewSchedule(genesisTime.AddDate(0, 7, 0).Unix(), sdk.NewDecWithPrec(16, 2)),
		types.NewSchedule(genesisTime.AddDate(0, 8, 0).Unix(), sdk.NewDecWithPrec(16, 2)),
		types.NewSchedule(genesisTime.AddDate(0, 9, 0).Unix(), sdk.NewDecWithPrec(20, 2)),
	})

	for i, bechAddr := range preseedAddresses {
		addr, _ := sdk.AccAddressFromBech32(bechAddr)
		baseAccount := &auth.BaseAccount{
			Address: addr,
			Coins:   preseedCoins[i],
		}

		baseVestingAcc := &auth.BaseVestingAccount{
			BaseAccount:     baseAccount,
			OriginalVesting: preseedOriginalVesting[i],
		}

		gradedVestingAccount := types.BaseGradedVestingAccount{
			BaseVestingAccount: baseVestingAcc,
			VestingSchedules:   []types.VestingSchedule{preseedSchedule},
		}

		preseedAccounts = append(preseedAccounts, gradedVestingAccount)
	}

	for i, bechAddr := range seedAddresses {
		addr, _ := sdk.AccAddressFromBech32(bechAddr)

		baseAccount := &auth.BaseAccount{
			Address: addr,
			Coins:   seedCoins[i],
		}

		baseVestingAcc := &auth.BaseVestingAccount{
			BaseAccount:     baseAccount,
			OriginalVesting: seedOriginalVesting[i],
		}

		gradedVestingAccount := types.BaseGradedVestingAccount{
			BaseVestingAccount: baseVestingAcc,
			VestingSchedules:   []types.VestingSchedule{seedSchedule},
		}

		if bechAddr == "terra1y9n2ywyu5dahtxar6k4z4jz97ynt8km4catuk6" {
			gradedVestingAccount.VestingSchedules = []types.VestingSchedule{seedSchedule2}
		}

		seedAccounts = append(seedAccounts, gradedVestingAccount)
	}

	for i, bechAddr := range normalAddress {
		addr, _ := sdk.AccAddressFromBech32(bechAddr)

		baseAccount := &auth.BaseAccount{
			Address: addr,
			Coins:   normalCoins[i],
		}

		baseVestingAcc := &auth.BaseVestingAccount{
			BaseAccount:     baseAccount,
			OriginalVesting: normalOriginalVesting[i],
		}

		gradedVestingAccount := types.BaseGradedVestingAccount{
			BaseVestingAccount: baseVestingAcc,
			VestingSchedules:   []types.VestingSchedule{privateSchedule},
		}

		normalAccounts = append(normalAccounts, gradedVestingAccount)
	}
}

func TestPreseedAccountUpdate(t *testing.T) {

	for _, acc := range preseedAccounts {
		lazyVestingSchedules := updatePreseedSchedules(&acc)

		require.Equal(t, 1, len(lazyVestingSchedules))
		require.Equal(t, assets.MicroLunaDenom, lazyVestingSchedules[0].GetDenom())
		lazySchedule := lazyVestingSchedules[0].LazySchedules

		require.Equal(t, genesisTime.AddDate(0, 1, 0).Unix(), lazySchedule[0].GetStartTime())
		require.Equal(t, genesisTime.AddDate(0, 2, 0).Unix(), lazySchedule[0].GetEndTime())
		require.Equal(t, sdk.NewDecWithPrec(10, 2), lazySchedule[0].GetRatio())

		require.Equal(t, genesisTime.AddDate(0, 2, 0).Unix(), lazySchedule[1].GetStartTime())
		require.Equal(t, genesisTime.AddDate(0, 11, 0).Unix(), lazySchedule[1].GetEndTime())
		require.Equal(t, sdk.NewDecWithPrec(27, 2), lazySchedule[1].GetRatio())

		require.Equal(t, genesisTime.AddDate(0, 11, 0).Unix(), lazySchedule[2].GetStartTime())
		require.Equal(t, genesisTime.AddDate(0, 13, 0).Unix(), lazySchedule[2].GetEndTime())
		require.Equal(t, sdk.NewDecWithPrec(8, 2), lazySchedule[2].GetRatio())

		require.Equal(t, genesisTime.AddDate(0, 13, 0).Unix(), lazySchedule[3].GetStartTime())
		require.Equal(t, genesisTime.AddDate(0, 17, 0).Unix(), lazySchedule[3].GetEndTime())
		require.Equal(t, sdk.NewDecWithPrec(40, 2), lazySchedule[3].GetRatio())

		require.Equal(t, genesisTime.AddDate(0, 17, 0).Unix(), lazySchedule[4].GetStartTime())
		require.Equal(t, genesisTime.AddDate(0, 18, 0).Unix(), lazySchedule[4].GetEndTime())
		require.Equal(t, sdk.NewDecWithPrec(15, 2), lazySchedule[4].GetRatio())

		require.True(t, lazyVestingSchedules[0].IsValid())
	}

}

func TestSeedAccountUpdate(t *testing.T) {

	for _, acc := range seedAccounts {
		lazyVestingSchedules := updateSeedSchedules(&acc)

		ratio := sdk.OneDec()

		if acc.GetAddress().String() == "terra1y9n2ywyu5dahtxar6k4z4jz97ynt8km4catuk6" {
			ratio = sdk.NewDecWithPrec(467, 3)
		}

		require.Equal(t, 1, len(lazyVestingSchedules))
		require.Equal(t, assets.MicroLunaDenom, lazyVestingSchedules[0].GetDenom())
		lazySchedule := lazyVestingSchedules[0].LazySchedules

		require.Equal(t, genesisTime.AddDate(0, 1, 0).Unix(), lazySchedule[0].GetStartTime())
		require.Equal(t, genesisTime.AddDate(0, 2, 0).Unix(), lazySchedule[0].GetEndTime())
		require.Equal(t, ratio.Mul(sdk.NewDecWithPrec(10, 2)), lazySchedule[0].GetRatio())

		require.Equal(t, genesisTime.AddDate(0, 2, 0).Unix(), lazySchedule[1].GetStartTime())
		require.Equal(t, genesisTime.AddDate(0, 4, 0).Unix(), lazySchedule[1].GetEndTime())
		require.Equal(t, ratio.Mul(sdk.NewDecWithPrec(6, 2)), lazySchedule[1].GetRatio())

		require.Equal(t, genesisTime.AddDate(0, 4, 0).Unix(), lazySchedule[2].GetStartTime())
		require.Equal(t, genesisTime.AddDate(0, 10, 0).Unix(), lazySchedule[2].GetEndTime())
		require.Equal(t, ratio.Mul(sdk.NewDecWithPrec(24, 2)), lazySchedule[2].GetRatio())

		require.Equal(t, genesisTime.AddDate(0, 10, 0).Unix(), lazySchedule[3].GetStartTime())
		require.Equal(t, genesisTime.AddDate(0, 13, 0).Unix(), lazySchedule[3].GetEndTime())
		require.Equal(t, ratio.Mul(sdk.NewDecWithPrec(60, 2)), lazySchedule[3].GetRatio())

		if !ratio.Equal(sdk.OneDec()) {
			require.Equal(t, genesisTime.AddDate(0, 18, 0).Unix(), lazySchedule[4].GetStartTime())
			require.Equal(t, genesisTime.AddDate(0, 19, 0).Unix(), lazySchedule[4].GetEndTime())
			require.Equal(t, sdk.OneDec().Sub(ratio), lazySchedule[4].GetRatio())
		}

		require.True(t, lazyVestingSchedules[0].IsValid())
	}

}

func TestUpdate230000(t *testing.T) {
	input := setup(t)

	for _, acc := range preseedAccounts {
		input.accKeeper.SetAccount(input.ctx, acc)
	}

	for _, acc := range seedAccounts {
		input.accKeeper.SetAccount(input.ctx, acc)
	}

	for _, acc := range normalAccounts {
		input.accKeeper.SetAccount(input.ctx, acc)
	}

	Update230000(input.ctx.WithBlockHeight(229999), input.accKeeper, input.oracleKeeper)
	require.Equal(t, sdk.NewDecWithPrec(1, 2), input.oracleKeeper.GetParams(input.ctx).OracleRewardBand)

	// not yet changed
	input.accKeeper.IterateAccounts(input.ctx, func(acc auth.Account) (stop bool) {
		stop = false

		vacc, ok := acc.(auth.VestingAccount)
		require.True(t, ok)

		_, ok = vacc.(types.GradedVestingAccount)
		require.True(t, ok)
		return
	})

	Update230000(input.ctx.WithBlockHeight(230000), input.accKeeper, input.oracleKeeper)
	require.Equal(t, sdk.NewDecWithPrec(2, 2), input.oracleKeeper.GetParams(input.ctx).OracleRewardBand)

	// not yet changed
	input.accKeeper.IterateAccounts(input.ctx, func(acc auth.Account) (stop bool) {
		stop = false

		vacc, ok := acc.(auth.VestingAccount)
		require.True(t, ok)

		_, ok = vacc.(types.LazyGradedVestingAccount)
		require.True(t, ok)
		return
	})

}
