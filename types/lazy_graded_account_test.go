package types

import (
	"testing"
	"time"

	"github.com/terra-project/core/types/assets"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/stretchr/testify/require"
)

var (
	angelLazySchedule        LazyVestingSchedule
	seedLazySchedule         LazyVestingSchedule
	privateLazySchedule      LazyVestingSchedule
	privateBonusLazySchedule LazyVestingSchedule
)

// initialize the times!
func init() {
	var err error
	timeLayoutString := "2006-01-02 15:04:05 -0700 MST"
	timeGenesis, err = time.Parse(timeLayoutString, timeGenesisString)
	if err != nil {
		panic(err)
	}

	monthlyTimes = []int64{}
	for i := 0; i < 4; i++ {
		for j := 0; j < 12; j++ {
			monthlyTimes = append(monthlyTimes, timeGenesis.AddDate(i, j, 0).Unix())
		}
	}

	angelLazySchedule = NewLazyVestingSchedule(assets.MicroLunaDenom, []LazySchedule{
		NewLazySchedule(monthlyTimes[1], monthlyTimes[2], sdk.NewDecWithPrec(10, 2)),
		NewLazySchedule(monthlyTimes[2], monthlyTimes[3], sdk.NewDecWithPrec(10, 2)),
		NewLazySchedule(monthlyTimes[3], monthlyTimes[4], sdk.NewDecWithPrec(10, 2)),
		NewLazySchedule(monthlyTimes[12], monthlyTimes[13], sdk.NewDecWithPrec(70, 2)),
	})

	seedLazySchedule = NewLazyVestingSchedule(assets.MicroLunaDenom, []LazySchedule{
		NewLazySchedule(monthlyTimes[1], monthlyTimes[2], sdk.NewDecWithPrec(10, 2)),
		NewLazySchedule(monthlyTimes[2], monthlyTimes[3], sdk.NewDecWithPrec(10, 2)),
		NewLazySchedule(monthlyTimes[3], monthlyTimes[4], sdk.NewDecWithPrec(10, 2)),
		NewLazySchedule(monthlyTimes[10], monthlyTimes[11], sdk.NewDecWithPrec(70, 2)),
	})

	privateLazySchedule = NewLazyVestingSchedule(assets.MicroLunaDenom, []LazySchedule{
		NewLazySchedule(monthlyTimes[3], monthlyTimes[9], sdk.NewDec(1)),
	})

	privateBonusLazySchedule = NewLazyVestingSchedule(assets.MicroLunaDenom, []LazySchedule{
		NewLazySchedule(monthlyTimes[6], monthlyTimes[18], sdk.NewDec(1)),
	})

}

func TestGetVestedCoinsLazyGradVestingAcc(t *testing.T) {
	genesis := timeGenesis

	_, _, addr := keyPubAddr()
	origCoins := sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 10000)}
	bacc := auth.NewBaseAccountWithAddress(addr)
	bacc.SetCoins(origCoins)

	// require no coins are vested until schedule maturation
	lgva := NewBaseLazyGradedVestingAccount(&bacc, []LazyVestingSchedule{
		angelLazySchedule,
	})
	vestedCoins := lgva.GetVestedCoins(genesis)
	require.Nil(t, vestedCoins)

	// require coins be vested at the expected cliff
	vestedCoins = lgva.GetVestedCoins(timeGenesis.AddDate(0, 1, 0))
	require.True(t, vestedCoins.Empty())

	vestedCoins = lgva.GetVestedCoins(timeGenesis.AddDate(0, 2, 0))
	require.Equal(t, scaleCoins(0.1, assets.MicroLunaDenom, origCoins), vestedCoins)

	vestedCoins = lgva.GetVestedCoins(timeGenesis.AddDate(0, 2, 15))
	require.Equal(t, scaleCoins(0.15, assets.MicroLunaDenom, origCoins), vestedCoins)

	vestedCoins = lgva.GetVestedCoins(timeGenesis.AddDate(0, 3, 0))
	require.Equal(t, scaleCoins(0.2, assets.MicroLunaDenom, origCoins), vestedCoins)

	vestedCoins = lgva.GetVestedCoins(timeGenesis.AddDate(0, 4, 0))
	require.Equal(t, scaleCoins(0.3, assets.MicroLunaDenom, origCoins), vestedCoins)

	vestedCoins = lgva.GetVestedCoins(timeGenesis.AddDate(0, 5, 0))
	require.Equal(t, scaleCoins(0.3, assets.MicroLunaDenom, origCoins), vestedCoins)

	vestedCoins = lgva.GetVestedCoins(timeGenesis.AddDate(1, 0, 0))
	require.Equal(t, scaleCoins(0.3, assets.MicroLunaDenom, origCoins), vestedCoins)

	vestedCoins = lgva.GetVestedCoins(timeGenesis.AddDate(1, 1, 0))
	require.Equal(t, scaleCoins(1.0, assets.MicroLunaDenom, origCoins), vestedCoins)
}

func TestGetVestingCoinsLazyGradVestingAcc(t *testing.T) {
	genesis := timeGenesis

	_, _, addr := keyPubAddr()
	origCoins := sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 10000)}
	bacc := auth.NewBaseAccountWithAddress(addr)
	bacc.SetCoins(origCoins)

	// require no coins are vested until schedule maturation
	lgva := NewBaseLazyGradedVestingAccount(&bacc, []LazyVestingSchedule{
		angelLazySchedule,
	})
	vestingCoins := lgva.GetVestingCoins(genesis)
	require.Equal(t, vestingCoins, origCoins)

	// require coins be vested at the expected cliff
	vestingCoins = lgva.GetVestingCoins(timeGenesis.AddDate(0, 1, 0))
	require.Equal(t, origCoins, vestingCoins)

	vestingCoins = lgva.GetVestingCoins(timeGenesis.AddDate(0, 2, 0))
	require.Equal(t, origCoins.Sub(scaleCoins(0.1, assets.MicroLunaDenom, origCoins)), vestingCoins)

	vestingCoins = lgva.GetVestingCoins(timeGenesis.AddDate(0, 2, 15))
	require.Equal(t, origCoins.Sub(scaleCoins(0.15, assets.MicroLunaDenom, origCoins)), vestingCoins)

	vestingCoins = lgva.GetVestingCoins(timeGenesis.AddDate(0, 3, 0))
	require.Equal(t, origCoins.Sub(scaleCoins(0.2, assets.MicroLunaDenom, origCoins)), vestingCoins)

	vestingCoins = lgva.GetVestingCoins(timeGenesis.AddDate(0, 4, 0))
	require.Equal(t, origCoins.Sub(scaleCoins(0.3, assets.MicroLunaDenom, origCoins)), vestingCoins)

	vestingCoins = lgva.GetVestingCoins(timeGenesis.AddDate(0, 5, 0))
	require.Equal(t, origCoins.Sub(scaleCoins(0.3, assets.MicroLunaDenom, origCoins)), vestingCoins)

	vestingCoins = lgva.GetVestingCoins(timeGenesis.AddDate(1, 0, 0))
	require.Equal(t, origCoins.Sub(scaleCoins(0.3, assets.MicroLunaDenom, origCoins)), vestingCoins)

	vestingCoins = lgva.GetVestingCoins(timeGenesis.AddDate(1, 1, 0))
	require.Equal(t, origCoins.Sub(scaleCoins(1.0, assets.MicroLunaDenom, origCoins)), vestingCoins)
}

func TestSpendableCoinsLazyGradVestingAcc(t *testing.T) {
	genesis := timeGenesis

	_, _, addr := keyPubAddr()
	origCoins := sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 10000), sdk.NewInt64Coin(assets.MicroSDRDenom, 10000)}
	sdrCoins := sdk.Coins{sdk.NewInt64Coin(assets.MicroSDRDenom, 10000)}
	bacc := auth.NewBaseAccountWithAddress(addr)
	bacc.SetCoins(origCoins)

	// require no coins are vested until schedule maturation
	lgva := NewBaseLazyGradedVestingAccount(&bacc, []LazyVestingSchedule{
		angelLazySchedule,
	})

	spendableCoins := lgva.SpendableCoins(genesis)
	require.Equal(t, sdrCoins, spendableCoins)

	// require that all coins are spendable after the maturation of the vesting
	// schedule
	spendableCoins = lgva.SpendableCoins(timeGenesis.AddDate(1, 1, 0))
	require.Equal(t, origCoins, spendableCoins)

	// require that all luna coins are still vesting after some time
	spendableCoins = lgva.SpendableCoins(genesis.Add(12 * time.Hour))
	require.Equal(t, spendableCoins, sdrCoins)

	// receive some coins
	relgvamt := sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 50)}
	lgva.SetCoins(lgva.GetCoins().Add(relgvamt))

	// require that only received coins and sdrCoins are spendable since the account is still
	// vesting
	spendableCoins = lgva.SpendableCoins(genesis.Add(12 * time.Hour))
	require.Equal(t, relgvamt.Add(sdrCoins), spendableCoins)

	// spend all spendable coins
	lgva.SetCoins(lgva.GetCoins().Sub(spendableCoins))

	// require that no more coins are spendable
	spendableCoins = lgva.SpendableCoins(genesis.Add(12 * time.Hour))
	require.True(t, spendableCoins.Empty())
}

func TestTrackDelegationLazyGradVestingAcc(t *testing.T) {
	genesis := timeGenesis

	_, _, addr := keyPubAddr()
	origCoins := sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 10000)}
	bacc := auth.NewBaseAccountWithAddress(addr)

	bacc.SetCoins(origCoins)
	lgva := NewBaseLazyGradedVestingAccount(&bacc, []LazyVestingSchedule{
		angelLazySchedule,
	})
	lgva.TrackDelegation(genesis, origCoins)
	require.Equal(t, origCoins, lgva.DelegatedVesting)
	require.Nil(t, lgva.DelegatedFree)
	require.Nil(t, lgva.GetCoins())

	// require the ability to delegate all vested coins
	bacc.SetCoins(origCoins)
	lgva = NewBaseLazyGradedVestingAccount(&bacc, []LazyVestingSchedule{
		angelLazySchedule,
	})
	lgva.TrackDelegation(genesis.AddDate(1, 1, 0), origCoins)
	require.Nil(t, lgva.DelegatedVesting)
	require.Equal(t, origCoins, lgva.DelegatedFree)
	require.Nil(t, lgva.GetCoins())

	// require the ability to delegate all coins half way through the vesting
	// schedule
	bacc.SetCoins(origCoins)
	lgva = NewBaseLazyGradedVestingAccount(&bacc, []LazyVestingSchedule{
		angelLazySchedule,
	})
	lgva.TrackDelegation(genesis.AddDate(0, 3, 0), sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 3000)})
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 3000)}, lgva.DelegatedVesting)
	require.Nil(t, lgva.DelegatedFree)

	lgva.TrackDelegation(genesis.AddDate(0, 4, 0), sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 7000)})
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 7000)}, lgva.DelegatedVesting)
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 3000)}, lgva.DelegatedFree)
	require.Nil(t, lgva.GetCoins())

	// require no modifications when delegation amount is zero or not enough funds
	bacc.SetCoins(origCoins)
	lgva = NewBaseLazyGradedVestingAccount(&bacc, []LazyVestingSchedule{
		angelLazySchedule,
	})

	require.Panics(t, func() {
		lgva.TrackDelegation(genesis.AddDate(1, 0, 0), sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 1000000)})
	})
	require.Nil(t, lgva.DelegatedVesting)
	require.Nil(t, lgva.DelegatedFree)
	require.Equal(t, origCoins, lgva.GetCoins())
}

func TestTrackUndelegationLazyGradVestingAcc(t *testing.T) {
	genesis := timeGenesis

	_, _, addr := keyPubAddr()
	origCoins := sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 10000), sdk.NewInt64Coin(assets.MicroSDRDenom, 10000)}
	bacc := auth.NewBaseAccountWithAddress(addr)
	bacc.SetCoins(origCoins)

	// require the ability to undelegate all vesting coins
	lgva := NewBaseLazyGradedVestingAccount(&bacc, []LazyVestingSchedule{
		angelLazySchedule,
	})
	lgva.TrackDelegation(genesis, origCoins)
	lgva.TrackUndelegation(origCoins)
	require.Nil(t, lgva.DelegatedFree)
	require.Nil(t, lgva.DelegatedVesting)
	require.Equal(t, origCoins, lgva.GetCoins())

	// require the ability to undelegate all vested coins
	bacc.SetCoins(origCoins)
	lgva = NewBaseLazyGradedVestingAccount(&bacc, []LazyVestingSchedule{
		angelLazySchedule,
	})
	lgva.TrackDelegation(genesis.AddDate(1, 1, 0), origCoins)
	lgva.TrackUndelegation(origCoins)
	require.Nil(t, lgva.DelegatedFree)
	require.Nil(t, lgva.DelegatedVesting)
	require.Equal(t, origCoins, lgva.GetCoins())

	// require no modifications when the undelegation amount is zero
	bacc.SetCoins(origCoins)
	lgva = NewBaseLazyGradedVestingAccount(&bacc, []LazyVestingSchedule{
		angelLazySchedule,
	})
	require.Panics(t, func() {
		lgva.TrackUndelegation(sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 0)})
	})
	require.Nil(t, lgva.DelegatedFree)
	require.Nil(t, lgva.DelegatedVesting)
	require.Equal(t, origCoins, lgva.GetCoins())

	// vest 50% and delegate to two validators
	lgva = NewBaseLazyGradedVestingAccount(&bacc, []LazyVestingSchedule{
		angelLazySchedule,
	})
	lgva.TrackDelegation(genesis.AddDate(0, 4, 0), sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 5000)})
	lgva.TrackDelegation(genesis.AddDate(0, 4, 0), sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 5000)})

	// undelegate from one validator that got slashed 50%
	lgva.TrackUndelegation(sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 2500)})
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 500)}, lgva.DelegatedFree)
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 7000)}, lgva.DelegatedVesting)
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 2500), sdk.NewInt64Coin(assets.MicroSDRDenom, 10000)}, lgva.GetCoins())

	// undelegate from the other validator that did not get slashed
	lgva.TrackUndelegation(sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 5000)})
	require.Nil(t, lgva.DelegatedFree)
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 2500)}, lgva.DelegatedVesting)
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 7500), sdk.NewInt64Coin(assets.MicroSDRDenom, 10000)}, lgva.GetCoins())
}

func TestStringLazyGradVestingAcc(t *testing.T) {

	_, _, addr := keyPubAddr()
	origCoins := sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 10000), sdk.NewInt64Coin(assets.MicroSDRDenom, 10000)}
	bacc := auth.NewBaseAccountWithAddress(addr)
	bacc.SetCoins(origCoins)

	// require the ability to undelegate all vesting coins
	lgva := NewBaseLazyGradedVestingAccount(&bacc, []LazyVestingSchedule{
		angelLazySchedule,
	})
	require.NotNil(t, lgva.String())

	lazyVestingSchedule, found := lgva.GetLazyVestingSchedule(assets.MicroLunaDenom)
	require.True(t, found)
	require.NotNil(t, lazyVestingSchedule.String())
}

func TestIsValidLazyGradVestingAcc(t *testing.T) {

	_, _, addr := keyPubAddr()
	origCoins := sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 10000), sdk.NewInt64Coin(assets.MicroSDRDenom, 10000)}
	bacc := auth.NewBaseAccountWithAddress(addr)
	bacc.SetCoins(origCoins)

	// require the ability to undelegate all vesting coins
	angelAccount := NewBaseLazyGradedVestingAccount(&bacc, []LazyVestingSchedule{
		angelLazySchedule,
	})

	lazyVestingSchedule, found := angelAccount.GetLazyVestingSchedule(assets.MicroLunaDenom)
	require.True(t, found)
	require.True(t, lazyVestingSchedule.IsValid())

	seedAccount := NewBaseLazyGradedVestingAccount(&bacc, []LazyVestingSchedule{
		seedLazySchedule,
	})

	lazyVestingSchedule, found = seedAccount.GetLazyVestingSchedule(assets.MicroLunaDenom)
	require.True(t, found)
	require.True(t, lazyVestingSchedule.IsValid())

	privateAccount := NewBaseLazyGradedVestingAccount(&bacc, []LazyVestingSchedule{
		privateLazySchedule,
	})

	lazyVestingSchedule, found = privateAccount.GetLazyVestingSchedule(assets.MicroLunaDenom)
	require.True(t, found)
	require.True(t, lazyVestingSchedule.IsValid())

	employeeAccount := NewBaseLazyGradedVestingAccount(&bacc, []LazyVestingSchedule{
		privateBonusLazySchedule,
	})

	lazyVestingSchedule, found = employeeAccount.GetLazyVestingSchedule(assets.MicroLunaDenom)
	require.True(t, found)
	require.True(t, lazyVestingSchedule.IsValid())
}
