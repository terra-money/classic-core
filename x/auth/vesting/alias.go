package vesting

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"

	"github.com/terra-money/core/x/auth/vesting/types"
)

var (
	// functions aliases
	RegisterCodec                  = types.RegisterCodec
	NewBaseVestingAccount          = authtypes.NewBaseVestingAccount
	NewLazyGradedVestingAccountRaw = types.NewLazyGradedVestingAccountRaw
	NewLazyGradedVestingAccount    = types.NewLazyGradedVestingAccount

	// variable aliases
	VestingCdc = types.VestingCdc
)

type (
	BaseVestingAccount       = authtypes.BaseVestingAccount
	LazyGradedVestingAccount = types.LazyGradedVestingAccount

	LazySchedule     = types.LazySchedule
	LazySchedules    = types.LazySchedules
	VestingSchedule  = types.VestingSchedule
	VestingSchedules = types.VestingSchedules
)
