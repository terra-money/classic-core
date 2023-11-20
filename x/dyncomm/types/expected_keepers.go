package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// AccountKeeper is expected keeper for auth module
type StakingKeeper interface {
	MinCommissionRate(ctx sdk.Context) sdk.Dec
	GetLastTotalPower(ctx sdk.Context) math.Int
	PowerReduction(ctx sdk.Context) math.Int
	IterateValidators(sdk.Context, func(index int64, validator stakingtypes.ValidatorI) (stop bool))
	SetValidator(ctx sdk.Context, validator stakingtypes.Validator)
}
