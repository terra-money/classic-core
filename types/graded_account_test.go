package types

import (
	"testing"
	"time"

	"github.com/terra-project/core/types/assets"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/stretchr/testify/require"
)

func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := secp256k1.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}

var (
	angelSchedule        VestingSchedule
	seedSchedule         VestingSchedule
	privateSchedule      VestingSchedule
	privateBonusSchedule VestingSchedule
	employeeSchedule     VestingSchedule
	timeGenesisString    = "2019-04-23 23:00:00 -0800 PST"
	monthlyTimes         []int64
	timeGenesis          time.Time
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

	angelSchedule = NewVestingSchedule(assets.MicroLunaDenom, []Schedule{
		NewSchedule(monthlyTimes[1], sdk.NewDecWithPrec(10, 2)),
		NewSchedule(monthlyTimes[2], sdk.NewDecWithPrec(10, 2)),
		NewSchedule(monthlyTimes[3], sdk.NewDecWithPrec(10, 2)),
		NewSchedule(monthlyTimes[12], sdk.NewDecWithPrec(70, 2)),
	})

	seedSchedule = NewVestingSchedule(assets.MicroLunaDenom, []Schedule{
		NewSchedule(monthlyTimes[1], sdk.NewDecWithPrec(10, 2)),
		NewSchedule(monthlyTimes[2], sdk.NewDecWithPrec(10, 2)),
		NewSchedule(monthlyTimes[3], sdk.NewDecWithPrec(10, 2)),
		NewSchedule(monthlyTimes[10], sdk.NewDecWithPrec(70, 2)),
	})

	privateSchedule = NewVestingSchedule(assets.MicroLunaDenom, []Schedule{
		NewSchedule(monthlyTimes[3], sdk.NewDecWithPrec(16, 2)),
		NewSchedule(monthlyTimes[4], sdk.NewDecWithPrec(17, 2)),
		NewSchedule(monthlyTimes[5], sdk.NewDecWithPrec(16, 2)),
		NewSchedule(monthlyTimes[6], sdk.NewDecWithPrec(17, 2)),
		NewSchedule(monthlyTimes[7], sdk.NewDecWithPrec(17, 2)),
		NewSchedule(monthlyTimes[8], sdk.NewDecWithPrec(17, 2)),
	})

	privateBonusSchedule = NewVestingSchedule(assets.MicroLunaDenom, []Schedule{
		NewSchedule(monthlyTimes[6], sdk.NewDecWithPrec(8, 2)),
		NewSchedule(monthlyTimes[7], sdk.NewDecWithPrec(8, 2)),
		NewSchedule(monthlyTimes[8], sdk.NewDecWithPrec(8, 2)),
		NewSchedule(monthlyTimes[9], sdk.NewDecWithPrec(8, 2)),
		NewSchedule(monthlyTimes[10], sdk.NewDecWithPrec(8, 2)),
		NewSchedule(monthlyTimes[11], sdk.NewDecWithPrec(8, 2)),
		NewSchedule(monthlyTimes[12], sdk.NewDecWithPrec(8, 2)),
		NewSchedule(monthlyTimes[13], sdk.NewDecWithPrec(8, 2)),
		NewSchedule(monthlyTimes[14], sdk.NewDecWithPrec(9, 2)),
		NewSchedule(monthlyTimes[15], sdk.NewDecWithPrec(9, 2)),
		NewSchedule(monthlyTimes[16], sdk.NewDecWithPrec(9, 2)),
		NewSchedule(monthlyTimes[17], sdk.NewDecWithPrec(9, 2)),
	})

	employeeSchedule = NewVestingSchedule(assets.MicroLunaDenom, []Schedule{
		NewSchedule(monthlyTimes[0], sdk.NewDecWithPrec(5, 2)),
		NewSchedule(monthlyTimes[12], sdk.NewDecWithPrec(29, 2)),
		NewSchedule(monthlyTimes[24], sdk.NewDecWithPrec(33, 2)),
		NewSchedule(monthlyTimes[36], sdk.NewDecWithPrec(33, 2)),
	})

}

func scaleCoins(scale float64, denom string, input sdk.Coins) sdk.Coins {
	output := sdk.Coins{}
	for _, coin := range input {
		if coin.Denom != denom {
			continue
		}

		decScale := sdk.NewDecWithPrec(int64(scale*100), 2)
		output = append(output, sdk.NewCoin(coin.Denom, decScale.MulInt(coin.Amount).RoundInt()))
	}
	return output
}

func TestGetVestedCoinsGradVestingAcc(t *testing.T) {
	genesis := timeGenesis

	_, _, addr := keyPubAddr()
	origCoins := sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 10000)}
	bacc := auth.NewBaseAccountWithAddress(addr)
	bacc.SetCoins(origCoins)

	// require no coins are vested until schedule maturation
	gva := NewBaseGradedVestingAccount(&bacc, []VestingSchedule{
		angelSchedule,
	})
	vestedCoins := gva.GetVestedCoins(genesis)
	require.Nil(t, vestedCoins)

	// require coins be vested at the expected cliff
	vestedCoins = gva.GetVestedCoins(timeGenesis.AddDate(0, 1, 0))
	require.Equal(t, scaleCoins(0.1, assets.MicroLunaDenom, origCoins), vestedCoins)

	vestedCoins = gva.GetVestedCoins(timeGenesis.AddDate(0, 2, 0))
	require.Equal(t, scaleCoins(0.2, assets.MicroLunaDenom, origCoins), vestedCoins)

	vestedCoins = gva.GetVestedCoins(timeGenesis.AddDate(0, 3, 0))
	require.Equal(t, scaleCoins(0.3, assets.MicroLunaDenom, origCoins), vestedCoins)

	vestedCoins = gva.GetVestedCoins(timeGenesis.AddDate(0, 4, 0))
	require.Equal(t, scaleCoins(0.3, assets.MicroLunaDenom, origCoins), vestedCoins)

	vestedCoins = gva.GetVestedCoins(timeGenesis.AddDate(0, 5, 0))
	require.Equal(t, scaleCoins(0.3, assets.MicroLunaDenom, origCoins), vestedCoins)

	vestedCoins = gva.GetVestedCoins(timeGenesis.AddDate(1, 0, 0))
	require.Equal(t, scaleCoins(1.0, assets.MicroLunaDenom, origCoins), vestedCoins)
}

func TestGetVestingCoinsGradVestingAcc(t *testing.T) {
	genesis := timeGenesis

	_, _, addr := keyPubAddr()
	origCoins := sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 10000)}
	bacc := auth.NewBaseAccountWithAddress(addr)
	bacc.SetCoins(origCoins)

	// require no coins are vested until schedule maturation
	gva := NewBaseGradedVestingAccount(&bacc, []VestingSchedule{
		angelSchedule,
	})
	vestingCoins := gva.GetVestingCoins(genesis)
	require.Equal(t, vestingCoins, origCoins)

	// require coins be vested at the expected cliff
	vestingCoins = gva.GetVestingCoins(timeGenesis.AddDate(0, 1, 0))
	require.Equal(t, origCoins.Sub(scaleCoins(0.1, assets.MicroLunaDenom, origCoins)), vestingCoins)

	vestingCoins = gva.GetVestingCoins(timeGenesis.AddDate(0, 2, 0))
	require.Equal(t, origCoins.Sub(scaleCoins(0.2, assets.MicroLunaDenom, origCoins)), vestingCoins)

	vestingCoins = gva.GetVestingCoins(timeGenesis.AddDate(0, 3, 0))
	require.Equal(t, origCoins.Sub(scaleCoins(0.3, assets.MicroLunaDenom, origCoins)), vestingCoins)

	vestingCoins = gva.GetVestingCoins(timeGenesis.AddDate(0, 4, 0))
	require.Equal(t, origCoins.Sub(scaleCoins(0.3, assets.MicroLunaDenom, origCoins)), vestingCoins)

	vestingCoins = gva.GetVestingCoins(timeGenesis.AddDate(0, 5, 0))
	require.Equal(t, origCoins.Sub(scaleCoins(0.3, assets.MicroLunaDenom, origCoins)), vestingCoins)

	vestingCoins = gva.GetVestingCoins(timeGenesis.AddDate(1, 0, 0))
	require.Equal(t, origCoins.Sub(scaleCoins(1.0, assets.MicroLunaDenom, origCoins)), vestingCoins)
}

func TestSpendableCoinsGradVestingAcc(t *testing.T) {
	genesis := timeGenesis

	_, _, addr := keyPubAddr()
	origCoins := sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 10000), sdk.NewInt64Coin(assets.MicroSDRDenom, 10000)}
	sdrCoins := sdk.Coins{sdk.NewInt64Coin(assets.MicroSDRDenom, 10000)}
	bacc := auth.NewBaseAccountWithAddress(addr)
	bacc.SetCoins(origCoins)

	// require no coins are vested until schedule maturation
	gva := NewBaseGradedVestingAccount(&bacc, []VestingSchedule{
		angelSchedule,
	})

	spendableCoins := gva.SpendableCoins(genesis)
	require.Equal(t, sdrCoins, spendableCoins)

	// require that all coins are spendable after the maturation of the vesting
	// schedule
	spendableCoins = gva.SpendableCoins(timeGenesis.AddDate(1, 0, 0))
	require.Equal(t, origCoins, spendableCoins)

	// require that all luna coins are still vesting after some time
	spendableCoins = gva.SpendableCoins(genesis.Add(12 * time.Hour))
	require.Equal(t, spendableCoins, sdrCoins)

	// receive some coins
	regvamt := sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 50)}
	gva.SetCoins(gva.GetCoins().Add(regvamt))

	// require that only received coins and sdrCoins are spendable since the account is still
	// vesting
	spendableCoins = gva.SpendableCoins(genesis.Add(12 * time.Hour))
	require.Equal(t, regvamt.Add(sdrCoins), spendableCoins)

	// spend all spendable coins
	gva.SetCoins(gva.GetCoins().Sub(spendableCoins))

	// require that no more coins are spendable
	spendableCoins = gva.SpendableCoins(genesis.Add(12 * time.Hour))
	require.Nil(t, spendableCoins)
}

func TestTrackDelegationGradVestingAcc(t *testing.T) {
	genesis := timeGenesis

	_, _, addr := keyPubAddr()
	origCoins := sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 10000)}
	bacc := auth.NewBaseAccountWithAddress(addr)

	bacc.SetCoins(origCoins)
	gva := NewBaseGradedVestingAccount(&bacc, []VestingSchedule{
		angelSchedule,
	})
	gva.TrackDelegation(genesis, origCoins)
	require.Equal(t, origCoins, gva.DelegatedVesting)
	require.Nil(t, gva.DelegatedFree)
	require.Nil(t, gva.GetCoins())

	// require the ability to delegate all vested coins
	bacc.SetCoins(origCoins)
	gva = NewBaseGradedVestingAccount(&bacc, []VestingSchedule{
		angelSchedule,
	})
	gva.TrackDelegation(timeGenesis.AddDate(1, 0, 0), origCoins)
	require.Nil(t, gva.DelegatedVesting)
	require.Equal(t, origCoins, gva.DelegatedFree)
	require.Nil(t, gva.GetCoins())

	// require the ability to delegate all coins half way through the vesting
	// schedule
	bacc.SetCoins(origCoins)
	gva = NewBaseGradedVestingAccount(&bacc, []VestingSchedule{
		angelSchedule,
	})
	gva.TrackDelegation(genesis.AddDate(0, 3, 0), sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 3000)})
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 3000)}, gva.DelegatedVesting)
	require.Nil(t, gva.DelegatedFree)

	gva.TrackDelegation(genesis.AddDate(0, 3, 0), sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 7000)})
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 7000)}, gva.DelegatedVesting)
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 3000)}, gva.DelegatedFree)
	require.Nil(t, gva.GetCoins())

	// require no modifications when delegation amount is zero or not enough funds
	bacc.SetCoins(origCoins)
	gva = NewBaseGradedVestingAccount(&bacc, []VestingSchedule{
		angelSchedule,
	})

	require.Panics(t, func() {
		gva.TrackDelegation(genesis.AddDate(1, 0, 0), sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 1000000)})
	})
	require.Nil(t, gva.DelegatedVesting)
	require.Nil(t, gva.DelegatedFree)
	require.Equal(t, origCoins, gva.GetCoins())
}

func TestTrackUndelegationGradVestingAcc(t *testing.T) {
	genesis := timeGenesis

	_, _, addr := keyPubAddr()
	origCoins := sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 10000), sdk.NewInt64Coin(assets.MicroSDRDenom, 10000)}
	bacc := auth.NewBaseAccountWithAddress(addr)
	bacc.SetCoins(origCoins)

	// require the ability to undelegate all vesting coins
	gva := NewBaseGradedVestingAccount(&bacc, []VestingSchedule{
		angelSchedule,
	})
	gva.TrackDelegation(genesis, origCoins)
	gva.TrackUndelegation(origCoins)
	require.Nil(t, gva.DelegatedFree)
	require.Nil(t, gva.DelegatedVesting)
	require.Equal(t, origCoins, gva.GetCoins())

	// require the ability to undelegate all vested coins
	bacc.SetCoins(origCoins)
	gva = NewBaseGradedVestingAccount(&bacc, []VestingSchedule{
		angelSchedule,
	})
	gva.TrackDelegation(genesis.AddDate(1, 0, 0), origCoins)
	gva.TrackUndelegation(origCoins)
	require.Nil(t, gva.DelegatedFree)
	require.Nil(t, gva.DelegatedVesting)
	require.Equal(t, origCoins, gva.GetCoins())

	// require no modifications when the undelegation amount is zero
	bacc.SetCoins(origCoins)
	gva = NewBaseGradedVestingAccount(&bacc, []VestingSchedule{
		angelSchedule,
	})
	require.Panics(t, func() {
		gva.TrackUndelegation(sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 0)})
	})
	require.Nil(t, gva.DelegatedFree)
	require.Nil(t, gva.DelegatedVesting)
	require.Equal(t, origCoins, gva.GetCoins())

	// vest 50% and delegate to two validators
	gva = NewBaseGradedVestingAccount(&bacc, []VestingSchedule{
		angelSchedule,
	})
	gva.TrackDelegation(genesis.AddDate(0, 3, 0), sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 5000)})
	gva.TrackDelegation(genesis.AddDate(0, 3, 0), sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 5000)})

	// undelegate from one validator that got slashed 50%
	gva.TrackUndelegation(sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 2500)})
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 500)}, gva.DelegatedFree)
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 7000)}, gva.DelegatedVesting)
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 2500), sdk.NewInt64Coin(assets.MicroSDRDenom, 10000)}, gva.GetCoins())

	// undelegate from the other validator that did not get slashed
	gva.TrackUndelegation(sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 5000)})
	require.Nil(t, gva.DelegatedFree)
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 2500)}, gva.DelegatedVesting)
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 7500), sdk.NewInt64Coin(assets.MicroSDRDenom, 10000)}, gva.GetCoins())
}

func TestStringGradVestingAcc(t *testing.T) {

	_, _, addr := keyPubAddr()
	origCoins := sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 10000), sdk.NewInt64Coin(assets.MicroSDRDenom, 10000)}
	bacc := auth.NewBaseAccountWithAddress(addr)
	bacc.SetCoins(origCoins)

	// require the ability to undelegate all vesting coins
	gva := NewBaseGradedVestingAccount(&bacc, []VestingSchedule{
		angelSchedule,
	})
	require.NotNil(t, gva.String())

	vestingSchedule, found := gva.GetVestingSchedule(assets.MicroLunaDenom)
	require.True(t, found)
	require.NotNil(t, vestingSchedule.String())
}

func TestIsValidGradVestingAcc(t *testing.T) {

	_, _, addr := keyPubAddr()
	origCoins := sdk.Coins{sdk.NewInt64Coin(assets.MicroLunaDenom, 10000), sdk.NewInt64Coin(assets.MicroSDRDenom, 10000)}
	bacc := auth.NewBaseAccountWithAddress(addr)
	bacc.SetCoins(origCoins)

	// require the ability to undelegate all vesting coins
	angelAccount := NewBaseGradedVestingAccount(&bacc, []VestingSchedule{
		angelSchedule,
	})

	vestingSchedule, found := angelAccount.GetVestingSchedule(assets.MicroLunaDenom)
	require.True(t, found)
	require.True(t, vestingSchedule.IsValid())

	seedAccount := NewBaseGradedVestingAccount(&bacc, []VestingSchedule{
		seedSchedule,
	})

	vestingSchedule, found = seedAccount.GetVestingSchedule(assets.MicroLunaDenom)
	require.True(t, found)
	require.True(t, vestingSchedule.IsValid())

	privateAccount := NewBaseGradedVestingAccount(&bacc, []VestingSchedule{
		privateSchedule,
	})

	vestingSchedule, found = privateAccount.GetVestingSchedule(assets.MicroLunaDenom)
	require.True(t, found)
	require.True(t, vestingSchedule.IsValid())

	employeeAccount := NewBaseGradedVestingAccount(&bacc, []VestingSchedule{
		employeeSchedule,
	})

	vestingSchedule, found = employeeAccount.GetVestingSchedule(assets.MicroLunaDenom)
	require.True(t, found)
	require.True(t, vestingSchedule.IsValid())
}
