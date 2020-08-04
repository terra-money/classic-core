package types

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authvesttypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
)

var (
	stakeDenom = "stake"
	feeDenom   = "fee"
)

func TestGetVestedCoinsLazyVestingAcc(t *testing.T) {
	now := tmtime.Now()
	endTime := now.Add(24 * time.Hour)

	_, _, addr := KeyTestPubAddr()
	origCoins := sdk.Coins{sdk.NewInt64Coin(feeDenom, 1000), sdk.NewInt64Coin(stakeDenom, 100)}
	bacc := authtypes.NewBaseAccountWithAddress(addr)
	bacc.SetCoins(origCoins)
	lgva := NewLazyGradedVestingAccount(&bacc, VestingSchedules{
		NewVestingSchedule(feeDenom, []LazySchedule{
			NewLazySchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
		NewVestingSchedule(stakeDenom, []LazySchedule{
			NewLazySchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
	})

	// require no coins vested in the very beginning of the vesting schedule
	vestedCoins := lgva.GetVestedCoins(now)
	require.Nil(t, vestedCoins)

	// require all coins vested at the end of the vesting schedule
	vestedCoins = lgva.GetVestedCoins(endTime)
	require.Equal(t, origCoins, vestedCoins)

	// require 50% of coins vested
	vestedCoins = lgva.GetVestedCoins(now.Add(12 * time.Hour))
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(feeDenom, 500), sdk.NewInt64Coin(stakeDenom, 50)}, vestedCoins)

	// require 100% of coins vested
	vestedCoins = lgva.GetVestedCoins(now.Add(48 * time.Hour))
	require.Equal(t, origCoins, vestedCoins)
}

func TestGetVestingCoinsLazyVestingAcc(t *testing.T) {
	now := tmtime.Now()
	endTime := now.Add(24 * time.Hour)

	_, _, addr := KeyTestPubAddr()
	origCoins := sdk.Coins{sdk.NewInt64Coin(feeDenom, 1000), sdk.NewInt64Coin(stakeDenom, 100)}
	bacc := authtypes.NewBaseAccountWithAddress(addr)
	bacc.SetCoins(origCoins)
	lgva := NewLazyGradedVestingAccount(&bacc, VestingSchedules{
		NewVestingSchedule(feeDenom, []LazySchedule{
			NewLazySchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
		NewVestingSchedule(stakeDenom, []LazySchedule{
			NewLazySchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
	})

	// require all coins vesting in the beginning of the vesting schedule
	vestingCoins := lgva.GetVestingCoins(now)
	require.Equal(t, origCoins, vestingCoins)

	// require no coins vesting at the end of the vesting schedule
	vestingCoins = lgva.GetVestingCoins(endTime)
	require.Nil(t, vestingCoins)

	// require 50% of coins vesting
	vestingCoins = lgva.GetVestingCoins(now.Add(12 * time.Hour))
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(feeDenom, 500), sdk.NewInt64Coin(stakeDenom, 50)}, vestingCoins)
}

func TestSpendableCoinsLazyVestingAcc(t *testing.T) {
	now := tmtime.Now()
	endTime := now.Add(24 * time.Hour)

	_, _, addr := KeyTestPubAddr()
	origCoins := sdk.Coins{sdk.NewInt64Coin(feeDenom, 1000), sdk.NewInt64Coin(stakeDenom, 100)}
	bacc := authtypes.NewBaseAccountWithAddress(addr)
	bacc.SetCoins(origCoins)
	lgva := NewLazyGradedVestingAccount(&bacc, VestingSchedules{
		NewVestingSchedule(feeDenom, []LazySchedule{
			NewLazySchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
		NewVestingSchedule(stakeDenom, []LazySchedule{
			NewLazySchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
	})

	// require that there exist no spendable coins in the beginning of the
	// vesting schedule
	spendableCoins := lgva.SpendableCoins(now)
	require.Nil(t, spendableCoins)

	// require that all original coins are spendable at the end of the vesting
	// schedule
	spendableCoins = lgva.SpendableCoins(endTime)
	require.Equal(t, origCoins, spendableCoins)

	// require that all vested coins (50%) are spendable
	spendableCoins = lgva.SpendableCoins(now.Add(12 * time.Hour))
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(feeDenom, 500), sdk.NewInt64Coin(stakeDenom, 50)}, spendableCoins)

	// receive some coins
	recvAmt := sdk.Coins{sdk.NewInt64Coin(stakeDenom, 50)}
	lgva.SetCoins(lgva.GetCoins().Add(recvAmt...))

	// require that all vested coins (50%) are spendable plus any received
	spendableCoins = lgva.SpendableCoins(now.Add(12 * time.Hour))
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(feeDenom, 500), sdk.NewInt64Coin(stakeDenom, 100)}, spendableCoins)

	// spend all spendable coins
	lgva.SetCoins(lgva.GetCoins().Sub(spendableCoins))

	// require that no more coins are spendable
	spendableCoins = lgva.SpendableCoins(now.Add(12 * time.Hour))
	require.Nil(t, spendableCoins)
}

func TestTrackDelegationLazyVestingAcc(t *testing.T) {
	now := tmtime.Now()
	endTime := now.Add(24 * time.Hour)

	_, _, addr := KeyTestPubAddr()
	origCoins := sdk.Coins{sdk.NewInt64Coin(feeDenom, 1000), sdk.NewInt64Coin(stakeDenom, 100)}
	bacc := authtypes.NewBaseAccountWithAddress(addr)
	bacc.SetCoins(origCoins)

	// require the ability to delegate all vesting coins
	lgva := NewLazyGradedVestingAccount(&bacc, VestingSchedules{
		NewVestingSchedule(feeDenom, []LazySchedule{
			NewLazySchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
		NewVestingSchedule(stakeDenom, []LazySchedule{
			NewLazySchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
	})

	lgva.TrackDelegation(now, origCoins)
	require.Equal(t, origCoins, lgva.DelegatedVesting)
	require.Nil(t, lgva.DelegatedFree)

	// require the ability to delegate all vested coins
	bacc.SetCoins(origCoins)
	lgva = NewLazyGradedVestingAccount(&bacc, VestingSchedules{
		NewVestingSchedule(feeDenom, []LazySchedule{
			NewLazySchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
		NewVestingSchedule(stakeDenom, []LazySchedule{
			NewLazySchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
	})

	lgva.TrackDelegation(endTime, origCoins)
	require.Nil(t, lgva.DelegatedVesting)
	require.Equal(t, origCoins, lgva.DelegatedFree)

	// require the ability to delegate all vesting coins (50%) and all vested coins (50%)
	bacc.SetCoins(origCoins)
	lgva = NewLazyGradedVestingAccount(&bacc, VestingSchedules{
		NewVestingSchedule(feeDenom, []LazySchedule{
			NewLazySchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
		NewVestingSchedule(stakeDenom, []LazySchedule{
			NewLazySchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
	})

	lgva.TrackDelegation(now.Add(12*time.Hour), sdk.Coins{sdk.NewInt64Coin(stakeDenom, 50)})
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(stakeDenom, 50)}, lgva.DelegatedVesting)
	require.Nil(t, lgva.DelegatedFree)

	lgva.TrackDelegation(now.Add(12*time.Hour), sdk.Coins{sdk.NewInt64Coin(stakeDenom, 50)})
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(stakeDenom, 50)}, lgva.DelegatedVesting)
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(stakeDenom, 50)}, lgva.DelegatedFree)

	// require no modifications when delegation amount is zero or not enough funds
	bacc.SetCoins(origCoins)
	lgva = NewLazyGradedVestingAccount(&bacc, VestingSchedules{
		NewVestingSchedule(feeDenom, []LazySchedule{
			NewLazySchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
		NewVestingSchedule(stakeDenom, []LazySchedule{
			NewLazySchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
	})

	require.Panics(t, func() {
		lgva.TrackDelegation(endTime, sdk.Coins{sdk.NewInt64Coin(stakeDenom, 1000000)})
	})
	require.Nil(t, lgva.DelegatedVesting)
	require.Nil(t, lgva.DelegatedFree)
}

func TestTrackUndelegationLazyVestingAcc(t *testing.T) {
	now := tmtime.Now()
	endTime := now.Add(24 * time.Hour)

	_, _, addr := KeyTestPubAddr()
	origCoins := sdk.Coins{sdk.NewInt64Coin(feeDenom, 1000), sdk.NewInt64Coin(stakeDenom, 100)}
	bacc := authtypes.NewBaseAccountWithAddress(addr)
	bacc.SetCoins(origCoins)

	// require the ability to undelegate all vesting coins
	lgva := NewLazyGradedVestingAccount(&bacc, VestingSchedules{
		NewVestingSchedule(feeDenom, []LazySchedule{
			NewLazySchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
		NewVestingSchedule(stakeDenom, []LazySchedule{
			NewLazySchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
	})

	lgva.TrackDelegation(now, origCoins)
	lgva.TrackUndelegation(origCoins)
	require.Nil(t, lgva.DelegatedFree)
	require.Nil(t, lgva.DelegatedVesting)

	// require the ability to undelegate all vested coins
	bacc.SetCoins(origCoins)
	lgva = NewLazyGradedVestingAccount(&bacc, VestingSchedules{
		NewVestingSchedule(feeDenom, []LazySchedule{
			NewLazySchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
		NewVestingSchedule(stakeDenom, []LazySchedule{
			NewLazySchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
	})

	lgva.TrackDelegation(endTime, origCoins)
	lgva.TrackUndelegation(origCoins)
	require.Nil(t, lgva.DelegatedFree)
	require.Nil(t, lgva.DelegatedVesting)

	// require no modifications when the undelegation amount is zero
	bacc.SetCoins(origCoins)
	lgva = NewLazyGradedVestingAccount(&bacc, VestingSchedules{
		NewVestingSchedule(feeDenom, []LazySchedule{
			NewLazySchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
		NewVestingSchedule(stakeDenom, []LazySchedule{
			NewLazySchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
	})

	require.Panics(t, func() {
		lgva.TrackUndelegation(sdk.Coins{sdk.NewInt64Coin(stakeDenom, 0)})
	})
	require.Nil(t, lgva.DelegatedFree)
	require.Nil(t, lgva.DelegatedVesting)

	// vest 50% and delegate to two validators
	lgva = NewLazyGradedVestingAccount(&bacc, VestingSchedules{
		NewVestingSchedule(feeDenom, []LazySchedule{
			NewLazySchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
		NewVestingSchedule(stakeDenom, []LazySchedule{
			NewLazySchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
	})

	lgva.TrackDelegation(now.Add(12*time.Hour), sdk.Coins{sdk.NewInt64Coin(stakeDenom, 50)})
	lgva.TrackDelegation(now.Add(12*time.Hour), sdk.Coins{sdk.NewInt64Coin(stakeDenom, 50)})

	// undelegate from one validator that got slashed 50%
	lgva.TrackUndelegation(sdk.Coins{sdk.NewInt64Coin(stakeDenom, 25)})
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(stakeDenom, 25)}, lgva.DelegatedFree)
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(stakeDenom, 50)}, lgva.DelegatedVesting)

	// undelegate from the other validator that did not get slashed
	lgva.TrackUndelegation(sdk.Coins{sdk.NewInt64Coin(stakeDenom, 50)})
	require.Nil(t, lgva.DelegatedFree)
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(stakeDenom, 25)}, lgva.DelegatedVesting)
}

func TestNewBaseVestingAccount(t *testing.T) {
	pubkey := secp256k1.GenPrivKey().PubKey()
	addr := sdk.AccAddress(pubkey.Address())
	_, err := authvesttypes.NewBaseVestingAccount(
		authtypes.NewBaseAccount(addr, sdk.NewCoins(), pubkey, 0, 0),
		sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 50)}, 100,
	)
	require.Equal(t, errors.New("vesting amount cannot be greater than total amount"), err)

	_, err = authvesttypes.NewBaseVestingAccount(
		authtypes.NewBaseAccount(addr, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 10)), pubkey, 0, 0),
		sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 50)}, 100,
	)
	require.Equal(t, errors.New("vesting amount cannot be greater than total amount"), err)

	_, err = authvesttypes.NewBaseVestingAccount(
		authtypes.NewBaseAccount(addr, sdk.NewCoins(sdk.NewInt64Coin("uatom", 50), sdk.NewInt64Coin("eth", 50)), pubkey, 0, 0),
		sdk.NewCoins(sdk.NewInt64Coin("uatom", 100), sdk.NewInt64Coin("eth", 20)), 100,
	)
	require.Equal(t, errors.New("vesting amount cannot be greater than total amount"), err)

	_, err = authvesttypes.NewBaseVestingAccount(
		authtypes.NewBaseAccount(addr, sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 50)}, pubkey, 0, 0),
		sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 50)}, 100,
	)
	require.NoError(t, err)

}

func TestGenesisAccountValidate(t *testing.T) {
	pubkey := secp256k1.GenPrivKey().PubKey()
	addr := sdk.AccAddress(pubkey.Address())
	baseAcc := authtypes.NewBaseAccount(addr, nil, pubkey, 0, 0)
	baseAccWithCoins := authtypes.NewBaseAccount(addr, nil, pubkey, 0, 0)
	baseAccWithCoins.SetCoins(sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 50)})
	baseVestingWithCoins, _ := authvesttypes.NewBaseVestingAccount(baseAccWithCoins, baseAccWithCoins.Coins, 100)
	tests := []struct {
		name   string
		acc    authexported.GenesisAccount
		expErr error
	}{
		{
			"valid base account",
			baseAcc,
			nil,
		},
		{
			"invalid base valid account",
			authtypes.NewBaseAccount(addr, sdk.NewCoins(), secp256k1.GenPrivKey().PubKey(), 0, 0),
			errors.New("pubkey and address pair is invalid"),
		},
		{
			"valid base vesting account",
			baseVestingWithCoins,
			nil,
		},
		{
			"valid continuous vesting account",
			NewLazyGradedVestingAccount(baseAcc, VestingSchedules{VestingSchedule{sdk.DefaultBondDenom, LazySchedules{LazySchedule{1554668078, 1654668078, sdk.OneDec()}}}}),
			nil,
		},
		{
			"invalid vesting times",
			NewLazyGradedVestingAccount(baseAcc, VestingSchedules{VestingSchedule{sdk.DefaultBondDenom, LazySchedules{LazySchedule{1654668078, 1554668078, sdk.OneDec()}}}}),
			errors.New("vesting start-time cannot be before end-time"),
		},
		{
			"invalid vesting times 2",
			NewLazyGradedVestingAccount(baseAcc, VestingSchedules{VestingSchedule{sdk.DefaultBondDenom, LazySchedules{LazySchedule{-1, 1554668078, sdk.OneDec()}}}}),
			errors.New("vesting start-time cannot be negative"),
		},
		{
			"invalid vesting ratio",
			NewLazyGradedVestingAccount(baseAcc, VestingSchedules{VestingSchedule{sdk.DefaultBondDenom, LazySchedules{LazySchedule{1554668078, 1654668078, sdk.ZeroDec()}}}}),
			errors.New("vesting ratio cannot be smaller than or equal with zero"),
		},
		{
			"invalid vesting sum of ratio",
			NewLazyGradedVestingAccount(baseAcc, VestingSchedules{VestingSchedule{sdk.DefaultBondDenom, LazySchedules{LazySchedule{1554668078, 1654668078, sdk.NewDecWithPrec(1, 1)}}}}),
			errors.New("vesting total ratio must be one"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := tt.acc.Validate()
			require.Equal(t, tt.expErr, err)
		})
	}
}

func TestBaseVestingAccountJSON(t *testing.T) {
	pubkey := secp256k1.GenPrivKey().PubKey()
	addr := sdk.AccAddress(pubkey.Address())
	coins := sdk.NewCoins(sdk.NewInt64Coin("test", 5))
	baseAcc := authtypes.NewBaseAccount(addr, coins, pubkey, 10, 50)

	acc, err := authvesttypes.NewBaseVestingAccount(baseAcc, coins, time.Now().Unix())
	require.NoError(t, err)

	bz, err := json.Marshal(acc)
	require.NoError(t, err)

	bz1, err := acc.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, string(bz1), string(bz))

	var a authvesttypes.BaseVestingAccount
	require.NoError(t, json.Unmarshal(bz, &a))
	require.Equal(t, acc.String(), a.String())
}

func TestLazyGradedVestingAccountJSON(t *testing.T) {
	now := tmtime.Now()
	endTime := now.Add(24 * time.Hour)

	pubkey := secp256k1.GenPrivKey().PubKey()
	addr := sdk.AccAddress(pubkey.Address())
	coins := sdk.NewCoins(sdk.NewInt64Coin("test", 5))
	baseAcc := authtypes.NewBaseAccount(addr, coins, pubkey, 10, 50)

	baseVesting, err := authvesttypes.NewBaseVestingAccount(baseAcc, coins, now.Unix())
	acc := NewLazyGradedVestingAccountRaw(baseVesting, VestingSchedules{
		NewVestingSchedule(feeDenom, []LazySchedule{
			NewLazySchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
		NewVestingSchedule(stakeDenom, []LazySchedule{
			NewLazySchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
	})
	require.NoError(t, err)

	bz, err := codec.Cdc.MarshalJSON(acc)
	require.NoError(t, err)

	bz1, err := acc.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, string(bz1), string(bz))

	var a LazyGradedVestingAccount
	require.NoError(t, codec.Cdc.UnmarshalJSON(bz, &a))
	require.Equal(t, acc.String(), a.String())
}
