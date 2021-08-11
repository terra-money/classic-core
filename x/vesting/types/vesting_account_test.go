package types_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authvestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"

	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/vesting/types"
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
	lgva := types.NewLazyGradedVestingAccount(bacc, origCoins, types.VestingSchedules{
		types.NewVestingSchedule(feeDenom, []types.Schedule{
			types.NewSchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
		types.NewVestingSchedule(stakeDenom, []types.Schedule{
			types.NewSchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
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
	lgva := types.NewLazyGradedVestingAccount(bacc, origCoins, types.VestingSchedules{
		types.NewVestingSchedule(feeDenom, []types.Schedule{
			types.NewSchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
		types.NewVestingSchedule(stakeDenom, []types.Schedule{
			types.NewSchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
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

func TestLockedCoinsLazyVestingAcc(t *testing.T) {
	now := tmtime.Now()
	endTime := now.Add(24 * time.Hour)

	_, _, addr := KeyTestPubAddr()
	origCoins := sdk.Coins{sdk.NewInt64Coin(feeDenom, 1000), sdk.NewInt64Coin(stakeDenom, 100)}
	bacc := authtypes.NewBaseAccountWithAddress(addr)
	lgva := types.NewLazyGradedVestingAccount(bacc, origCoins, types.VestingSchedules{
		types.NewVestingSchedule(feeDenom, []types.Schedule{
			types.NewSchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
		types.NewVestingSchedule(stakeDenom, []types.Schedule{
			types.NewSchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
	})

	// require that there exist no spendable coins in the beginning of the
	// vesting schedule
	lockedCoins := lgva.LockedCoins(now)
	require.Equal(t, origCoins, lockedCoins)

	// require that all original coins are spendable at the end of the vesting
	// schedule
	lockedCoins = lgva.LockedCoins(endTime)
	require.Equal(t, sdk.NewCoins(), lockedCoins)

	// require that all vested coins (50%) are spendable
	lockedCoins = lgva.LockedCoins(now.Add(12 * time.Hour))
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(feeDenom, 500), sdk.NewInt64Coin(stakeDenom, 50)}, lockedCoins)
}

func TestTrackDelegationLazyVestingAcc(t *testing.T) {
	now := tmtime.Now()
	endTime := now.Add(24 * time.Hour)

	_, _, addr := KeyTestPubAddr()
	origCoins := sdk.Coins{sdk.NewInt64Coin(feeDenom, 1000), sdk.NewInt64Coin(stakeDenom, 100)}
	bacc := authtypes.NewBaseAccountWithAddress(addr)

	// require the ability to delegate all vesting coins
	lgva := types.NewLazyGradedVestingAccount(bacc, origCoins, types.VestingSchedules{
		types.NewVestingSchedule(feeDenom, []types.Schedule{
			types.NewSchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
		types.NewVestingSchedule(stakeDenom, []types.Schedule{
			types.NewSchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
	})

	lgva.TrackDelegation(now, origCoins, origCoins)
	require.Equal(t, origCoins, lgva.DelegatedVesting)
	require.True(t, lgva.DelegatedFree.Empty())

	// require the ability to delegate all vested coins
	lgva = types.NewLazyGradedVestingAccount(bacc, origCoins, types.VestingSchedules{
		types.NewVestingSchedule(feeDenom, []types.Schedule{
			types.NewSchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
		types.NewVestingSchedule(stakeDenom, []types.Schedule{
			types.NewSchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
	})

	lgva.TrackDelegation(endTime, origCoins, origCoins)
	require.True(t, lgva.DelegatedVesting.Empty())
	require.Equal(t, origCoins, lgva.DelegatedFree)

	// require the ability to delegate all vesting coins (50%) and all vested coins (50%)
	lgva = types.NewLazyGradedVestingAccount(bacc, origCoins, types.VestingSchedules{
		types.NewVestingSchedule(feeDenom, []types.Schedule{
			types.NewSchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
		types.NewVestingSchedule(stakeDenom, []types.Schedule{
			types.NewSchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
	})

	lgva.TrackDelegation(now.Add(12*time.Hour), origCoins, sdk.Coins{sdk.NewInt64Coin(stakeDenom, 50)})
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(stakeDenom, 50)}, lgva.DelegatedVesting)
	require.True(t, lgva.DelegatedFree.Empty())

	lgva.TrackDelegation(now.Add(12*time.Hour), origCoins, sdk.Coins{sdk.NewInt64Coin(stakeDenom, 50)})
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(stakeDenom, 50)}, lgva.DelegatedVesting)
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(stakeDenom, 50)}, lgva.DelegatedFree)

	// require no modifications when delegation amount is zero or not enough funds
	lgva = types.NewLazyGradedVestingAccount(bacc, origCoins, types.VestingSchedules{
		types.NewVestingSchedule(feeDenom, []types.Schedule{
			types.NewSchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
		types.NewVestingSchedule(stakeDenom, []types.Schedule{
			types.NewSchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
	})

	require.Panics(t, func() {
		lgva.TrackDelegation(endTime, origCoins, sdk.Coins{sdk.NewInt64Coin(stakeDenom, 1000000)})
	})
	require.True(t, lgva.DelegatedVesting.Empty())
	require.True(t, lgva.DelegatedFree.Empty())
}

func TestTrackUndelegationLazyVestingAcc(t *testing.T) {
	now := tmtime.Now()
	endTime := now.Add(24 * time.Hour)

	_, _, addr := KeyTestPubAddr()
	origCoins := sdk.Coins{sdk.NewInt64Coin(feeDenom, 1000), sdk.NewInt64Coin(stakeDenom, 100)}
	bacc := authtypes.NewBaseAccountWithAddress(addr)

	// require the ability to undelegate all vesting coins
	lgva := types.NewLazyGradedVestingAccount(bacc, origCoins, types.VestingSchedules{
		types.NewVestingSchedule(feeDenom, []types.Schedule{
			types.NewSchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
		types.NewVestingSchedule(stakeDenom, []types.Schedule{
			types.NewSchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
	})

	lgva.TrackDelegation(now, origCoins, origCoins)
	lgva.TrackUndelegation(origCoins)
	require.True(t, lgva.DelegatedFree.Empty())
	require.True(t, lgva.DelegatedVesting.Empty())

	// require the ability to undelegate all vested coins
	lgva = types.NewLazyGradedVestingAccount(bacc, origCoins, types.VestingSchedules{
		types.NewVestingSchedule(feeDenom, []types.Schedule{
			types.NewSchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
		types.NewVestingSchedule(stakeDenom, []types.Schedule{
			types.NewSchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
	})

	lgva.TrackDelegation(endTime, origCoins, origCoins)
	lgva.TrackUndelegation(origCoins)
	require.True(t, lgva.DelegatedFree.Empty())
	require.True(t, lgva.DelegatedVesting.Empty())

	// require no modifications when the undelegation amount is zero
	lgva = types.NewLazyGradedVestingAccount(bacc, origCoins, types.VestingSchedules{
		types.NewVestingSchedule(feeDenom, []types.Schedule{
			types.NewSchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
		types.NewVestingSchedule(stakeDenom, []types.Schedule{
			types.NewSchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
	})

	require.Panics(t, func() {
		lgva.TrackUndelegation(sdk.Coins{sdk.NewInt64Coin(stakeDenom, 0)})
	})
	require.True(t, lgva.DelegatedFree.Empty())
	require.True(t, lgva.DelegatedVesting.Empty())

	// vest 50% and delegate to two validators
	lgva = types.NewLazyGradedVestingAccount(bacc, origCoins, types.VestingSchedules{
		types.NewVestingSchedule(feeDenom, []types.Schedule{
			types.NewSchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
		types.NewVestingSchedule(stakeDenom, []types.Schedule{
			types.NewSchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
	})

	lgva.TrackDelegation(now.Add(12*time.Hour), origCoins, sdk.Coins{sdk.NewInt64Coin(stakeDenom, 50)})
	lgva.TrackDelegation(now.Add(12*time.Hour), origCoins, sdk.Coins{sdk.NewInt64Coin(stakeDenom, 50)})

	// undelegate from one validator that got slashed 50%
	lgva.TrackUndelegation(sdk.Coins{sdk.NewInt64Coin(stakeDenom, 25)})
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(stakeDenom, 25)}, lgva.DelegatedFree)
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(stakeDenom, 50)}, lgva.DelegatedVesting)

	// undelegate from the other validator that did not get slashed
	lgva.TrackUndelegation(sdk.Coins{sdk.NewInt64Coin(stakeDenom, 50)})
	require.True(t, lgva.DelegatedFree.Empty())
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(stakeDenom, 25)}, lgva.DelegatedVesting)
}

func TestGenesisAccountValidate(t *testing.T) {
	pubkey := secp256k1.GenPrivKey().PubKey()
	addr := sdk.AccAddress(pubkey.Address())
	baseAcc := authtypes.NewBaseAccount(addr, pubkey, 0, 0)
	initialVesting := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 50))
	baseVestingWithCoins := authvestingtypes.NewBaseVestingAccount(baseAcc, initialVesting, 100)
	tests := []struct {
		name   string
		acc    authtypes.GenesisAccount
		expErr error
	}{
		{
			"valid base account",
			baseAcc,
			nil,
		},
		{
			"invalid base valid account",
			authtypes.NewBaseAccount(addr, secp256k1.GenPrivKey().PubKey(), 0, 0),
			errors.New("account address and pubkey address do not match"),
		},
		{
			"valid base vesting account",
			baseVestingWithCoins,
			nil,
		},
		{
			"valid continuous vesting account",
			types.NewLazyGradedVestingAccount(baseAcc, initialVesting, types.VestingSchedules{types.VestingSchedule{core.MicroLunaDenom, types.Schedules{types.Schedule{1554668078, 1654668078, sdk.OneDec()}}}}),
			nil,
		},
		{
			"invalid vesting times",
			types.NewLazyGradedVestingAccount(baseAcc, initialVesting, types.VestingSchedules{types.VestingSchedule{core.MicroLunaDenom, types.Schedules{types.Schedule{1654668078, 1554668078, sdk.OneDec()}}}}),
			errors.New("vesting start-time cannot be before end-time"),
		},
		{
			"invalid vesting times 2",
			types.NewLazyGradedVestingAccount(baseAcc, initialVesting, types.VestingSchedules{types.VestingSchedule{core.MicroLunaDenom, types.Schedules{types.Schedule{-1, 1554668078, sdk.OneDec()}}}}),
			errors.New("vesting start-time cannot be negative"),
		},
		{
			"invalid vesting ratio",
			types.NewLazyGradedVestingAccount(baseAcc, initialVesting, types.VestingSchedules{types.VestingSchedule{core.MicroLunaDenom, types.Schedules{types.Schedule{1554668078, 1654668078, sdk.ZeroDec()}}}}),
			errors.New("vesting ratio cannot be smaller than or equal with zero"),
		},
		{
			"invalid vesting sum of ratio",
			types.NewLazyGradedVestingAccount(baseAcc, initialVesting, types.VestingSchedules{types.VestingSchedule{core.MicroLunaDenom, types.Schedules{types.Schedule{1554668078, 1654668078, sdk.NewDecWithPrec(1, 1)}}}}),
			errors.New("vesting total ratio must be one"),
		},
		{
			"multiple vesting schedule for a denom",
			types.NewLazyGradedVestingAccount(baseAcc, initialVesting, types.VestingSchedules{
				{core.MicroLunaDenom, types.Schedules{types.Schedule{1554668078, 1654668078, sdk.OneDec()}}},
				{core.MicroLunaDenom, types.Schedules{types.Schedule{1554668078, 1654668078, sdk.OneDec()}}},
			}),
			errors.New("cannot have multiple vesting schedules for uluna"),
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

func TestBaseVestingAccountMarshal(t *testing.T) {
	pubkey := secp256k1.GenPrivKey().PubKey()
	addr := sdk.AccAddress(pubkey.Address())
	coins := sdk.NewCoins(sdk.NewInt64Coin("test", 5))
	baseAcc := authtypes.NewBaseAccount(addr, pubkey, 10, 50)

	acc := authvestingtypes.NewBaseVestingAccount(baseAcc, coins, time.Now().Unix())

	cdc := MakeTestCodec(t)
	bz, err := cdc.MarshalInterface(acc)
	require.Nil(t, err)

	var acc2 authtypes.AccountI
	err = cdc.UnmarshalInterface(bz, &acc2)

	require.Nil(t, err)
	require.IsType(t, &authvestingtypes.BaseVestingAccount{}, acc2)
	require.Equal(t, acc.String(), acc2.String())
}

func TestLazyGradedVestingAccountMarshal(t *testing.T) {
	now := tmtime.Now()
	endTime := now.Add(24 * time.Hour)

	pubkey := secp256k1.GenPrivKey().PubKey()
	addr := sdk.AccAddress(pubkey.Address())
	coins := sdk.NewCoins(sdk.NewInt64Coin("test", 5))
	baseAcc := authtypes.NewBaseAccount(addr, pubkey, 10, 50)

	baseVesting := authvestingtypes.NewBaseVestingAccount(baseAcc, coins, now.Unix())
	acc := types.NewLazyGradedVestingAccountRaw(baseVesting, types.VestingSchedules{
		types.NewVestingSchedule(feeDenom, []types.Schedule{
			types.NewSchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
		types.NewVestingSchedule(stakeDenom, []types.Schedule{
			types.NewSchedule(now.Unix(), endTime.Unix(), sdk.NewDec(1)),
		}),
	})

	cdc := MakeTestCodec(t)
	bz, err := cdc.MarshalInterface(acc)
	require.Nil(t, err)

	var acc2 authtypes.AccountI
	err = cdc.UnmarshalInterface(bz, &acc2)

	require.Nil(t, err)
	require.IsType(t, &types.LazyGradedVestingAccount{}, acc2)
	require.Equal(t, acc.String(), acc2.String())
}
